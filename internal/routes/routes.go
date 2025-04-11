package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes настраивает все маршруты приложения
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(db)

	// Группа API
	api := app.Group("/api")

	// Маршруты аутентификации (публичные)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Защищенные маршруты
	protected := api.Group("/user", middlewares.Protected())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		return c.JSON(fiber.Map{
			"message": "Защищенный маршрут",
			"userID":  userID,
		})
	})

	// Здесь можно добавить другие маршруты
}