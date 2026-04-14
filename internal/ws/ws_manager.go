package ws

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
)

type WsShard struct {
	sync.RWMutex
	clients map[string]map[string]*WsClient
	rooms   map[string]map[string]bool
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
	RoomID  string          `json:"room_id"`
	Payload json.RawMessage `json:"payload"`
}

const ShardCount = 512

var (
	ctx    = context.Background()
	NodeID = func() string {
		id := os.Getenv("NODE_ID")
		if id == "" {
			return "node_" + uuid.New().String()
		}
		return "node_" + id
	}()
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

func (m *WsManager) RegisterEventHandler(eventType string, handler EventHandler) {
	m.handlers[eventType] = handler
}

func (m *WsManager) GetShard(userID string) *WsShard {
	h := fnv.New32a()
	h.Write([]byte(userID))
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

func (m *WsManager) JoinGroup(userID string, groupID string) {
	m.redis.SAdd("group:users:"+groupID, userID)

	shard := m.GetShard(userID)
	shard.Lock()
	defer shard.Unlock()

	if _, ok := shard.rooms[groupID]; !ok {
		shard.rooms[groupID] = make(map[string]bool)
	}
	shard.rooms[groupID][userID] = true
}

func (m *WsManager) LeaveGroup(userID string, roomID string) {
	m.redis.SRem("group:users:"+roomID, userID)

	shard := m.GetShard(userID)
	shard.Lock()
	defer shard.Unlock()

	if usersInRoom, ok := shard.rooms[roomID]; ok {
		delete(usersInRoom, userID)

		if len(usersInRoom) == 0 {
			delete(shard.rooms, roomID)
		}
	}
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
	}
	go m.PublishMessageToUser(targetUserID, payload)
}

func (m *WsManager) SendSmartGroup(groupID string, payload []byte) {
	broadcastMsg := &GlobalRoomMessage{
		RoomID:  groupID,
		Payload: payload,
	}
	msg, _ := sonic.Marshal(broadcastMsg)
	m.redis.Publish(ctx, "global_room_channel", msg)
}

func (m *WsManager) SendLocalMessageToUser(targetUserID string, payload []byte) {
	shard := m.GetShard(targetUserID)
	shard.RLock()
	if clients, exists := shard.clients[targetUserID]; exists {
		for _, client := range clients {
			client.Send <- payload
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

func (m *WsManager) ListenGroupMessages() {
	pubsub, err := m.redis.Subscribe(ctx, "global_room_channel")
	if err != nil {
		return
	}
	for msg := range pubsub {
		var data GlobalRoomMessage
		if err := sonic.Unmarshal([]byte(msg.Payload), &data); err != nil {
			continue
		}

		for i := 0; i < ShardCount; i++ {
			shard := m.shards[i]

			shard.RLock()
			if userIDs, ok := shard.rooms[data.RoomID]; ok {
				for userID := range userIDs {
					if devices, ok := shard.clients[userID]; ok {
						for _, client := range devices {
							client.Send <- data.Payload
						}
					}
				}
			}
			shard.RUnlock()
		}
	}
}
