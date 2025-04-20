package main

import (
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"engkids/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"log"
	"os"
)

func main() {
	// Инициализация логгера (пишет в файл logs/app.log)
	appLogger, err := logger.NewLogger("engkids")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Logger initialized")

	// Подключение к базе данных
	db := database.ConnectDB()

	// Настройка Fiber
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.LoggingMiddleware(appLogger))

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Настройка маршрутов
	routes.SetupRoutes(app, db, appLogger)

	// Получаем порт из .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appLogger.WithField("port", port).Info("Starting HTTP server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal("Failed to start server: ", err)
	}
}
