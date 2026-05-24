package service

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}
type MockUserRepo struct {
	mock.Mock
}

func (m *MockRepository) SaveMember(ctx context.Context, member *domain.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockRepository) GetMemberByID(ctx context.Context, id string) (*domain.Member, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Member), args.Error(1)
}

func (m *MockUserRepo) SaveUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}
