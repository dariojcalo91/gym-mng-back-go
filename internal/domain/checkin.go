package domain

import (
	"strings"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusNotApplicable PaymentStatus = "NOT_APPLICABLE"
	PaymentStatusPaid          PaymentStatus = "PAID"
	PaymentStatusPending       PaymentStatus = "ON_CREDIT"
)

type CheckIn struct {
	ID            string
	GymID         string
	MemberID      string
	CheckedInAt   time.Time
	PaymentStatus PaymentStatus
	PaymentNote   string
	CreatedAt     time.Time
}

func (c *CheckIn) Validate() error {
	if strings.TrimSpace(c.MemberID) == "" {
		return ErrInvalidMember
	}
	if strings.TrimSpace(c.GymID) == "" {
		return ErrInvalidGym
	}
	switch c.PaymentStatus {
	case PaymentStatusNotApplicable, PaymentStatusPaid, PaymentStatusPending:
	default:
		return ErrInvalidPaymentStatus
	}
	return nil
}
