package domain

import (
	"strings"
	"time"
)

type Gym struct {
	ID          string
	Name        string
	OwnerUserID string
	CreatedAt   time.Time
}

func (g *Gym) Validate() error {
	if len(strings.TrimSpace(g.Name)) < 3 {
		return ErrInvalidName
	}
	if strings.TrimSpace(g.OwnerUserID) == "" {
		return ErrInvalidOwner
	}
	return nil
}

func ValidateGymName(name string) error {
	if len(strings.TrimSpace(name)) < 3 {
		return ErrInvalidGymName
	}
	return nil
}
