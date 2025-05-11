package interfaces

import (
	"engkids/internal/models"
)

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