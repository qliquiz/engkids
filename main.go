package main

import (
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		log.Fatal("Error loading .env file")
	}
	db := database.ConnectDB()
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}
	log.Fatal(app.Listen(":" + port))
}
