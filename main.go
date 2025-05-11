package main

import (
	_ "engkids/docs"
	"engkids/internal/repositories"
	"engkids/internal/routes"
	"engkids/internal/services"
	"engkids/pkg/database"
	"engkids/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"log"
	"os"
)

func main() {
	appLogger, err := logger.NewLogger("engkids")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Logger initialized")

	db := database.ConnectDB()

	// Создание репозиториев
	authRepo := repositories.NewAuthGormRepository(db)
	userRepo := repositories.NewUserGormRepository(db) // Предполагается, что этот метод существует или будет создан

	// Инициализация сервисов с репозиториями
	authService := services.NewAuthService(authRepo)
	//userService := services.NewUserService(userRepo)
	userService := services.NewUserService(userRepo.DB)


	app := fiber.New()
	app.Use(requestid.New())
	app.Use(logger.LoggingMiddleware(appLogger))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Правильная передача сервисов в маршруты
	routes.SetupRoutes(app, authService, userService, appLogger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	appLogger.WithField("port", port).Info("Starting HTTP server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal("Failed to start server: ", err)
	}
}