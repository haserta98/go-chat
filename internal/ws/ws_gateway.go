package ws

import (
	"encoding/json"
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/cmd"
)

type Echo struct {
	Val int    `json:"val"`
	To  string `json:"to"`
}

type WsGateway struct {
	httpServer *cmd.HTTPServerImpl
	Manager    *WsManager
}

func NewWsGateway(httpServer *cmd.HTTPServerImpl, manager *WsManager) *WsGateway {
	return &WsGateway{
		httpServer: httpServer,
		Manager:    manager,
	}
}

func (g *WsGateway) Start() {

	g.Manager.RegisterEventHandler("echo", func(client *WsClient, payload json.RawMessage) {
		var echo Echo
		if err := sonic.Unmarshal(payload, &echo); err != nil {
			log.Printf("Invalid echo payload: %v", err)
			return
		}
		echo.Val++
		response, _ := sonic.Marshal(echo)

		g.Manager.SendSmart(echo.To, response)
	})

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
		userID := c.Query("userID")
		connID := uuid.New().String()
		log.Print("New WebSocket connection: userID=", userID, " connID=", connID)

		if userID == "" {
			c.Close()
			return
		}

		client := NewWsClient(connID, userID, c)

		g.Manager.AddClient(client)
		go client.WritePump()
		client.ReadPump(g.Manager)
	}))
}
