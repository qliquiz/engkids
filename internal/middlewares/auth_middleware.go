package middlewares

import (
	"engkids/internal/services"
	"engkids/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var authService *services.AuthService

// InjectAuthService нужен для инициализации authService
func InjectAuthService(s *services.AuthService) {
	authService = s
}

// Protected middleware с авто-обновлением токена (для mobile)
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauthorized("Необходим access токен")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return unauthorized("Неверный формат Authorization")
		}

		accessToken := parts[1]
		claims, err := jwt.ValidateToken(accessToken)

		if err == nil {
			// access валиден
			setLocals(c, claims)
			return c.Next()
		}

		// access просрочен — пробуем refresh
		refreshToken := c.Get("X-Refresh-Token")
		if refreshToken == "" {
			return unauthorized("Access истёк, refresh не передан")
		}

		if authService == nil {
			return unauthorized("AuthService не инициализирован")
		}

		resp, err := authService.Refresh(refreshToken)
		if err != nil {
			return unauthorized("Refresh токен невалиден")
		}

		// Обновляем токены клиенту
		c.Set("X-New-Access-Token", resp.AccessToken)
		c.Set("X-New-Refresh-Token", resp.RefreshToken)

		user := resp.User
		c.Locals("userID", user.ID)
		c.Locals("email", user.Email)
		c.Locals("role", user.Role)

		return c.Next()
	}
}

func setLocals(c *fiber.Ctx, claims *jwt.Claims) {
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)
}

func unauthorized(msg string) error {
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}
