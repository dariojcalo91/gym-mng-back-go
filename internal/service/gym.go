package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
)

type Repository interface {
	SaveMember(ctx context.Context, member *domain.Member) error
	GetMemberByID(ctx context.Context, id string) (*domain.Member, error)
}

type GymService struct {
	repo         Repository
	emailChannel chan string
	wg           sync.WaitGroup
}

func NewService(r Repository) *GymService {
	return &GymService{
		repo:         r,
		emailChannel: make(chan string, 5),
	}
}

func (s *GymService) StartWorkers(ctx context.Context) {
	s.wg.Add(3)
	for i := range 3 {
		go s.emailWorker(ctx, i+1)
	}
}

func (s *GymService) emailWorker(ctx context.Context, id int) {
	defer s.wg.Done()

	slog.Info("worker ready to process emails", "worker_id", id)
	for {
		select {
		case email, ok := <-s.emailChannel:
			if !ok {
				return
			}
			if err := s.processEmail(ctx, email, id); err != nil {
				return
			}
		case <-ctx.Done():
			slog.Info("worker shutting down", "worker_id", id, "reason", ctx.Err())
			return
		}
	}
}

func (s *GymService) processEmail(ctx context.Context, email string, workerID int) error {
	select {
	case <-time.After(2 * time.Second):
		slog.Info("email sent", "worker_id", workerID, "to", email)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *GymService) RegisterMember(ctx context.Context, m *domain.Member) error {
	if err := m.Validate(); err != nil {
		return err
	}

	m.Status = "active"

	err := s.repo.SaveMember(ctx, m)
	if err != nil {
		return err
	}

	select {
	case s.emailChannel <- m.Email:
	default:
		slog.Warn("email channel full, dropping email", "email", m.Email)
	}

	return nil
}

func (s *GymService) GetMemberStatus(ctx context.Context, id string) (*domain.Member, error) {
	return s.repo.GetMemberByID(ctx, id)
}

func (s *GymService) Shutdown() {
	slog.Info("closing email channel and waiting for workers")
	close(s.emailChannel)
	s.wg.Wait()
	slog.Info("all workers completed")
}
