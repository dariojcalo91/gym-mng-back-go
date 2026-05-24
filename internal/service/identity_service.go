package service

import (
	"context"
	"errors"

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

type LoginResponse struct {
	Token string
	User  domain.User
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

	encryptedEmail, err := utils.EncryptEmail(u.Email)
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

func (s *IdentityService) LoginUser(ctx context.Context, username string, password string) (LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return LoginResponse{}, errors.New("no user found whit that username")
	}

	if !(utils.CheckPasswordHash(password, user.Password)) {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	token, err := tokenGenerator(user)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func tokenGenerator(user domain.User) (string, error) {
	jwtManager := utils.NewJWTManager()
	data := utils.User{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	token, err := jwtManager.Generate(data)
	if err != nil {
		return "", err
	}

	return token, nil
}
