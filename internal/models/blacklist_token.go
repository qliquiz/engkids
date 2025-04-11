package models

import "time"

type BlacklistedToken struct {
	Token     string `gorm:"primaryKey"`
	ExpiresAt time.Time
}
