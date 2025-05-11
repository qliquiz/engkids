package services

import (
	"engkids/internal/dto"
	"engkids/internal/interfaces"
	"engkids/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	err := db.Find(&users).Error
	return users, err
}

// UserService сервис для работы с пользователями
type UserService struct {
	UserRepo interfaces.UserRepository
}

// NewUserService создает новый сервис пользователя
func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// GetUserProfile получает профиль пользователя с его статистикой
func (s *UserService) GetUserProfile(userID uint) (*dto.UserProfileResponse, error) {
	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
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
	inventoryItems, err := s.UserRepo.GetInventoryItems(userID)
	if err != nil {
		return nil, err
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
		User:       *user,
		Statistics: statsDTO,
		Inventory:  inventoryDTO,
	}, nil
}

// GetOrCreateUserStatistics получает или создает статистику пользователя
func (s *UserService) GetOrCreateUserStatistics(userID uint) (*models.UserStatistics, error) {
	// Пробуем найти статистику
	stats, err := s.UserRepo.GetUserStatistics(userID)
	
	// Если не найдена, создаем новую
	if err != nil || stats == nil {
		newStats := &models.UserStatistics{
			UserID:     userID,
			Level:      1,
			Experience: 0,
			Coins:      100, // Начальный бонус
			LastActive: time.Now(),
		}
		
		if err := s.UserRepo.CreateUserStatistics(newStats); err != nil {
			return nil, fiber.ErrInternalServerError
		}
		
		return newStats, nil
	}

	// Обновляем последнюю активность
	stats.LastActive = time.Now()
	if err := s.UserRepo.UpdateLastActive(userID, time.Now()); err != nil {
		return nil, err
	}

	return stats, nil
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
	if err := s.UserRepo.UpdateUserStatistics(stats); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return stats, nil
}

// GetUserInventory получает инвентарь пользователя
func (s *UserService) GetUserInventory(userID uint) ([]dto.InventoryItemDTO, error) {
	inventoryItems, err := s.UserRepo.GetInventoryItems(userID)
	if err != nil {
		return nil, err
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
	// Находим предмет в инвентаре пользователя
	inventoryItem, err := s.UserRepo.GetInventoryItem(userID, req.ItemID)
	if err != nil {
		return err
	}

	// Если пытаемся экипировать предмет, нужно снять все предметы этой категории
	if req.IsEquipped {
		item, err := s.UserRepo.GetItemByID(inventoryItem.ItemID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// Снимаем все предметы этой категории
		if err := s.UserRepo.UnequipItemsByCategory(userID, item.Category, item.ID); err != nil {
			return fiber.ErrInternalServerError
		}
	}

	// Обновляем статус экипировки
	inventoryItem.IsEquipped = req.IsEquipped
	if err := s.UserRepo.UpdateInventoryItem(inventoryItem); err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}

// PurchaseItem покупает новый предмет для пользователя
func (s *UserService) PurchaseItem(userID uint, req *dto.PurchaseItemRequest) error {
	// Начинаем транзакцию
	tx, err := s.UserRepo.BeginTx()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				// Логирование ошибки откатывания транзакции
				// log.Printf("Failed to rollback tx: %v", err)
			}
		}
	}()

	// Находим предмет
	item, err := tx.GetItemByID(req.ItemID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return err
	}

	// Проверяем, есть ли уже этот предмет у пользователя
	exists, err := tx.CheckInventoryItemExists(userID, req.ItemID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.ErrInternalServerError
	}

	if exists {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.NewError(fiber.StatusConflict, "Вы уже владеете этим предметом")
	}

	// Получаем статистику для проверки монет
	stats, err := tx.GetUserStatistics(userID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.ErrInternalServerError
	}

	// Проверяем, достаточно ли монет
	if stats.Coins < item.Price {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.NewError(fiber.StatusBadRequest, "Недостаточно монет для покупки")
	}

	// Вычитаем монеты
	stats.Coins -= item.Price
	if err := tx.UpdateUserStatistics(stats); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.ErrInternalServerError
	}

	// Добавляем предмет в инвентарь
	inventoryItem := &models.InventoryItem{
		UserID:     userID,
		ItemID:     item.ID,
		IsEquipped: false,
		AcquiredAt: time.Now(),
	}

	if err := tx.CreateInventoryItem(inventoryItem); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Логирование ошибки откатывания транзакции
			// log.Printf("Failed to rollback tx: %v", rollbackErr)
		}
		return fiber.ErrInternalServerError
	}

	// Подтверждаем транзакцию
	return tx.Commit()
}

// GetUserWords получает список слов пользователя
func (s *UserService) GetUserWords(userID uint) ([]dto.WordDTO, error) {
	userWords, err := s.UserRepo.GetUserWords(userID)
	if err != nil {
		return nil, err
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
	_, err := s.UserRepo.GetWordByID(req.WordID)
	if err != nil {
		return err
	}

	// Ищем слово в словаре пользователя
	userWord, err := s.UserRepo.GetUserWord(userID, req.WordID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if userWord == nil {
		// Создаем новую запись слова для пользователя
		nextReview := time.Now().Add(24 * time.Hour) // Следующее повторение через 24 часа
		newUserWord := &models.UserWord{
			UserID:         userID,
			WordID:         req.WordID,
			KnowledgeLevel: req.KnowledgeLevel,
			RepeatCount:    1,
			NextReviewAt:   nextReview,
		}

		if err := s.UserRepo.CreateUserWord(newUserWord); err != nil {
			return fiber.ErrInternalServerError
		}

		// Если слово новое и уровень знания выше 0, увеличиваем счетчик изученных слов
		if req.KnowledgeLevel > 0 {
			stats, err := s.UserRepo.GetUserStatistics(userID)
			if err != nil {
				return fiber.ErrInternalServerError
			}
			stats.WordsLearned++
			if err := s.UserRepo.UpdateUserStatistics(stats); err != nil {
				return fiber.ErrInternalServerError
			}
		}
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

		if err := s.UserRepo.UpdateUserWord(userWord); err != nil {
			return fiber.ErrInternalServerError
		}

		// Если слово достигло максимального уровня знания впервые, увеличиваем счетчик
		if req.KnowledgeLevel == 5 && userWord.KnowledgeLevel < 5 {
			stats, err := s.UserRepo.GetUserStatistics(userID)
			if err != nil {
				return fiber.ErrInternalServerError
			}
			stats.WordsLearned++
			if err := s.UserRepo.UpdateUserStatistics(stats); err != nil {
				return fiber.ErrInternalServerError
			}
		}
	}

	return nil
}