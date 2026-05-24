package domain

import "errors"

var (
	ErrInvalidEmail = errors.New("email is required")
	ErrInvalidName  = errors.New("name must be at least 3 characters long")
	ErrInvalidPlan  = errors.New("plan is required")
)
