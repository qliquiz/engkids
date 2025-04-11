package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"engkids/internal/models" // Импортируем модели
)

// ConnectDB устанавливает подключение к базе данных
func ConnectDB() *gorm.DB {
	// Получаем строку подключения из переменной окружения
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN environment variable is not set")
	}

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Печатаем сообщение об успешном подключении
	fmt.Println("Database connection established")

	// Выполняем миграции для всех моделей
	err = db.AutoMigrate(&models.User{}, &models.Child{}, &models.Progress{})
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
