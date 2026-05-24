package service

import (
	"context"
	"log"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
)

// The interface (Port) lives where it is used
type Repository interface {
	SaveMember(ctx context.Context, member *domain.Member) error
}

type GymService struct {
	repo         Repository
	emailChannel chan string // Channel for email addresses
}

func NewService(r Repository) *GymService {
	s := &GymService{
		repo:         r,
		emailChannel: make(chan string, 5), // Buffered channel with capacity of 5
	}

	// Initialize workers to process emails from the channel
	// For example, 3 workers processing emails in parallel
	for i := 1; i <= 3; i++ {
		go s.emailWorker(i)
	}

	return s
}

// The Worker: Listens to the channel indefinitely
func (s *GymService) emailWorker(id int) {
	log.Printf("Worker %d ready to process emails...", id)
	for email := range s.emailChannel {
		// Heavy task simulation
		time.Sleep(2 * time.Second)
		log.Printf("Worker %d sent email to: %s", id, email)
	}
}

func (s *GymService) RegisterMember(ctx context.Context, m *domain.Member) error {
	// validate bussines rule (domain)
	if err := m.Validate(); err != nil {
		return err
	}

	// define inicial status for new members
	m.Status = "active"

	// persist the member using the repository (adapter) (synchronous/blocking operation)
	err := s.repo.SaveMember(ctx, m)
	if err != nil {
		// log here original error
		return err
	}

	// Send to channel (Asynchronous)
	// It does NOT block the user and the worker will take it when it's free
	s.emailChannel <- m.Email

	return nil
}

func (s *GymService) Shutdown() {
	log.Println("Closing email channel...")
	close(s.emailChannel) // This notifies the workers that the channel is done
}
