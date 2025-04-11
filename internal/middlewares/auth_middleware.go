package middlewares

import (
	"engkids/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected middleware для проверки JWT токена
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		
		// Проверяем наличие заголовка
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Необходима авторизация",
			})
		}
		
		// Проверяем формат (Bearer {token})
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Неверный формат токена",
			})
		}
		
		// Извлекаем токен
		tokenString := parts[1]
		
		// Валидируем токен
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Недействительный токен",
			})
		}
		
		// Добавляем данные пользователя в контекст
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		
		return c.Next()
	}
}