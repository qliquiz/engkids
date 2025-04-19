package routes

import (
	"engkids/internal/handlers"
	"engkids/internal/middlewares"
	"engkids/pkg/elasticsearch"
	"engkids/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes sets up routes for the API
func SetupRoutes(app *fiber.App, db *gorm.DB, log *logger.Logger, es *elasticsearch.Client) {
	app.Get("/", func(c *fiber.Ctx) error {
		log.WithField("path", "/").Info("Index page accessed")
		return c.SendString("hi from the server)")
	})

	api := app.Group("/api")

	// @Summary Register user
	// @Description Register a new user
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param user body handlers.RegisterRequest true "User Registration"
	// @Success 200 {object} handlers.RegisterResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Router /register [post]
	api.Post("/register", handlers.Register(db, log))

	// @Summary Login user
	// @Description Login a user to get a token
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param user body handlers.LoginRequest true "User Login"
	// @Success 200 {object} handlers.LoginResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Router /login [post]
	api.Post("/login", handlers.Login(db, log))

	// @Summary Protected route
	// @Description A route that requires authentication
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Success 200 {string} string "You are authorized!"
	// @Failure 401 {string} string "Unauthorized"
	// @Router /protected [get]
	api.Get("/protected", middlewares.Protected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user")
		log.WithFields(map[string]interface{}{
			"user_id": userID,
			"path":    "/api/protected",
		}).Info("Protected endpoint accessed")
		return c.SendString("You are authorized!")
	})

	// Если есть подключение к Elasticsearch, добавляем маршрут для поиска логов
	if es != nil {
		api.Get("/logs", middlewares.Protected(), handlers.GetLogs(es, log))
	}
}
