package routes

import (
	//"engkids/internal/handlers"
	"engkids/internal/middlewares"
	//"engkids/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, logger *logrus.Logger) {
	//authService := services.NewAuthService(db)

	//authHandler := handlers.NewAuthHandler(authService)

	//middlewares.InjectAuthService(authService)

	api := app.Group("/api")

	// Публичные маршруты
	/*auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)*/

	// Защищённые маршруты
	protected := api.Group("/user", middlewares.Protected())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("userID")
		logger.WithField("userID", userID).Info("accessed protected profile route")
		return c.JSON(fiber.Map{
			"message": "Защищённый маршрут",
			"userID":  userID,
		})
	})
}
