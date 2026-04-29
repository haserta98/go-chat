package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/nats-io/nats.go"
)

type GlobalMessage struct {
	TargetUserID string          `json:"targetUserID"`
	EventType    string          `json:"eventType,omitempty"` // If empty, defaults to "send_user_message"
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

type WsManager struct {
	clientsMu  sync.RWMutex
	clients    map[string]map[string]*WsClient
	roomsMu    sync.RWMutex
	rooms      map[string]*WsGroup
	redis      *internal.RedisClient
	handlers   map[string]EventHandler
	cluster    *cmd.Cluster
	natsConn   *nats.Conn
	subscriber *internal.NatsSubscriber
	publisher  *internal.NatsPublisher
	appContext *cmd.AppContext
}

func NewWsManager(redis *internal.RedisClient, cluster *cmd.Cluster, natsConn *nats.Conn, appCtx *cmd.AppContext) *WsManager {
	return &WsManager{
		clients:    make(map[string]map[string]*WsClient),
		rooms:      make(map[string]*WsGroup),
		redis:      redis,
		handlers:   make(map[string]EventHandler),
		cluster:    cluster,
		natsConn:   natsConn,
		publisher:  internal.NewNatsPublisher(natsConn),
		appContext: appCtx,
	}
}

func (m *WsManager) Start() {
	go m.ListenInboundMessages()
	go m.ListenClusterGroupActions()
}

func (m *WsManager) JoinGroup(groupID string, client *WsClient) {
	m.roomsMu.Lock()
	if _, ok := m.rooms[groupID]; !ok {
		m.rooms[groupID] = NewWsGroup(groupID, m)
		go m.rooms[groupID].ListenGroupMessages()
	}
	group := m.rooms[groupID]
	m.roomsMu.Unlock()

	group.AddClient(client)
	client.addGroup(groupID)
}

func (m *WsManager) LeaveGroup(groupID string, client *WsClient) {
	m.roomsMu.RLock()
	group, ok := m.rooms[groupID]
	m.roomsMu.RUnlock()

	if ok {
		group.RemoveClient(client)
	}
	client.removeGroup(groupID)
}

func (m *WsManager) RegisterEventHandler(eventType string, handler EventHandler) {
	m.handlers[eventType] = handler
}

func (m *WsManager) AddClient(client *WsClient) {
	m.clientsMu.Lock()
	isFirstConnection := len(m.clients[client.UserID]) == 0
	if _, exists := m.clients[client.UserID]; !exists {
		m.clients[client.UserID] = make(map[string]*WsClient)
	}
	m.clients[client.UserID][client.ID] = client
	m.clientsMu.Unlock()

	if m.appContext != nil {
		myGroups, err := m.appContext.GetRepository("Group").(*repository.GroupRepository).GetMyGroups(client.UserID)
		if err == nil {
			for _, group := range myGroups {
				m.JoinGroup(group.ID, client)
			}
		}
	}

	m.redis.SAdd("user:nodes:"+client.UserID, m.cluster.NodeID)

	// Set online presence in Redis and notify contacts
	if isFirstConnection {
		m.redis.Set("user:online:"+client.UserID, "1", 0)
		go m.BroadcastPresence(client.UserID, "online")
	}
}

func (m *WsManager) RemoveClient(client *WsClient) {
	client.Close()

	for _, groupID := range client.getGroups() {
		m.LeaveGroup(groupID, client)
	}

	m.clientsMu.Lock()
	if userClients, exists := m.clients[client.UserID]; exists {
		delete(userClients, client.ID)
		if len(userClients) == 0 {
			delete(m.clients, client.UserID)
			m.clientsMu.Unlock()
			m.redis.SRem("user:nodes:"+client.UserID, m.cluster.NodeID)

			// Clear online presence and notify contacts
			m.redis.Del("user:online:" + client.UserID)
			go m.BroadcastPresence(client.UserID, "offline")
			return
		}
	}
	m.clientsMu.Unlock()
}

func (m *WsManager) IsLocalUser(userID string) bool {
	m.clientsMu.RLock()
	_, exists := m.clients[userID]
	m.clientsMu.RUnlock()
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
	m.roomsMu.RLock()
	group, exists := m.rooms[groupID]
	m.roomsMu.RUnlock()

	if exists {
		group.SendMessage(from, payload)
	} else {
		m.BroadcastToGroup(from, groupID, payload)
	}
}

func (m *WsManager) SendLocalMessageToUser(targetUserID string, payload []byte) {
	eventRequest := EventRequest{
		Type:    "send_user_message",
		Payload: payload,
	}
	m.SendLocalEventToUser(targetUserID, &eventRequest)
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

		eventType := globalMsg.EventType
		if eventType == "" {
			eventType = "send_user_message"
		}

		m.SendLocalEventToUser(globalMsg.TargetUserID, &EventRequest{
			Type:    eventType,
			Payload: globalMsg.Payload,
		})
	}
}

