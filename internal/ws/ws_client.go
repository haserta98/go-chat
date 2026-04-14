package ws

import (
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/v3/websocket"
)

type WsClient struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

func NewWsClient(id string, userID string, conn *websocket.Conn) *WsClient {
	return &WsClient{
		ID:     id,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
}

func (c *WsClient) WritePump() {
	pongWait := 5 * time.Second
	pingPeriod := (pongWait * 9) / 10

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			log.Printf("Sending ping to userID=%s connID=%s", c.UserID, c.ID)
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WsClient) ReadPump(manager *WsManager) {

	pongWait := 5 * time.Second

	c.Conn.SetReadLimit(4096)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		manager.RemoveClient(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		if len(msg) == 0 {
			continue
		}

		var request EventRequest
		if err := sonic.Unmarshal(msg, &request); err != nil {
			log.Printf("Geçersiz JSON formatı: %v", err)
			continue
		}

		if handler, exists := manager.handlers[request.Type]; exists {
			handler(c, request.Payload)
		} else {
			log.Printf("Kayıtlı olmayan event türü: %s", request.Type)
		}
	}
}
