package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"
	"engkids/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, logger *logrus.Logger) {
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info("get hi from /")
		return c.SendString("another hi")
	})

	// Инициализация сервисов
	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Настройка middleware
	middlewares.InjectAuthService(authService)

	api := app.Group("/api")

	// Публичные маршруты
	//api.Get("/logs", handlers.GetLogs(es, logger))
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Защищённые маршруты
	protected := api.Group("/user", middlewares.Protected())

	// Оставляем оригинальный обработчик профиля для обратной совместимости
	protected.Get("/profile-legacy", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		logger.WithField("userID", userID).Info("accessed protected profile route")
		return c.JSON(fiber.Map{
			"message": "Защищённый маршрут",
			"userID":  userID,
			"email":   c.Locals("email"),
		})
	})

	// Новые маршруты
	// Профиль пользователя и статистика
	protected.Get("/profile", userHandler.GetUserProfile)

	// Инвентарь и гардероб
	protected.Get("/inventory", userHandler.GetUserInventory)
	protected.Put("/inventory/item", userHandler.UpdateInventoryItem)
	protected.Post("/inventory/purchase", userHandler.PurchaseItem)

	// Словарь пользователя
	protected.Get("/words", userHandler.GetUserWords)
	protected.Post("/words/learn", userHandler.LearnWord)
}
