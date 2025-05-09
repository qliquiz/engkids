package dto

import (
	"engkids/internal/models"
	"time"
)

// UserProfileResponse - полный ответ с данными профиля пользователя
type UserProfileResponse struct {
	User       models.User        `json:"user"`
	Statistics *UserStatisticsDTO `json:"statistics"`
	Inventory  []InventoryItemDTO `json:"inventory,omitempty"`
}

// UserStatisticsDTO - DTO для статистики пользователя
type UserStatisticsDTO struct {
	Level            int       `json:"level"`
	Experience       int       `json:"experience"`
	Coins            int       `json:"coins"`
	WordsLearned     int       `json:"words_learned"`
	LessonsCompleted int       `json:"lessons_completed"`
	DaysStreak       int       `json:"days_streak"`
	LastActive       time.Time `json:"last_active"`
}

// InventoryItemDTO - DTO для предметов инвентаря
type InventoryItemDTO struct {
	ID         uint      `json:"id"`
	Item       ItemDTO   `json:"item"`
	IsEquipped bool      `json:"is_equipped"`
	AcquiredAt time.Time `json:"acquired_at"`
}

// ItemDTO - DTO для предмета
type ItemDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Rarity      string `json:"rarity"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

// WordDTO - DTO для слова
type WordDTO struct {
	ID                 uint       `json:"id"`
	EnglishWord        string     `json:"english_word"`
	RussianTranslation string     `json:"russian_translation"`
	Difficulty         int        `json:"difficulty"`
	Category           string     `json:"category"`
	KnowledgeLevel     int        `json:"knowledge_level,omitempty"`
	NextReviewAt       *time.Time `json:"next_review_at,omitempty"`
}

// UpdateInventoryRequest - запрос на обновление инвентаря
type UpdateInventoryRequest struct {
	ItemID     uint `json:"item_id" validate:"required"`
	IsEquipped bool `json:"is_equipped"`
}

// PurchaseItemRequest - запрос на покупку предмета
type PurchaseItemRequest struct {
	ItemID uint `json:"item_id" validate:"required"`
}

// LearnWordRequest - запрос на изучение слова
type LearnWordRequest struct {
	WordID         uint `json:"word_id" validate:"required"`
	KnowledgeLevel int  `json:"knowledge_level" validate:"min=0,max=5"`
}
