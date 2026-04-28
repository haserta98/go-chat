package cmd

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type HTTPServerImpl struct {
	app  *fiber.App
	port int
}

func NewHTTPServer(port int) *HTTPServerImpl {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	return &HTTPServerImpl{
		app:  app,
		port: port,
	}
}

func (s *HTTPServerImpl) Listen() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

func (s *HTTPServerImpl) GetInstance() *fiber.App {
	return s.app
}
