package interfaces

import (
	"engkids/internal/models"
)

// AuthRepository интерфейс для работы с авторизацией
type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	SaveRefreshToken(rt *models.RefreshToken) error
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByUserID(userID uint) error
	GetUserByID(userID uint) (*models.User, error)

	// Для создания статистики
	CreateUserStatistics(stats *models.UserStatistics) error
}