package main

import (
	"encoding/json"
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

type LogMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

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
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))

	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	msg := LogMessage{
		Level:   "INFO",
		Message: "Hello from Go!",
	}
	_ = json.NewEncoder(conn).Encode(msg)
}
