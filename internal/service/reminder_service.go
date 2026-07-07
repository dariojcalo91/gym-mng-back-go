package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/service/dto"
)

type ExpiringMembershipRepository interface {
	ListExpiringMemberships(ctx context.Context, within time.Duration) ([]*dto.MemberWithOwnerEmail, error)
}

type ReminderService struct {
	repo     ExpiringMembershipRepository
	notifier Notifier
}

func NewReminderService(repo ExpiringMembershipRepository, notifier Notifier) *ReminderService {
	return &ReminderService{repo: repo, notifier: notifier}
}

// RunDaily should be called once a day (from a cron entrypoint in cmd/, not from a gRPC handler).
func (s *ReminderService) RunDaily(ctx context.Context) error {
	expiring, err := s.repo.ListExpiringMemberships(ctx, 3*24*time.Hour)
	if err != nil {
		return err
	}
	for _, m := range expiring {
		s.notifier.Enqueue(EmailJob{
			To:      m.OwnerEmail,
			Subject: "Membresía por vencer",
			Body:    fmt.Sprintf("%s vence el %s", m.Member.Name, m.Member.MembershipEnd.Format("02/01/2006")),
		})
	}
	return nil
}
