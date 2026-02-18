package domain

import (
	"strings"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
)

type User struct {
	ID        string
	Username  string
	Password  string // Here the Bcrypt Hash will go
	Email     string // Here the Encrypted email will go
	Role      string // SUPER_USER, ADMIN, TRAINER
	CreatedAt time.Time
}

func (m *User) Validate() error {
	if strings.TrimSpace(m.Email) == "" {
		return utils.ErrInvalidEmail
	}
	if len(strings.TrimSpace(m.Username)) < 3 {
		return utils.ErrInvalidName
	}
	return nil
}
