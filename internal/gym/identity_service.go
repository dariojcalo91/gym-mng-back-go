package gym

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
)

type IdentityRepository interface {
	SaveUser(ctx context.Context, u *domain.User) error
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type IdentityService struct {
	repo IdentityRepository
}

func NewIdentityService(r IdentityRepository) *IdentityService {
	return &IdentityService{repo: r}
}

func (s *IdentityService) RegisterUser(ctx context.Context, u *domain.User) error {
	// validate business rules (domain)
	if err := u.Validate(); err != nil {
		return err
	}

	// persist the user using the repository (adapter)
	err := s.repo.SaveUser(ctx, u)
	if err != nil {
		// log here original error
		return err
	}
	return nil
}

func (s *IdentityService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}
