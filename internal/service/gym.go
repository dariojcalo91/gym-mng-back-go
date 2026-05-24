package service

import (
	"context"
	"log"
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

	log.Printf("Worker %d ready to process emails...", id)
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
			log.Printf("Worker %d shutting down: %v", id, ctx.Err())
			return
		}
	}
}

func (s *GymService) processEmail(ctx context.Context, email string, workerID int) error {
	select {
	case <-time.After(2 * time.Second):
		log.Printf("Worker %d sent email to: %s", workerID, email)
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
		log.Println("Warning: email channel full, dropping email for", m.Email)
	}

	return nil
}

func (s *GymService) GetMemberStatus(ctx context.Context, id string) (*domain.Member, error) {
	return s.repo.GetMemberByID(ctx, id)
}

func (s *GymService) Shutdown() {
	log.Println("Closing email channel...")
	close(s.emailChannel)
	s.wg.Wait()
	log.Println("All workers completed.")
}
