package services

import (
	"engkids/internal/dto"
	"engkids/internal/models"
	"engkids/pkg/jwt"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.FullAuthResponse, error) {
	var existing models.User
	if err := s.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Пользователь уже существует")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("DB error:", err)
		return nil, fiber.ErrInternalServerError
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Хеширование пароля не удалось")
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     "user",
	}

	if err := s.DB.Create(&user).Error; err != nil {
		log.Println("DB create error:", err)
		return nil, fiber.ErrInternalServerError
	}

	return s.buildFullAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.FullAuthResponse, error) {
	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный email или пароль")
		}
		log.Println("DB error:", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный email или пароль")
	}

	return s.buildFullAuthResponse(&user)
}

func (s *AuthService) Refresh(oldRefresh string) (*dto.FullAuthResponse, error) {
	var rt models.RefreshToken
	err := s.DB.Where("token = ?", oldRefresh).First(&rt).Error
	if err != nil || rt.ExpiresAt.Before(time.Now()) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный или просроченный refresh токен")
	}

	var user models.User
	if err := s.DB.First(&user, rt.UserID).Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}
	log.Println("refresh: ", &rt)
	s.DB.Delete(&rt)
	return s.buildFullAuthResponse(&user)
}

func (s *AuthService) buildFullAuthResponse(user *models.User) (*dto.FullAuthResponse, error) {
	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	refreshToken := uuid.NewString()

	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	err = s.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).Create(&rt).Error
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	user.Password = ""
	return &dto.FullAuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	result := s.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return fiber.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}
	return nil
}
