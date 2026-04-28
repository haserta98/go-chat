package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/haserta98/go-rest/internal"
)

type WsGroup struct {
	ID      string
	Clients map[string]map[string]*WsClient
	Manager *WsManager
	mu      sync.RWMutex
	stop    chan struct{}
}

func NewWsGroup(id string, manager *WsManager) *WsGroup {
	return &WsGroup{
		ID:      id,
		Clients: make(map[string]map[string]*WsClient),
		Manager: manager,
		stop:    make(chan struct{}),
	}
}

func (g *WsGroup) Stop() {
	close(g.stop)
}

func (g *WsGroup) AddClient(client *WsClient) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.Clients[client.UserID]; !ok {
		g.Clients[client.UserID] = make(map[string]*WsClient)
	}
	g.Clients[client.UserID][client.ID] = client
	g.Manager.redis.SAdd("group:users:"+g.ID, client.UserID)
}

func (g *WsGroup) RemoveClient(client *WsClient) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.Clients[client.UserID]; ok {
		delete(g.Clients[client.UserID], client.ID)
		if len(g.Clients[client.UserID]) == 0 {
			g.Manager.redis.SRem("group:users:"+g.ID, client.UserID)
			delete(g.Clients, client.UserID)
		}
	}
}

func (g *WsGroup) SendMessage(from *WsClient, message []byte) {
	broadcastMsg := &GlobalRoomMessage{
		RoomID:        g.ID,
		SenderID:      from.UserID,
		ExcludeSender: true,
		Payload:       message,
	}
	msg, _ := sonic.Marshal(broadcastMsg)
	channel := fmt.Sprintf("global_room_channel:%s", g.ID)

	g.Manager.publisher.Publish(channel, msg)
}

func (m *WsManager) BroadcastToGroup(from *WsClient, groupID string, payload []byte) {
	broadcastMsg := &GlobalRoomMessage{
		RoomID:        groupID,
		SenderID:      from.UserID,
		ExcludeSender: true,
		Payload:       payload,
	}
	msg, _ := sonic.Marshal(broadcastMsg)
	channel := fmt.Sprintf("global_room_channel:%s", groupID)
	m.publisher.Publish(channel, msg)
}

func (m *WsGroup) ListenGroupMessages() {
	channel := fmt.Sprintf("global_room_channel:%s", m.ID)
	subscriber := internal.NewNatsSubscriber(m.Manager.natsConn, channel)
	ch, err := subscriber.Subscribe()
	defer subscriber.Unsubscribe()

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-m.stop:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var data GlobalRoomMessage
			if err := sonic.Unmarshal(msg.Data, &data); err != nil {
				log.Println("hata", err)
				continue
			}
			m.mu.RLock()
			for _, clients := range m.Clients {
				for _, client := range clients {
					if data.ExcludeSender && client.UserID == data.SenderID {
						continue
					}
					select {
					case client.Send <- &EventRequest{
						Type:    "send_group_message",
						Payload: data.Payload,
					}:
					default:
						// Buffer full, drop to prevent deadlock
					}
				}
			}
			m.mu.RUnlock()
		}
	}
}
