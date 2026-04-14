package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	userDTO := new(models.UserCreateDTO)
	if err := c.Bind().WithAutoHandling().Body(userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if userDTO.Name == "" || userDTO.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name ve Email alanları zorunludur",
		})
	}
	user, err := h.service.CreateUser(userDTO.Name, userDTO.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı oluşturulurken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "not-found",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcılar alınırken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	id := c.Params("id")
	userDTO := new(models.UserUpdateDTO)
	if err := c.Bind().WithAutoHandling().Body(userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if userDTO.Name == "" || userDTO.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name ve Email alanları zorunludur",
		})
	}
	updateErr := h.service.UpdateUser(id, userDTO.Name, userDTO.Email)
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı güncellenirken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")
	deleteErr := h.service.DeleteUser(id)
	if deleteErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı silinirken bir hata oluştu",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}
