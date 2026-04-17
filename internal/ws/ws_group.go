package ws

import (
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
)

type WsGroup struct {
	ID      string
	Clients map[string]map[string]*WsClient
	Manager *WsManager
	mu      sync.RWMutex
}

func NewWsGroup(id string, manager *WsManager) *WsGroup {
	return &WsGroup{
		ID:      id,
		Clients: make(map[string]map[string]*WsClient),
		Manager: manager,
		mu:      sync.RWMutex{},
	}
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
	if _, ok := g.Clients[client.UserID]; ok {
		delete(g.Clients[client.UserID], client.ID)
		if len(g.Clients[client.UserID]) == 0 {
			g.Manager.redis.SRem("group:users:"+g.ID, client.UserID)
			delete(g.Clients, client.UserID)
		}
	}
}

func (g *WsGroup) SendMessage(message []byte) {

}

func (m *WsGroup) ListenGroupMessages() {
	channel := fmt.Sprintf("global_room_channel:%s", m.ID)
	pubsub, err := m.Manager.redis.Subscribe(ctx, channel)
	if err != nil {
		return
	}
	for msg := range pubsub {
		var data GlobalRoomMessage
		if err := sonic.Unmarshal([]byte(msg.Payload), &data); err != nil {
			continue
		}

		if clientsInRoom, ok := m.Clients[data.RoomID]; ok {
			for _, client := range clientsInRoom {
				if data.ExcludeSender && client.UserID == data.SenderID {
					continue
				}
				select {
				case client.Send <- data.Payload:
				default:
					// Buffer full, drop to prevent deadlock
				}
			}
		}
	}
}
