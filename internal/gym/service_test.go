package gym

import (
	"context"
	"errors"
	"testing"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterMember(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: Member registered", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		member := &domain.Member{Name: "Alex", Email: "alex@mail.com", Plan: "Gold"}

		// mockRepo response successful save
		mockRepo.On("SaveMember", ctx, member).Return(nil)

		err := service.RegisterMember(ctx, member)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: name too short", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		member := &domain.Member{Name: "Al", Email: "alex@mail.com", Plan: "Gold"}

		err := service.RegisterMember(ctx, member)

		assert.Error(t, err)
		// We use domine rules to check the error type, not the exact message
		assert.Equal(t, domain.ErrInvalidName, err)
		// verify that SaveMember was never called due to validation failure
		mockRepo.AssertNotCalled(t, "SaveMember", mock.Anything, mock.Anything)
	})

	t.Run("Error: Database error - db connection lost", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		member := &domain.Member{Name: "Alex", Email: "alex@mail.com", Plan: "Gold"}
		dbError := errors.New("db connection lost")

		// Simulamos que la DB falla
		mockRepo.On("SaveMember", ctx, member).Return(dbError)

		err := service.RegisterMember(ctx, member)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db connection lost")
		mockRepo.AssertExpectations(t)
	})
}
