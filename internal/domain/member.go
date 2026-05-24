package domain

import "strings"

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
	if strings.TrimSpace(m.Plan) == "" {
		return ErrInvalidPlan
	}
	return nil
}
