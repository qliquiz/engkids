package repositories

import (
	"engkids/internal/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
}

// AuthGormRepository реализация AuthRepository с использованием GORM
type AuthGormRepository struct {
	DB *gorm.DB
}

// NewAuthGormRepository создает новый экземпляр репозитория авторизации
func NewAuthGormRepository(db *gorm.DB) *AuthGormRepository {
	return &AuthGormRepository{DB: db}
}

// GetUserByEmail получает пользователя по email
func (r *AuthGormRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный email или пароль")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &user, nil
}

// CreateUser создает нового пользователя
func (r *AuthGormRepository) CreateUser(user *models.User) error {
	var existing models.User
	if err := r.DB.Where("email = ?", user.Email).First(&existing).Error; err == nil {
		return fiber.NewError(fiber.StatusConflict, "Пользователь уже существует")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.ErrInternalServerError
	}

	if err := r.DB.Create(user).Error; err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

// GetRefreshToken получает refresh токен
func (r *AuthGormRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.DB.Where("token = ?", token).First(&rt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный refresh токен")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &rt, nil
}

// SaveRefreshToken сохраняет refresh токен
func (r *AuthGormRepository) SaveRefreshToken(rt *models.RefreshToken) error {
	err := r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).Create(rt).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return nil
}

// DeleteRefreshToken удаляет refresh токен
func (r *AuthGormRepository) DeleteRefreshToken(token string) error {
	result := r.DB.Where("token = ?", token).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return fiber.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "Неверный refresh токен")
	}
	return nil
}

// DeleteRefreshTokenByUserID удаляет все refresh токены пользователя
func (r *AuthGormRepository) DeleteRefreshTokenByUserID(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

// GetUserByID получает пользователя по ID
func (r *AuthGormRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Пользователь не найден")
		}
		return nil, fiber.ErrInternalServerError
	}
	return &user, nil
}