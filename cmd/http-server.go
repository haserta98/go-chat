package cmd

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type HTTPServerImpl struct {
	app  *fiber.App
	port int
}

func NewHTTPServer(port int) *HTTPServerImpl {
	return &HTTPServerImpl{
		app:  fiber.New(),
		port: port,
	}
}

func (s *HTTPServerImpl) Listen() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}

func (s *HTTPServerImpl) GetInstance() *fiber.App {
	return s.app
}
