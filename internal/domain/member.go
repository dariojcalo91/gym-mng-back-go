package domain

import (
	"strings"
	"time"
)

type MemberType string

const (
	MemberTypeMonthly    MemberType = "MONTHLY"
	MemberTypeOccasional MemberType = "OCCASIONAL"
)

type Member struct {
	ID              string
	GymID           string
	Name            string
	Phone           string
	Type            MemberType
	MembershipStart *time.Time
	MembershipEnd   *time.Time
	CreatedAt       time.Time
}

func (m *Member) Validate() error {
	if len(strings.TrimSpace(m.Name)) < 3 {
		return ErrInvalidName
	}
	if strings.TrimSpace(m.Phone) == "" {
		return ErrInvalidPhone
	}
	if strings.TrimSpace(m.GymID) == "" {
		return ErrInvalidGym
	}
	switch m.Type {
	case MemberTypeMonthly, MemberTypeOccasional:
	default:
		return ErrInvalidMemberType
	}
	if m.Type == MemberTypeMonthly && (m.MembershipStart == nil || m.MembershipEnd == nil) {
		return ErrInvalidMembershipDates
	}
	return nil
}

// IsMembershipActive reports whether a monthly member's plan currently covers "now".
// Used by the check-in flow to decide if payment should be asked.
func (m *Member) IsMembershipActive(now time.Time) bool {
	if m.Type != MemberTypeMonthly || m.MembershipEnd == nil {
		return false
	}
	return !now.After(*m.MembershipEnd)
}
