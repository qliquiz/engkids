package database

import (
	"engkids/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB открывает подключение к базе данных
func ConnectDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbname == "" {
		log.Fatalf("Missing one or more required DB environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Remove models.Progress from the migration since it's commented out
	err = db.AutoMigrate(
		&models.User{},
		&models.Child{},
		// &models.Progress{},  // Removed
		&models.RefreshToken{},
		&models.UserStatistics{},
		&models.InventoryItem{},
		&models.Item{},
		&models.UserWord{},
		&models.Word{},
	)
	if err != nil {
		log.Fatal("Error during migration: ", err)
	}

	fmt.Println("Migrations applied successfully")
	return db
}

// DB глобальная переменная для хранения подключения к базе данных
var DB *gorm.DB

// InitDB инициализирует подключение и миграцию базы данных
func InitDB() {
	DB = ConnectDB()
}

// CloseDB закрывает подключение к базе данных
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to close database connection: ", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Failed to close database connection: ", err)
	}

	fmt.Println("Database connection closed")
}
