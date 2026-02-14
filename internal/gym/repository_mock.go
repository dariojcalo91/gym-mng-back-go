package gym

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveMember(ctx context.Context, member *domain.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockRepository) GetMemberByID(ctx context.Context, id string) (*domain.Member, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Member), args.Error(1)
}
