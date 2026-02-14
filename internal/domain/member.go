package domain

import (
	"errors"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("email is required")
	ErrInvalidName  = errors.New("name must be at least 3 characters long")
)

type Member struct {
	ID     string
	Name   string
	Email  string
	Plan   string
	Status string
}

func (m *Member) Validate() error {
	if strings.TrimSpace(m.Email) == "" {
		return ErrInvalidEmail
	}
	if len(strings.TrimSpace(m.Name)) < 3 {
		return ErrInvalidName
	}
	return nil
}
