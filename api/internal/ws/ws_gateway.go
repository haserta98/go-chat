package ws

import (
	"encoding/json"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/cmd"
)

type Echo struct {
	Val int    `json:"val"`
	To  string `json:"to"`
}

type SendGroupMessage struct {
	GroupID string          `json:"groupId"`
	Payload json.RawMessage `json:"payload"`
}

type WsGateway struct {
	httpServer   *cmd.HTTPServerImpl
	Manager      *WsManager
	onConnect    func(client *WsClient) bool
	onDisconnect func(client *WsClient)
}

func NewWsGateway(httpServer *cmd.HTTPServerImpl, manager *WsManager) *WsGateway {
	return &WsGateway{
		httpServer:   httpServer,
		Manager:      manager,
		onConnect:    nil,
		onDisconnect: nil,
	}
}

func (g *WsGateway) Start() {

}

func (g *WsGateway) OnConnect(handler func(client *WsClient) bool) {
	g.onConnect = handler
}

func (g *WsGateway) OnDisconnect(handler func(client *WsClient)) {
	g.onDisconnect = handler
}

func (g *WsGateway) HandleWebSocket() {
	g.httpServer.GetInstance().Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	g.httpServer.GetInstance().Get("/ws", websocket.New(func(c *websocket.Conn) {
		connID := uuid.New().String()

		sessionID := c.Cookies("session_id")

		// Fallback for query param (not recommended but useful if credentials can't be sent)
		if sessionID == "" {
			sessionID = c.Query("session_id")
		}

		if sessionID == "" {
			c.Close()
			return
		}

		// Look up session in Redis
		sessionData, err := g.Manager.redis.Get("session:" + sessionID)
		if err != nil {
			c.Close()
			return
		}

		var claims map[string]interface{}
		if err := json.Unmarshal([]byte(sessionData), &claims); err != nil {
			c.Close()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.Close()
			return
		}

		client := NewWsClient(connID, userID, c)

		if g.onConnect != nil && !g.onConnect(client) {
			c.Close()
			return
		}

		g.Manager.AddClient(client)
		go client.WritePump()

		defer func() {
			g.Manager.RemoveClient(client)
			if g.onDisconnect != nil {
				g.onDisconnect(client)
			}
		}()

		client.ReadPump(g.Manager)
	}))
}
