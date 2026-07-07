package domain

import "errors"

var (
	ErrInvalidEmail           = errors.New("email is required")
	ErrInvalidName            = errors.New("name must be at least 3 characters long")
	ErrInvalidOwner           = errors.New("owner user id is required")
	ErrInvalidPhone           = errors.New("phone is required")
	ErrInvalidGym             = errors.New("gym id is required")
	ErrInvalidMember          = errors.New("member id is required")
	ErrInvalidMemberType      = errors.New("member type must be monthly or occasional")
	ErrInvalidMembershipDates = errors.New("monthly members require a membership start and end date")
	ErrInvalidPaymentStatus   = errors.New("invalid payment status")
	ErrMemberNotFound         = errors.New("member not found")
	ErrInvalidGymName         = errors.New("gym name must be at least 3 characters long")
)
