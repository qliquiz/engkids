package main

import (
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title EngKids API
// @version 1.0
// @description This is an API for the EngKids project
// @host localhost:3030
// @BasePath /
// В main.go
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к базе данных
	db := database.ConnectDB()

	// Инициализация приложения Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Добавление обработчика для корневого маршрута
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to EngKids API!")
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Настройка маршрутов
	routes.SetupRoutes(app, db)

	// Чтение порта из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
	}

	// Запуск сервера
	log.Fatal(app.Listen(":" + port))
}
