package services

import (
	"engkids/internal/dto"
	"engkids/internal/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

// UserService сервис для работы с пользователями
type UserService struct {
	DB *gorm.DB
}

// NewUserService создает новый сервис пользователя
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// GetUserProfile получает профиль пользователя с его статистикой
func (s *UserService) GetUserProfile(userID uint) (*dto.UserProfileResponse, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Пользователь не найден")
		}
		return nil, fiber.ErrInternalServerError
	}

	// Получаем статистику или создаем если нет
	stats, err := s.GetOrCreateUserStatistics(userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в DTO
	statsDTO := &dto.UserStatisticsDTO{
		Level:            stats.Level,
		Experience:       stats.Experience,
		Coins:            stats.Coins,
		WordsLearned:     stats.WordsLearned,
		LessonsCompleted: stats.LessonsCompleted,
		DaysStreak:       stats.DaysStreak,
		LastActive:       stats.LastActive,
	}

	// Получаем инвентарь пользователя
	var inventoryItems []models.InventoryItem
	if err := s.DB.Preload("Item").Where("user_id = ?", userID).Find(&inventoryItems).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Преобразуем инвентарь в DTO
	inventoryDTO := make([]dto.InventoryItemDTO, len(inventoryItems))
	for i, item := range inventoryItems {
		inventoryDTO[i] = dto.InventoryItemDTO{
			ID:         item.ID,
			IsEquipped: item.IsEquipped,
			AcquiredAt: item.AcquiredAt,
			Item: dto.ItemDTO{
				ID:          item.Item.ID,
				Name:        item.Item.Name,
				Type:        item.Item.Type,
				Category:    item.Item.Category,
				Rarity:      item.Item.Rarity,
				ImageURL:    item.Item.ImageURL,
				Description: item.Item.Description,
			},
		}
	}

	return &dto.UserProfileResponse{
		User:       user,
		Statistics: statsDTO,
		Inventory:  inventoryDTO,
	}, nil
}

// GetOrCreateUserStatistics получает или создает статистику пользователя
func (s *UserService) GetOrCreateUserStatistics(userID uint) (*models.UserStatistics, error) {
	var stats models.UserStatistics

	// Пробуем найти статистику
	result := s.DB.Where("user_id = ?", userID).First(&stats)
	
	// Если не найдена, создаем новую
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		stats = models.UserStatistics{
			UserID:     userID,
			Level:      1,
			Experience: 0,
			Coins:      100, // Начальный бонус
			LastActive: time.Now(),
		}
		
		if err := s.DB.Create(&stats).Error; err != nil {
			return nil, fiber.ErrInternalServerError
		}
		
		return &stats, nil
	} else if result.Error != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Обновляем последнюю активность
	stats.LastActive = time.Now()
	s.DB.Save(&stats)

	return &stats, nil
}

// UpdateUserStatistics обновляет статистику пользователя
func (s *UserService) UpdateUserStatistics(userID uint, experience, coins int) (*models.UserStatistics, error) {
	stats, err := s.GetOrCreateUserStatistics(userID)
	if err != nil {
		return nil, err
	}

	// Обновляем опыт и монеты
	stats.Experience += experience
	stats.Coins += coins

	// Проверяем, нужно ли повысить уровень
	// Формула для следующего уровня: nextLevel = currentLevel * 100
	nextLevelThreshold := stats.Level * 100
	if stats.Experience >= nextLevelThreshold {
		stats.Level++
		// Можно добавить бонус за новый уровень
		stats.Coins += stats.Level * 10
	}

	// Сохраняем изменения
	if err := s.DB.Save(stats).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return stats, nil
}

// GetUserInventory получает инвентарь пользователя
func (s *UserService) GetUserInventory(userID uint) ([]dto.InventoryItemDTO, error) {
	var inventoryItems []models.InventoryItem
	if err := s.DB.Preload("Item").Where("user_id = ?", userID).Find(&inventoryItems).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Преобразуем в DTO
	inventoryDTO := make([]dto.InventoryItemDTO, len(inventoryItems))
	for i, item := range inventoryItems {
		inventoryDTO[i] = dto.InventoryItemDTO{
			ID:         item.ID,
			IsEquipped: item.IsEquipped,
			AcquiredAt: item.AcquiredAt,
			Item: dto.ItemDTO{
				ID:          item.Item.ID,
				Name:        item.Item.Name,
				Type:        item.Item.Type,
				Category:    item.Item.Category,
				Rarity:      item.Item.Rarity,
				ImageURL:    item.Item.ImageURL,
				Description: item.Item.Description,
			},
		}
	}

	return inventoryDTO, nil
}

// UpdateInventoryItem обновляет предмет в инвентаре (например, экипирует/снимает)
func (s *UserService) UpdateInventoryItem(userID uint, req *dto.UpdateInventoryRequest) error {
	var inventoryItem models.InventoryItem
	
	// Находим предмет в инвентаре пользователя
	result := s.DB.Where("user_id = ? AND id = ?", userID, req.ItemID).First(&inventoryItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Предмет не найден в вашем инвентаре")
		}
		return fiber.ErrInternalServerError
	}

	// Если пытаемся экипировать предмет, нужно снять все предметы этой категории
	if req.IsEquipped {
		var item models.Item
		if err := s.DB.First(&item, inventoryItem.ItemID).Error; err != nil {
			return fiber.ErrInternalServerError
		}

		// Снимаем все предметы этой категории
		if err := s.DB.Model(&models.InventoryItem{}).
			Joins("JOIN items ON inventory_items.item_id = items.id").
			Where("inventory_items.user_id = ? AND items.category = ? AND items.id != ?", 
				userID, item.Category, item.ID).
			Update("is_equipped", false).Error; err != nil {
			return fiber.ErrInternalServerError
		}
	}

	// Обновляем статус экипировки
	inventoryItem.IsEquipped = req.IsEquipped
	if err := s.DB.Save(&inventoryItem).Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

