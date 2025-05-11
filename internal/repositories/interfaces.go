package repositories

import (
	"engkids/internal/models"
	"time"
)

// UserRepository определяет операции доступа к данным пользователя
type UserRepository interface {
	GetUserByID(userID uint) (*models.User, error)

	// Операции со статистикой
	GetUserStatistics(userID uint) (*models.UserStatistics, error)
	CreateUserStatistics(stats *models.UserStatistics) error
	UpdateUserStatistics(stats *models.UserStatistics) error
	UpdateLastActive(userID uint, lastActive time.Time) error

	// Операции с инвентарем
	GetInventoryItems(userID uint) ([]models.InventoryItem, error)
	GetInventoryItem(userID uint, itemID uint) (*models.InventoryItem, error)
	UpdateInventoryItem(item *models.InventoryItem) error
	CreateInventoryItem(item *models.InventoryItem) error
	CheckInventoryItemExists(userID uint, itemID uint) (bool, error)
	UnequipItemsByCategory(userID uint, category string, exceptItemID uint) error

	// Операции с предметами
	GetItemByID(itemID uint) (*models.Item, error)

	// Операции со словами
	GetUserWords(userID uint) ([]models.UserWord, error)
	GetUserWord(userID uint, wordID uint) (*models.UserWord, error)
	CreateUserWord(userWord *models.UserWord) error
	UpdateUserWord(userWord *models.UserWord) error
	GetWordByID(wordID uint) (*models.Word, error)

	// Управление транзакциями
	BeginTx() (Transaction, error)
}

// Transaction интерфейс для транзакций базы данных
type Transaction interface {
	Commit() error
	Rollback() error
	GetUserStatistics(userID uint) (*models.UserStatistics, error)
	UpdateUserStatistics(stats *models.UserStatistics) error
	GetItemByID(itemID uint) (*models.Item, error)
	CreateInventoryItem(item *models.InventoryItem) error
	CheckInventoryItemExists(userID uint, itemID uint) (bool, error)
}
