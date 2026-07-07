package service

import (
	"context"
	"errors"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
)

type IdentityRepository interface {
	// RegisterUserWithGym persists the user and its gym in a single transaction.
	RegisterUserWithGym(ctx context.Context, u *domain.User, gymName string) (*domain.Gym, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type IdentityService struct {
	repo IdentityRepository
}

type LoginResponse struct {
	Token string
	User  domain.User
}

type RegisterResult struct {
	UserID string
	GymID  string
}

func NewIdentityService(r IdentityRepository) *IdentityService {
	return &IdentityService{repo: r}
}

func (s *IdentityService) RegisterUser(ctx context.Context, u *domain.User, gymName string) (*RegisterResult, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := domain.ValidateGymName(gymName); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashedPassword

	encryptedEmail, err := utils.EncryptEmail(u.Email)
	if err != nil {
		return nil, err
	}
	u.Email = encryptedEmail

	gym, err := s.repo.RegisterUserWithGym(ctx, u, gymName)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{UserID: u.ID, GymID: gym.ID}, nil
}

func (s *IdentityService) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *IdentityService) LoginUser(ctx context.Context, username, password string) (LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return LoginResponse{}, errors.New("no user found with that username")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return LoginResponse{}, errors.New("invalid credentials")
	}
	token, err := tokenGenerator(user)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: token, User: user}, nil
}

func tokenGenerator(user domain.User) (string, error) {
	jwtManager := utils.NewJWTManager()
	data := utils.User{ID: user.ID, Username: user.Username, Role: user.Role}
	// GymID must travel in the JWT claims too — the interceptor reads it from
	// here on every subsequent request, it's never trusted from the request body.
	token, err := jwtManager.Generate(data)
	if err != nil {
		return "", err
	}
	return token, nil
}
