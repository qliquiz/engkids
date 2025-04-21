package main

import (
	_ "engkids/docs"
	"engkids/internal/routes"
	"engkids/pkg/database"
	"engkids/pkg/elasticsearch"
	"engkids/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"log"
	"os"
)

func main() {
	es, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	appLogger, err := logger.NewLogger("engkids")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Logger initialized")

	db := database.ConnectDB()

	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.LoggingMiddleware(appLogger))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetupRoutes(app, db, es, appLogger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	appLogger.WithField("port", port).Info("Starting HTTP server")
	if err := app.Listen(":" + port); err != nil {
		appLogger.Fatal("Failed to start server: ", err)
	}
}
