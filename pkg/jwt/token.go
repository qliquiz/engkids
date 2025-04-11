package jwt

import (
	"engkids/config"
	"time"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims // Используем StandardClaims для работы с зарегистрированными полями
}

// GenerateToken создает новый JWT токен
func GenerateToken(userID uint, email, role string) (string, error) {
	// Получаем секретный ключ из конфигурации
	secretKey := config.GetEnv("JWT_SECRET_KEY", "your-secret-key")
	
	
	// Создаем claims с данными пользователя
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	
	// Создаем токен с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) (*Claims, error) {
	secretKey := config.GetEnv("JWT_SECRET_KEY", "your-secret-key")
	
	claims := &Claims{}
	
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	
	return claims, nil
}
