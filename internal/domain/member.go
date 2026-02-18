package domain

import (
	"strings"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
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
		return utils.ErrInvalidEmail
	}
	if len(strings.TrimSpace(m.Name)) < 3 {
		return utils.ErrInvalidName
	}
	return nil
}
