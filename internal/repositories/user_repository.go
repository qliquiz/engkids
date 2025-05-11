package repositories

import (
	"engkids/internal/models"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserGormRepository реализует UserRepository с использованием GORM
type UserGormRepository struct {
	DB *gorm.DB
}

// NewUserGormRepository создает новый экземпляр репозитория
func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{DB: db}
}

// GetUserByID получает пользователя по ID
func (r *UserGormRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Пользователь не найден")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &user, nil
}

// GetUserStatistics получает статистику пользователя по ID
func (r *UserGormRepository) GetUserStatistics(userID uint) (*models.UserStatistics, error) {
	var stats models.UserStatistics
	err := r.DB.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Возвращаем nil, чтобы показать, что статистика не существует
		}
		return nil, fiber.ErrInternalServerError
	}
	return &stats, nil
}

// CreateUserStatistics создает новую статистику пользователя
func (r *UserGormRepository) CreateUserStatistics(stats *models.UserStatistics) error {
	return r.DB.Create(stats).Error
}

// UpdateUserStatistics обновляет статистику пользователя
func (r *UserGormRepository) UpdateUserStatistics(stats *models.UserStatistics) error {
	return r.DB.Save(stats).Error
}

// UpdateLastActive обновляет время последней активности
func (r *UserGormRepository) UpdateLastActive(userID uint, lastActive time.Time) error {
	return r.DB.Model(&models.UserStatistics{}).
		Where("user_id = ?", userID).
		Update("last_active", lastActive).Error
}

// GetInventoryItems получает все предметы инвентаря пользователя
func (r *UserGormRepository) GetInventoryItems(userID uint) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	err := r.DB.Preload("Item").Where("user_id = ?", userID).Find(&items).Error
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	return items, nil
}

// GetInventoryItem получает конкретный предмет инвентаря
func (r *UserGormRepository) GetInventoryItem(userID uint, itemID uint) (*models.InventoryItem, error) {
	var item models.InventoryItem
	err := r.DB.Where("user_id = ? AND id = ?", userID, itemID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Предмет не найден в вашем инвентаре")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &item, nil
}

// UpdateInventoryItem обновляет предмет инвентаря
func (r *UserGormRepository) UpdateInventoryItem(item *models.InventoryItem) error {
	return r.DB.Save(item).Error
}

// CreateInventoryItem создает новый предмет инвентаря
func (r *UserGormRepository) CreateInventoryItem(item *models.InventoryItem) error {
	return r.DB.Create(item).Error
}

// CheckInventoryItemExists проверяет, есть ли у пользователя уже предмет
func (r *UserGormRepository) CheckInventoryItemExists(userID uint, itemID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&models.InventoryItem{}).
		Where("user_id = ? AND item_id = ?", userID, itemID).
		Count(&count).Error
	if err != nil {
		return false, fiber.ErrInternalServerError
	}
	return count > 0, nil
}

// UnequipItemsByCategory снимает экипировку со всех предметов категории, кроме одного
func (r *UserGormRepository) UnequipItemsByCategory(userID uint, category string, exceptItemID uint) error {
	return r.DB.Model(&models.InventoryItem{}).
		Joins("JOIN items ON inventory_items.item_id = items.id").
		Where("inventory_items.user_id = ? AND items.category = ? AND items.id != ?",
			userID, category, exceptItemID).
		Update("is_equipped", false).Error
}

// GetItemByID получает предмет по ID
func (r *UserGormRepository) GetItemByID(itemID uint) (*models.Item, error) {
	var item models.Item
	err := r.DB.First(&item, itemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Предмет не найден")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &item, nil
}

// GetUserWords получает все слова пользователя
func (r *UserGormRepository) GetUserWords(userID uint) ([]models.UserWord, error) {
	var words []models.UserWord
	err := r.DB.Preload("Word").Where("user_id = ?", userID).Find(&words).Error
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	return words, nil
}

// GetUserWord получает конкретное слово пользователя
func (r *UserGormRepository) GetUserWord(userID uint, wordID uint) (*models.UserWord, error) {
	var word models.UserWord
	err := r.DB.Where("user_id = ? AND word_id = ?", userID, wordID).First(&word).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Возвращаем nil, чтобы показать, что слово не найдено
		}
		return nil, fiber.ErrInternalServerError
	}
	return &word, nil
}

// CreateUserWord создает новое слово пользователя
func (r *UserGormRepository) CreateUserWord(word *models.UserWord) error {
	return r.DB.Create(word).Error
}

// UpdateUserWord обновляет слово пользователя
func (r *UserGormRepository) UpdateUserWord(word *models.UserWord) error {
	return r.DB.Save(word).Error
}

// GetWordByID получает слово по ID
func (r *UserGormRepository) GetWordByID(wordID uint) (*models.Word, error) {
	var word models.Word
	err := r.DB.First(&word, wordID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Слово не найдено")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &word, nil
}

// BeginTx начинает новую транзакцию
func (r *UserGormRepository) BeginTx() (Transaction, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTransaction{tx: tx}, nil
}

// GormTransaction реализует интерфейс Transaction
type GormTransaction struct {
	tx *gorm.DB
}

func (t *GormTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *GormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *GormTransaction) GetUserStatistics(userID uint) (*models.UserStatistics, error) {
	var stats models.UserStatistics
	err := t.tx.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (t *GormTransaction) UpdateUserStatistics(stats *models.UserStatistics) error {
	return t.tx.Save(stats).Error
}

func (t *GormTransaction) GetItemByID(itemID uint) (*models.Item, error) {
	var item models.Item
	err := t.tx.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (t *GormTransaction) CreateInventoryItem(item *models.InventoryItem) error {
	return t.tx.Create(item).Error
}

func (t *GormTransaction) CheckInventoryItemExists(userID uint, itemID uint) (bool, error) {
	var count int64
	err := t.tx.Model(&models.InventoryItem{}).
		Where("user_id = ? AND item_id = ?", userID, itemID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


