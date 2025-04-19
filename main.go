package main

import (
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"engkids/pkg/elasticsearch"
	"engkids/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title EngKids API
// @version 1.0
// @description This is an API for the EngKids project
// @host localhost:5050
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// Инициализация логгера
	appLogger, err := logger.NewLogger("engkids")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Подключение к базе данных
	db := database.ConnectDB()

	// Подключение к Elasticsearch
	var esClient *elasticsearch.Client
	es, err := elasticsearch.NewClient()
	if err != nil {
		appLogger.Warn("Failed to connect to Elasticsearch (ELK monitoring will be limited): ", err)
	} else {
		esClient = es
		appLogger.Info("Successfully connected to Elasticsearch")
	}

	// Настройка Fiber
	app := fiber.New()

	// Добавляем middleware для генерации уникального ID запроса
	app.Use(requestid.New())

	// Добавляем middleware для логирования
	app.Use(logger.LoggingMiddleware(appLogger))

	// Настройка Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Настройка маршрутов
	routes.SetupRoutes(app, db, appLogger, esClient)

	// Получаем порт из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	appLogger.WithField("port", port).Info("Starting HTTP server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal("Failed to start server: ", err)
	}
}
