package ws

import (
	"encoding/json"
	"log"

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
		log.Printf("hello")
		userID := c.Query("userID")
		connID := uuid.New().String()
		log.Print("New WebSocket connection: userID=", userID, " connID=", connID)

		if userID == "" {
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
			c.Close()
			if g.onDisconnect != nil {
				g.onDisconnect(client)
			}
		}()

		client.ReadPump(g.Manager)
	}))
}
