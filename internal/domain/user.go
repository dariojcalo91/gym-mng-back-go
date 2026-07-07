package domain

import (
	"strings"
	"time"
)

type User struct {
	ID        string
	GymID     string
	Username  string
	Password  string
	Email     string
	Role      string
	CreatedAt time.Time
}

func (m *User) Validate() error {
	if strings.TrimSpace(m.Email) == "" {
		return ErrInvalidEmail
	}
	if len(strings.TrimSpace(m.Username)) < 3 {
		return ErrInvalidName
	}
	return nil
}
