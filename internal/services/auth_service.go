package services

import (
	"engkids/internal/dto"
	"engkids/internal/models"
	"engkids/internal/repositories"
	"engkids/pkg/jwt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repositories.AuthGormRepository
}

func NewAuthService(repo repositories.AuthGormRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.FullAuthResponse, error) {
	// Хешируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Хеширование пароля не удалось")
	}

	// Создаем пользователя
	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Role:     "user",
	}

	// Сохраняем пользователя через репозиторий
	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	// Создаем статистику пользователя
	stats := &models.UserStatistics{
		UserID: user.ID,
		// Инициализация статистики по умолчанию
	}

	if err := s.repo.CreateUserStatistics(stats); err != nil {
		return nil, err
	}

	return s.buildFullAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.FullAuthResponse, error) {
	// Получаем пользователя через репозиторий
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Неверный email или пароль")
	}

	return s.buildFullAuthResponse(user)
}

func (s *AuthService) Refresh(oldRefresh string) (*dto.FullAuthResponse, error) {
	// Получаем токен через репозиторий
	rt, err := s.repo.GetRefreshToken(oldRefresh)
	if err != nil {
		return nil, err
	}

	// Проверяем срок действия
	if rt.ExpiresAt.Before(time.Now()) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Просроченный refresh токен")
	}

	// Получаем пользователя
	user, err := s.repo.GetUserByID(rt.UserID)
	if err != nil {
		return nil, err
	}

	// Удаляем старый токен
	if err := s.repo.DeleteRefreshToken(oldRefresh); err != nil {
		log.Println("Ошибка при удалении старого токена:", err)
	}

	return s.buildFullAuthResponse(user)
}

func (s *AuthService) buildFullAuthResponse(user *models.User) (*dto.FullAuthResponse, error) {
	// Генерируем JWT токен
	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	// Создаем refresh токен
	refreshToken := uuid.NewString()
	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	// Сохраняем через репозиторий
	if err := s.repo.SaveRefreshToken(&rt); err != nil {
		return nil, err
	}

	// Не отправляем пароль
	user.Password = ""

	return &dto.FullAuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.repo.DeleteRefreshToken(refreshToken)
}