// PurchaseItem покупает новый предмет для пользователя
func (s *UserService) PurchaseItem(userID uint, req *dto.PurchaseItemRequest) error {
	// Начинаем транзакцию
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Находим предмет
	var item models.Item
	if err := tx.First(&item, req.ItemID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Предмет не найден")
		}
		return fiber.ErrInternalServerError
	}

	// Проверяем, есть ли уже этот предмет у пользователя
	var count int64
	if err := tx.Model(&models.InventoryItem{}).
		Where("user_id = ? AND item_id = ?", userID, req.ItemID).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		tx.Rollback()
		return fiber.NewError(fiber.StatusConflict, "Вы уже владеете этим предметом")
	}

	// Получаем статистику для проверки монет
	var stats models.UserStatistics
	if err := tx.Where("user_id = ?", userID).First(&stats).Error; err != nil {
		tx.Rollback()
		return fiber.ErrInternalServerError
	}

	// Проверяем, достаточно ли монет
	if stats.Coins < item.Price {
		tx.Rollback()
		return fiber.NewError(fiber.StatusBadRequest, "Недостаточно монет для покупки")
	}

	// Вычитаем монеты
	stats.Coins -= item.Price
	if err := tx.Save(&stats).Error; err != nil {
		tx.Rollback()
		return fiber.ErrInternalServerError
	}

	// Добавляем предмет в инвентарь
	inventoryItem := models.InventoryItem{
		UserID:     userID,
		ItemID:     item.ID,
		IsEquipped: false,
		AcquiredAt: time.Now(),
	}

	if err := tx.Create(&inventoryItem).Error; err != nil {
		tx.Rollback()
		return fiber.ErrInternalServerError
	}

	// Подтверждаем транзакцию
	return tx.Commit().Error
}

// GetUserWords получает список слов пользователя
func (s *UserService) GetUserWords(userID uint) ([]dto.WordDTO, error) {
	var userWords []models.UserWord
	if err := s.DB.Preload("Word").Where("user_id = ?", userID).Find(&userWords).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Преобразуем в DTO
	wordDTOs := make([]dto.WordDTO, len(userWords))
	for i, userWord := range userWords {
		wordDTOs[i] = dto.WordDTO{
			ID:                 userWord.Word.ID,
			EnglishWord:        userWord.Word.EnglishWord,
			RussianTranslation: userWord.Word.RussianTranslation,
			Difficulty:         userWord.Word.Difficulty,
			Category:           userWord.Word.Category,
			KnowledgeLevel:     userWord.KnowledgeLevel,
			NextReviewAt:       &userWord.NextReviewAt,
		}
	}

	return wordDTOs, nil
}

// LearnWord добавляет слово в словарь пользователя или обновляет уровень знания
func (s *UserService) LearnWord(userID uint, req *dto.LearnWordRequest) error {
	// Проверяем существование слова
	var word models.Word
	if err := s.DB.First(&word, req.WordID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Слово не найдено")
		}
		return fiber.ErrInternalServerError
	}

	// Ищем слово в словаре пользователя
	var userWord models.UserWord
	result := s.DB.Where("user_id = ? AND word_id = ?", userID, req.WordID).First(&userWord)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Создаем новую запись слова для пользователя
		nextReview := time.Now().Add(24 * time.Hour) // Следующее повторение через 24 часа
		userWord = models.UserWord{
			UserID:         userID,
			WordID:         req.WordID,
			KnowledgeLevel: req.KnowledgeLevel,
			RepeatCount:    1,
			NextReviewAt:   nextReview,
		}

		if err := s.DB.Create(&userWord).Error; err != nil {
			return fiber.ErrInternalServerError
		}

		// Если слово новое и уровень знания выше 0, увеличиваем счетчик изученных слов
		if req.KnowledgeLevel > 0 {
			if err := s.DB.Model(&models.UserStatistics{}).
				Where("user_id = ?", userID).
				UpdateColumn("words_learned", gorm.Expr("words_learned + ?", 1)).Error; err != nil {
				return fiber.ErrInternalServerError
			}
		}
	} else if result.Error != nil {
		return fiber.ErrInternalServerError
	} else {
		// Обновляем существующую запись
		now := time.Now()
		userWord.KnowledgeLevel = req.KnowledgeLevel
		userWord.RepeatCount++
		userWord.LastReviewedAt = &now

		// Расчет интервала для следующего повторения в зависимости от уровня знания
		// Используем алгоритм подобный Spaced Repetition
		var interval time.Duration
		switch req.KnowledgeLevel {
		case 0:
			interval = 6 * time.Hour
		case 1:
			interval = 12 * time.Hour
		case 2:
			interval = 24 * time.Hour
		case 3:
			interval = 3 * 24 * time.Hour
		case 4:
			interval = 7 * 24 * time.Hour
		case 5:
			interval = 14 * 24 * time.Hour
		default:
			interval = 24 * time.Hour
		}
		userWord.NextReviewAt = now.Add(interval)

		if err := s.DB.Save(&userWord).Error; err != nil {
			return fiber.ErrInternalServerError
		}

		// Если слово достигло максимального уровня знания впервые, увеличиваем счетчик
		if req.KnowledgeLevel == 5 && userWord.KnowledgeLevel < 5 {
			if err := s.DB.Model(&models.UserStatistics{}).
				Where("user_id = ?", userID).
				UpdateColumn("words_learned", gorm.Expr("words_learned + ?", 1)).Error; err != nil {
				return fiber.ErrInternalServerError
			}
		}
	}

	return nil
}