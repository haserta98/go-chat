package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haserta98/go-rest/internal/dto"
	"github.com/haserta98/go-rest/internal/service"
)

type GroupHandler struct {
	service *service.GroupService
}

func NewGroupHandler(service *service.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) CreateGroup(c fiber.Ctx) error {
	createDTO := new(dto.GroupCreateRequest)
	if err := c.Bind().WithAutoHandling().Body(createDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if createDTO.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name alanı zorunludur",
		})
	}
	group, err := h.service.CreateGroup(createDTO.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Grup oluşturulurken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   group,
	})
}

func (h *GroupHandler) GetGroupByID(c fiber.Ctx) error {
	id := c.Params("id")
	group, err := h.service.GetGroupByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "not-found",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   group,
	})
}

func (h *GroupHandler) GetAllGroups(c fiber.Ctx) error {
	groups, err := h.service.GetAllGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gruplar alınırken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   groups,
	})
}

func (h *GroupHandler) UpdateGroup(c fiber.Ctx) error {
	id := c.Params("id")
	updateDTO := new(dto.GroupCreateRequest)
	if err := c.Bind().WithAutoHandling().Body(updateDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if updateDTO.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name alanı zorunludur",
		})
	}
	err := h.service.UpdateGroup(id, updateDTO.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Grup güncellenirken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (h *GroupHandler) DeleteGroup(c fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteGroup(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Grup silinirken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}
