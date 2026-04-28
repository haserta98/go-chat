package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haserta98/go-rest/internal/service"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (h *MessageHandler) GetMessagesBetween(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	otherUserID := c.Params("otherUserID")
	if otherUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "otherUserID is required",
		})
	}

	messages, err := h.service.GetMessagesBetween(userID.(string), otherUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch messages",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   messages,
	})
}

func (h *MessageHandler) GetGroupMessages(c fiber.Ctx) error {
	groupID := c.Params("groupID")
	if groupID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupID is required",
		})
	}

	messages, err := h.service.GetGroupMessages(groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch group messages",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   messages,
	})
}
