package jwt

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(id uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func VerifyToken(t string) bool {
	_, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	return err == nil
}
