package handlers

import (
	"engkids/internal/models"
	"engkids/pkg/jwt"
	"engkids/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

// Register создает нового пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Невозможно обработать данные",
		})
	}

	// Валидация данных
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Проверка существования пользователя
	var existingUser models.User
	if result := h.DB.Where("email = ?", req.Email).First(&existingUser); result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Пользователь с таким email уже существует",
		})
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при обработке пароля",
		})
	}

	// Создание пользователя
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "parent",
	}

	result := h.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании пользователя",
		})
	}

	// Генерация JWT токена
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании токена",
		})
	}

	// Не возвращаем пароль
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login аутентифицирует пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Невозможно обработать данные",
		})
	}

	// Валидация данных
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Поиск пользователя
	var user models.User
	if result := h.DB.Where("email = ?", req.Email).First(&user); result.RowsAffected == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный email или пароль",
		})
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Неверный email или пароль",
		})
	}

	// Генерация JWT токена
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании токена",
		})
	}

	// Не возвращаем пароль
	user.Password = ""

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		Token: token,
		User:  user,
	})
}