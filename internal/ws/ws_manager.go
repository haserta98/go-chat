package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
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
	ctx    = context.Background()
	NodeID = os.Getenv("NODE_ID")
)

type WsManager struct {
	shards   [ShardCount]*WsShard
	redis    *internal.RedisClient
	handlers map[string]EventHandler
	cluster  *cmd.Cluster
}

func NewWsManager(redis *internal.RedisClient, cluster *cmd.Cluster) *WsManager {
	manager := &WsManager{
		redis:    redis,
		handlers: make(map[string]EventHandler),
		cluster:  cluster,
	}
	for i := 0; i < ShardCount; i++ {
		manager.shards[i] = &WsShard{
			clients: make(map[string]map[string]*WsClient),
		}
	}
	return manager
}

func (m *WsManager) Start() {
	go m.ListenRedis()
}

func (m *WsManager) JoinGroup(groupID string, client *WsClient) {
	shard := m.GetShard(client.UserID)
	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.rooms[groupID]; !ok {
		shard.rooms[groupID] = NewWsGroup(groupID, m)
		go shard.rooms[groupID].ListenGroupMessages()
	}
	shard.rooms[groupID].AddClient(client)
}

func (m *WsManager) LeaveGroup(groupID string, client *WsClient) {
	shard := m.GetShard(client.UserID)
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

	m.redis.SAdd("user:nodes:"+client.UserID, NodeID)
}

func (m *WsManager) RemoveClient(client *WsClient) {
	shard := m.GetShard(client.UserID)
	shard.Lock()

	if userClients, exists := shard.clients[client.UserID]; exists {
		delete(userClients, client.ID)
		if len(userClients) == 0 {
			delete(shard.clients, client.UserID)
			m.redis.SRem("user:nodes:"+client.UserID, NodeID)
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
		go m.PublishMessageToUser(targetUserID, payload)
	}
}

func (m *WsManager) SendSmartGroup(groupID string, payload []byte) {
	broadcastMsg := &GlobalRoomMessage{
		RoomID:        groupID,
		ExcludeSender: false,
		Payload:       payload,
	}
	msg, _ := sonic.Marshal(broadcastMsg)
	channel := fmt.Sprintf("global_room_channel:%s", groupID)
	m.redis.Publish(ctx, channel, msg)
}

func (m *WsManager) BroadcastToGroup(senderID string, groupID string, payload []byte) {
	broadcastMsg := &GlobalRoomMessage{
		RoomID:        groupID,
		SenderID:      senderID,
		ExcludeSender: true,
		Payload:       payload,
	}
	msg, _ := sonic.Marshal(broadcastMsg)
	channel := fmt.Sprintf("global_room_channel:%s", groupID)
	m.redis.Publish(ctx, channel, msg)
}

func (m *WsManager) SendLocalMessageToUser(targetUserID string, payload []byte) {
	shard := m.GetShard(targetUserID)
	shard.RLock()
	if clients, exists := shard.clients[targetUserID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- payload:
			default:
				// Buffer full, drop to prevent deadlock
			}
		}
	}
	shard.RUnlock()
}

func (m *WsManager) PublishMessageToUser(targetUserID string, payload []byte) {
	nodes, err := m.redis.SMembers("user:nodes:" + targetUserID)
	if err != nil {
		return
	}
	for _, node := range nodes {
		if node != NodeID {
			if alive, err := m.cluster.IsTargetNodeAlive(node); err != nil || !alive {
				m.redis.SRem("user:nodes:"+targetUserID, node)
				continue
			}
			msg, _ := sonic.Marshal(GlobalMessage{
				TargetUserID: targetUserID,
				Payload:      payload,
			})
			m.redis.Publish(ctx, "inbound:"+node, msg)
		}
	}
}

func (m *WsManager) ListenRedis() {
	pubsub, err := m.redis.Subscribe(ctx, "inbound:"+NodeID)
	if err != nil {
		return
	}
	for msg := range pubsub {
		var globalMsg GlobalMessage
		if err := sonic.Unmarshal([]byte(msg.Payload), &globalMsg); err != nil {
			continue
		}
		m.SendLocalMessageToUser(globalMsg.TargetUserID, globalMsg.Payload)
	}
}
