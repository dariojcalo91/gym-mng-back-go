package gym

import (
	"context"
	"os"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
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
	if err := u.Validate(); err != nil {
		return err
	}

	// Secure the password and email before saving
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		// log here original error
		return err
	}
	u.Password = hashedPassword

	encryptedEmail, err := utils.EncryptEmail(u.Email, os.Getenv("AES_MASTER_KEY"))
	if err != nil {
		// log here original error
		return err
	}
	u.Email = encryptedEmail

	// Save the user to the repository
	err = s.repo.SaveUser(ctx, u)
	if err != nil {
		// log here original error
		return err
	}
	return nil
}

func (s *IdentityService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}