// SendLocalEventToUser sends a raw EventRequest to all local connections of a user.
// Unlike SendLocalMessageToUser, this does not wrap the payload in a "send_user_message" event.
func (m *WsManager) SendLocalEventToUser(targetUserID string, event *EventRequest) {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	if clients, exists := m.clients[targetUserID]; exists {
		for _, client := range clients {
			select {
			case client.Send <- event:
			default:
			}
		}
	}
}

// BroadcastPresence notifies all users who have the given userID as a contact
// that their presence status has changed.
func (m *WsManager) BroadcastPresence(userID string, status string) {
	if m.appContext == nil {
		return
	}

	userRepo, ok := m.appContext.GetRepository("User").(*repository.UserRepository)
	if !ok {
		return
	}

	// Get all users who have this user as a contact (reverse lookup)
	ownerIDs, err := userRepo.GetContactOwners(userID)
	if err != nil {
		log.Printf("Failed to get contact owners for presence broadcast: %v", err)
		return
	}

	payload, _ := sonic.Marshal(map[string]string{
		"user_id": userID,
		"status":  status,
	})

	event := &EventRequest{
		Type:    "presence_change",
		Payload: payload,
	}

	for _, ownerID := range ownerIDs {
		if m.IsLocalUser(ownerID) {
			m.SendLocalEventToUser(ownerID, event)
		}

		// For remote nodes, publish via NATS with EventType set
		nodes, err := m.redis.SMembers("user:nodes:" + ownerID)
		if err != nil {
			continue
		}
		for _, node := range nodes {
			if node != m.cluster.NodeID {
				msg, _ := sonic.Marshal(GlobalMessage{
					TargetUserID: ownerID,
					EventType:    "presence_change",
					Payload:      payload,
				})
				m.publisher.Publish("inbound:"+node, msg)
			}
		}
	}
}

type ClusterGroupAction struct {
	Action  string `json:"action"`
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}

func (m *WsManager) NotifyGroupAction(action, groupID, userID string) {
	msg, _ := sonic.Marshal(ClusterGroupAction{
		Action:  action,
		GroupID: groupID,
		UserID:  userID,
	})
	m.publisher.Publish("cluster:group_actions", msg)
}

func (m *WsManager) ListenClusterGroupActions() {
	subscriber := internal.NewNatsSubscriber(m.natsConn, "cluster:group_actions")
	ch, err := subscriber.Subscribe()
	if err != nil {
		return
	}

	for msg := range ch {
		var action ClusterGroupAction
		if err := sonic.Unmarshal([]byte(msg.Data), &action); err != nil {
			continue
		}

		m.clientsMu.RLock()
		clients, ok := m.clients[action.UserID]
		m.clientsMu.RUnlock()

		if ok {
			for _, client := range clients {
				switch action.Action {
				case "join":
					m.JoinGroup(action.GroupID, client)
				case "leave":
					m.LeaveGroup(action.GroupID, client)
				}
			}
		}
	}
}
