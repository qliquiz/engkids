package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"
	"engkids/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(app *fiber.App, authService *services.AuthService, userService *services.UserService, logger *logrus.Logger) {
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info("get hi from /")
		return c.SendString("another hi")
	})

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Настройка middleware
	middlewares.InjectAuthService(authService)

	api := app.Group("/api")

	// Публичные маршруты
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Все маршруты теперь публичные
	user := api.Group("/user")

	// Профиль пользователя и статистика
	user.Get("/profile", userHandler.GetUserProfile)

	// Инвентарь и гардероб
	user.Get("/inventory", userHandler.GetUserInventory)
	user.Put("/inventory/item", userHandler.UpdateInventoryItem)
	user.Post("/inventory/purchase", userHandler.PurchaseItem)

	// Словарь пользователя
	user.Get("/words", userHandler.GetUserWords)
	user.Post("/words/learn", userHandler.LearnWord)
}
