package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"uniqueIndex"`
	Token     string    `gorm:"unique"`
	ExpiresAt time.Time `gorm:"not null"`
}
