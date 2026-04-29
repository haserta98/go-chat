package handler

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/haserta98/go-rest/internal"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/service"
)

type UserHandler struct {
	service     service.UserService
	redisClient *internal.RedisClient
}

func NewUserHandler(service service.UserService, redisClient *internal.RedisClient) *UserHandler {
	return &UserHandler{service: service, redisClient: redisClient}
}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	userDTO := new(models.UserCreateDTO)
	if err := c.Bind().WithAutoHandling().Body(userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if userDTO.Name == "" || userDTO.Email == "" || userDTO.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, Email ve Password alanları zorunludur",
		})
	}
	user, err := h.service.CreateUser(userDTO.Name, userDTO.Email, userDTO.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) LoginUser(c fiber.Ctx) error {
	loginDTO := new(models.UserLoginDTO)
	if err := c.Bind().WithAutoHandling().Body(loginDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}

	token, user, err := h.service.LoginUser(loginDTO.Name, loginDTO.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create session UUID
	sessionID := uuid.New().String()
	sessionData, _ := json.Marshal(map[string]string{
		"user_id": user.ID,
		"name":    user.Name,
		// We can also store the token if needed, but not strictly necessary since the session is now the source of truth
		"token": token,
	})

	h.redisClient.Set("session:"+sessionID, sessionData, 24*time.Hour)

	// Set HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  token, // Also keep the token in response for backward compatibility
		"data":   user,
	})
}

func (h *UserHandler) LogoutUser(c fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		// Delete session from Redis
		h.redisClient.Del("session:" + sessionID)
	}

	// Clear the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"status": "success",
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

func (h *UserHandler) AddContact(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	contactDTO := new(models.UserContactDTO)
	if err := c.Bind().WithAutoHandling().Body(contactDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if contactDTO.ContactID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ContactID alanı zorunludur",
		})
	}

	err := h.service.AddContact(userID.(string), contactDTO.ContactID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (h *UserHandler) RemoveContact(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	contactDTO := new(models.UserContactDTO)
	if err := c.Bind().WithAutoHandling().Body(contactDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz veri",
		})
	}
	if contactDTO.ContactID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ContactID alanı zorunludur",
		})
	}

	err := h.service.RemoveContact(userID.(string), contactDTO.ContactID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (h *UserHandler) GetContacts(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	users, err := h.service.GetContacts(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kişiler alınırken hata oluştu",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func (h *UserHandler) GetMe(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (h *UserHandler) GetContactsOnlineStatus(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	contacts, err := h.service.GetContacts(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kişiler alınırken hata oluştu",
		})
	}

	onlineStatuses := make(map[string]bool)
	for _, contact := range contacts {
		isOnline, _ := h.redisClient.Exists("user:online:" + contact.ID)
		onlineStatuses[contact.ID] = isOnline
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   onlineStatuses,
	})
}

