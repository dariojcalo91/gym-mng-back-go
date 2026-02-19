package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       string
	Username string
	Email    string
	Role     string
}

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{secretKey, duration}
}

func (m *JWTManager) Generate(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(m.tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}
