package middleware

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/haserta98/go-rest/internal"
	"github.com/redis/go-redis/v9"
)

func NewAuthMiddleware(redisClient *internal.RedisClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Read session from cookie first
		sessionID := c.Cookies("session_id")

		// Fallback for websocket connections if they pass it via query?
		// Websockets will automatically send cookies if credentials: include is used.
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Oturum bulunamadı"})
		}

		// Look up session in Redis
		sessionData, err := redisClient.Get("session:" + sessionID)
		if err == redis.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz veya süresi dolmuş oturum"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Oturum doğrulanırken hata oluştu"})
		}

		// sessionData contains a JSON string of {"user_id": "...", "name": "..."}
		var claims map[string]interface{}
		if err := json.Unmarshal([]byte(sessionData), &claims); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Oturum içeriği okunamadı"})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("name", claims["name"])

		return c.Next()
	}
}
