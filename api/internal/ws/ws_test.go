package ws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
	"github.com/nats-io/nats.go"
)

var rc = internal.NewRedisClient()
var cluster = cmd.NewCluster(rc, "node-1")
var nc, _ = nats.Connect(nats.DefaultURL)

func TestWsManagerAddClientAndSendLocalMessageToUser(t *testing.T) {
	manager := NewWsManager(rc, cluster, nc)
	client := &WsClient{
		ID:     "conn-1",
		UserID: "user-1",
		Send:   make(chan []byte, 1),
	}

	manager.AddClient(client)

	payload := []byte(`{"msg":"hello"}`)
	manager.SendLocalMessageToUser("user-1", payload)

	select {
	case got := <-client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}
}

func TestWsManagerJoinAndLeaveGroup(t *testing.T) {
	manager := NewWsManager(rc, cluster, nc)

	client := &WsClient{ID: "c1", UserID: "user-1"}
	manager.JoinGroup("group-1", client)

	shard := manager.GetShard("user-1")
	if _, ok := shard.rooms["group-1"]; !ok {
		t.Fatal("expected group to be created")
	}
	if _, ok := shard.rooms["group-1"].Clients[client.UserID][client.ID]; !ok {
		t.Fatal("expected user to join the group")
	}

	manager.LeaveGroup("group-1", client)

	if _, ok := shard.rooms["group-1"].Clients[client.UserID][client.ID]; ok {
		t.Fatal("expected user to leave the group")
	}
}

func TestWsGatewayEchoHandler(t *testing.T) {
	manager := NewWsManager(rc, cluster, nc)
	gateway := NewWsGateway(nil, manager)
	gateway.Start()

	receiver := &WsClient{
		ID:     "conn-2",
		UserID: "user-2",
		Send:   make(chan []byte, 1),
	}

	manager.AddClient(receiver)

	handler, ok := manager.handlers["echo"]
	if !ok {
		t.Fatal("expected echo handler to be registered")
	}

	payload, err := json.Marshal(Echo{Val: 41, To: "user-2"})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	handler(nil, payload)

	select {
	case msg := <-receiver.Send:
		var got Echo
		if err := json.Unmarshal(msg, &got); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if got.Val != 42 {
			t.Fatalf("expected echo value to increment to 42, got %d", got.Val)
		}
		if got.To != "user-2" {
			t.Fatalf("expected message target user-2, got %s", got.To)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for echo response")
	}
}

func TestWsManagerSendGroupMessage(t *testing.T) {

	manager := NewWsManager(rc, cluster, nc)
	manager.Start()

	user1Client := &WsClient{
		ID:     "conn-1",
		UserID: "user-1",
		Send:   make(chan []byte, 1),
	}

	user2Client := &WsClient{
		ID:     "conn-2",
		UserID: "user-2",
		Send:   make(chan []byte, 1),
	}

	manager.AddClient(user1Client)
	manager.AddClient(user2Client)

	manager.JoinGroup("group-1", user1Client)
	manager.JoinGroup("group-1", user2Client)

	payload := []byte(`{"msg":"hello"}`)
	go manager.SendSmartGroup(user1Client, "group-1", payload)

	select {
	case got := <-user1Client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}

	select {
	case got := <-user2Client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}
}

func TestWsManagerJoinAndLeaveGroupMessages(t *testing.T) {

	manager := NewWsManager(rc, cluster, nc)
	manager.Start()

	user1Client := &WsClient{
		ID:     "conn-1",
		UserID: "user-1",
		Send:   make(chan []byte, 1),
	}

	user2Client := &WsClient{
		ID:     "conn-2",
		UserID: "user-2",
		Send:   make(chan []byte, 1),
	}

	manager.AddClient(user1Client)
	manager.AddClient(user2Client)

	manager.JoinGroup("group-1", user1Client)
	manager.JoinGroup("group-1", user2Client)

	payload := []byte(`{"msg":"hello"}`)
	manager.SendSmartGroup(user1Client, "group-1", payload)

	select {
	case got := <-user1Client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}

	select {
	case got := <-user2Client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}

	manager.LeaveGroup("group-1", user2Client)
	manager.SendSmartGroup(user1Client, "group-1", payload)

	select {
	case got := <-user1Client.Send:
		if string(got) != string(payload) {
			t.Fatalf("expected %s, got %s", string(payload), string(got))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for websocket payload")
	}

	select {
	case got := <-user2Client.Send:
		t.Fatalf("Not expected message for user-2: %s", string(got))
	default:
		return
	}
}
