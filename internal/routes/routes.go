package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"
	"engkids/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes настраивает все маршруты приложения
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Создаем сервис аутентификации
	authService := services.NewAuthService(db)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)

	// Передаем authService в middleware
	middlewares.InjectAuthService(authService)

	// Группа API
	api := app.Group("/api")

	// Маршруты аутентификации (публичные)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)

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
