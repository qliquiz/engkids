package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'parent'"`
}

type Child struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Age       int            `json:"age"`
	ParentID  uint           `json:"parent_id"`
	Parent    User           `json:"-" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

/*
type Progress struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ChildID     uint       `json:"child_id"`
	LessonID    uint       `json:"lesson_id"`
	Completed   bool       `json:"completed" gorm:"default:false"`
	Score       int        `json:"score" gorm:"default:0"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
*/

// UserStatistics модель для хранения статистики пользователя
type UserStatistics struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	UserID           uint           `json:"user_id" gorm:"uniqueIndex"`
	User             User           `json:"-" gorm:"foreignKey:UserID"`
	Level            int            `json:"level" gorm:"default:1"`
	Experience       int            `json:"experience" gorm:"default:0"`
	Coins            int            `json:"coins" gorm:"default:0"`
	WordsLearned     int            `json:"words_learned" gorm:"default:0"`
	LessonsCompleted int            `json:"lessons_completed" gorm:"default:0"`
	DaysStreak       int            `json:"days_streak" gorm:"default:0"`
	LastActive       time.Time      `json:"last_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// Inventory модель для хранения инвентаря пользователя
type InventoryItem struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"index"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	ItemID     uint           `json:"item_id" gorm:"index"`
	Item       Item           `json:"item" gorm:"foreignKey:ItemID"`
	IsEquipped bool           `json:"is_equipped" gorm:"default:false"`
	AcquiredAt time.Time      `json:"acquired_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// Item модель для хранения предметов (одежда, аксессуары и т.д.)
type Item struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`     // clothes, accessory, etc.
	Category    string         `json:"category" gorm:"not null"` // hat, shirt, pants, etc.
	Rarity      string         `json:"rarity" gorm:"default:'common'"`
	Price       int            `json:"price" gorm:"default:0"`
	ImageURL    string         `json:"image_url"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserWord модель для связи пользователя со словами
type UserWord struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"index"`
	User           User           `json:"-" gorm:"foreignKey:UserID"`
	WordID         uint           `json:"word_id" gorm:"index"`
	Word           Word           `json:"word" gorm:"foreignKey:WordID"`
	KnowledgeLevel int            `json:"knowledge_level" gorm:"default:0"` // 0-5 уровень знания слова
	RepeatCount    int            `json:"repeat_count" gorm:"default:0"`
	NextReviewAt   time.Time      `json:"next_review_at"`
	LastReviewedAt *time.Time     `json:"last_reviewed_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// Word модель для хранения слов
type Word struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	EnglishWord        string         `json:"english_word" gorm:"not null;uniqueIndex"`
	RussianTranslation string         `json:"russian_translation" gorm:"not null"`
	Difficulty         int            `json:"difficulty" gorm:"default:1"` // 1-5
	Category           string         `json:"category"`
	ImageURL           *string        `json:"image_url"`
	ExampleSentence    *string        `json:"example_sentence"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}
