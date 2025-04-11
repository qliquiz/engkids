package services

import (
	"engkids/internal/models"
	"time"

	"gorm.io/gorm"
)

type TokenService struct {
	DB *gorm.DB
}

func (ts *TokenService) Blacklist(token string, exp time.Time) error {
	return ts.DB.Create(&models.BlacklistedToken{
		Token:     token,
		ExpiresAt: exp,
	}).Error
}

func (ts *TokenService) IsBlacklisted(token string) bool {
	var b models.BlacklistedToken
	err := ts.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&b).Error
	return err == nil
}
