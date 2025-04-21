package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"
	"engkids/internal/services"
	"engkids/pkg/elasticsearch"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, es *elasticsearch.Client, logger *logrus.Logger) {
	app.Get("/", func(c *fiber.Ctx) error {
		logger.Info("get hi from /")
		return c.SendString("hi from the server!")
	})

	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService)
	middlewares.InjectAuthService(authService)

	api := app.Group("/api")

	// Публичные маршруты
	api.Get("/logs", handlers.GetLogs(es, logger))

	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Защищённые маршруты
	protected := api.Group("/user", middlewares.Protected())

	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		logger.WithField("userID", userID).Info("accessed protected profile route")
		return c.JSON(fiber.Map{
			"message": "Защищённый маршрут",
			"userID":  userID,
			"email":   c.Locals("email"),
		})
	})
}
