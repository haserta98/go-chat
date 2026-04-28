package ws

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/nats-io/nats.go"
)

type WsShard struct {
	sync.RWMutex
	clients map[string]map[string]*WsClient
	rooms   map[string]*WsGroup
}

type GlobalMessage struct {
	TargetUserID string          `json:"targetUserID"`
	Payload      json.RawMessage `json:"payload"`
}

type EventRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(client *WsClient, payload json.RawMessage)

type GlobalRoomMessage struct {
	RoomID        string          `json:"room_id"`
	SenderID      string          `json:"sender_id,omitempty"`
	ExcludeSender bool            `json:"exclude_sender"`
	Payload       json.RawMessage `json:"payload"`
}

const ShardCount = 512

var (
	ctx = context.Background()
)

type WsManager struct {
	shards     [ShardCount]*WsShard
	redis      *internal.RedisClient
	handlers   map[string]EventHandler
	cluster    *cmd.Cluster
	natsConn   *nats.Conn
	subscriber *internal.NatsSubscriber
	publisher  *internal.NatsPublisher
	appContext *cmd.AppContext
}

func NewWsManager(redis *internal.RedisClient, cluster *cmd.Cluster, natsConn *nats.Conn, appCtx *cmd.AppContext) *WsManager {
	manager := &WsManager{
		redis:      redis,
		handlers:   make(map[string]EventHandler),
		cluster:    cluster,
		natsConn:   natsConn,
		publisher:  internal.NewNatsPublisher(natsConn),
		appContext: appCtx,
	}
	for i := 0; i < ShardCount; i++ {
		manager.shards[i] = &WsShard{
			clients: make(map[string]map[string]*WsClient),
			rooms:   make(map[string]*WsGroup),
		}
	}
	return manager
}

func (m *WsManager) Start() {
	go m.ListenInboundMessages()
}

func (m *WsManager) JoinGroup(groupID string, client *WsClient) {
	shard := m.GetShard(groupID)
	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.rooms[groupID]; !ok {
		shard.rooms[groupID] = NewWsGroup(groupID, m)
		go shard.rooms[groupID].ListenGroupMessages()
	}
	shard.rooms[groupID].AddClient(client)
}

func (m *WsManager) LeaveGroup(groupID string, client *WsClient) {
	shard := m.GetShard(groupID)
	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.rooms[groupID]; ok {
		shard.rooms[groupID].RemoveClient(client)
	}
}

func (m *WsManager) RegisterEventHandler(eventType string, handler EventHandler) {
	m.handlers[eventType] = handler
}

func (m *WsManager) GetShard(key string) *WsShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return m.shards[h.Sum32()%ShardCount]
}

func (m *WsManager) AddClient(client *WsClient) {
	shard := m.GetShard(client.UserID)

	shard.Lock()
	if _, exists := shard.clients[client.UserID]; !exists {
		shard.clients[client.UserID] = make(map[string]*WsClient)
	}
	shard.clients[client.UserID][client.ID] = client
	shard.Unlock()

	myGroups, err := m.appContext.GetRepository("Group").(*repository.GroupRepository).GetMyGroups(client.UserID)
	if err != nil {
		return
	}
	for _, group := range myGroups {
		m.JoinGroup(group.ID, client)
	}

	m.redis.SAdd("user:nodes:"+client.UserID, m.cluster.NodeID)
}

func (m *WsManager) RemoveClient(client *WsClient) {
	client.Close()

	var groupIDs []string
	for i := 0; i < ShardCount; i++ {
		shard := m.shards[i]
		shard.RLock()
		for groupID, group := range shard.rooms {
			group.mu.RLock()
			if conns, ok := group.Clients[client.UserID]; ok {
				if _, exists := conns[client.ID]; exists {
					groupIDs = append(groupIDs, groupID)
				}
			}
			group.mu.RUnlock()
		}
		shard.RUnlock()
	}

	for _, groupID := range groupIDs {
		m.LeaveGroup(groupID, client)
	}

	shard := m.GetShard(client.UserID)
	shard.Lock()
	if userClients, exists := shard.clients[client.UserID]; exists {
		delete(userClients, client.ID)
		if len(userClients) == 0 {
			delete(shard.clients, client.UserID)
			shard.Unlock()
			m.redis.SRem("user:nodes:"+client.UserID, m.cluster.NodeID)
			return
		}
	}
	shard.Unlock()
}

func (m *WsManager) IsLocalUser(userID string) bool {
	shard := m.GetShard(userID)
	shard.RLock()
	_, exists := shard.clients[userID]
	shard.RUnlock()
	return exists
}

func (m *WsManager) SendSmart(targetUserID string, payload []byte) {
	if m.IsLocalUser(targetUserID) {
		m.SendLocalMessageToUser(targetUserID, payload)
	} else {
		m.PublishMessageToUser(targetUserID, payload)
	}
}

func (m *WsManager) SendSmartGroup(from *WsClient, groupID string, payload []byte) {
	shard := m.GetShard(from.UserID)
	shard.RLock()
	group, exists := shard.rooms[groupID]
	shard.RUnlock()

	if exists {
		group.SendMessage(from, payload)
	} else {
		m.BroadcastToGroup(from, groupID, payload)
	}
}

func (m *WsManager) SendLocalMessageToUser(targetUserID string, payload []byte) {
	shard := m.GetShard(targetUserID)
	shard.RLock()
	defer shard.RUnlock()

	if clients, exists := shard.clients[targetUserID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- &EventRequest{
				Type:    "send_user_message",
				Payload: payload,
			}:
			default:
			}
		}
	}
}

func (m *WsManager) PublishMessageToUser(targetUserID string, payload []byte) {
	nodes, err := m.redis.SMembers("user:nodes:" + targetUserID)
	if err != nil {
		return
	}

	for _, node := range nodes {
		if node != m.cluster.NodeID {
			if alive, err := m.cluster.IsTargetNodeAlive(node); err != nil || !alive {
				m.redis.SRem("user:nodes:"+targetUserID, node)
				continue
			}
			msg, _ := sonic.Marshal(GlobalMessage{
				TargetUserID: targetUserID,
				Payload:      payload,
			})

			m.publisher.Publish("inbound:"+node, msg)
		}
	}
}

func (m *WsManager) ListenInboundMessages() {

	subscriber := internal.NewNatsSubscriber(m.natsConn, "inbound:"+m.cluster.NodeID)
	ch, err := subscriber.Subscribe()

	if err != nil {
		return
	}
	m.subscriber = subscriber

	for msg := range ch {
		var globalMsg GlobalMessage
		if err := sonic.Unmarshal([]byte(msg.Data), &globalMsg); err != nil {
			continue
		}
		m.SendLocalMessageToUser(globalMsg.TargetUserID, globalMsg.Payload)
	}
}
