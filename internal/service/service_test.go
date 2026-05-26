package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("AES_MASTER_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	_ = os.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	os.Exit(m.Run())
}

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

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: User registered", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		user := &domain.User{Username: "admin", Password: "hashedpassword", Email: "admin@mail.com"}
		mockRepo.On("SaveUser", ctx, user).Return(nil)

		err := service.RegisterUser(ctx, user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: email is required", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		user := &domain.User{Username: "admin", Password: "hashedpassword", Email: ""}
		err := service.RegisterUser(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidEmail, err)
		mockRepo.AssertNotCalled(t, "SaveUser", mock.Anything, mock.Anything)
	})

	t.Run("Error: name too short", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		user := &domain.User{Username: "ad", Password: "hashedpassword", Email: "admin@mail.com"}
		err := service.RegisterUser(ctx, user)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidName, err)
		mockRepo.AssertNotCalled(t, "SaveUser", mock.Anything, mock.Anything)
	})

	t.Run("Error: Database error - db connection lost", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		user := &domain.User{Username: "admin", Password: "hashedpassword", Email: "admin@mail.com"}
		dbError := errors.New("db connection lost")

		mockRepo.On("SaveUser", ctx, user).Return(dbError)

		err := service.RegisterUser(ctx, user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db connection lost")
		mockRepo.AssertExpectations(t)
	})
}

func TestRegisterMember_EmptyPlan(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	member := &domain.Member{Name: "Alex", Email: "alex@mail.com", Plan: ""}
	err := service.RegisterMember(ctx, member)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidPlan, err)
	mockRepo.AssertNotCalled(t, "SaveMember", mock.Anything, mock.Anything)
}

func TestGetMemberStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: member found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		expected := &domain.Member{ID: "1", Name: "Alex", Email: "alex@mail.com", Plan: "Gold", Status: "active"}
		mockRepo.On("GetMemberByID", ctx, "1").Return(expected, nil)

		member, err := service.GetMemberStatus(ctx, "1")

		assert.NoError(t, err)
		assert.Equal(t, expected, member)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: member not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetMemberByID", ctx, "999").Return((*domain.Member)(nil), errors.New("member not found"))

		member, err := service.GetMemberStatus(ctx, "999")

		assert.Error(t, err)
		assert.Nil(t, member)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: valid credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		hashedPassword, _ := utils.HashPassword("secret123")
		user := domain.User{ID: "1", Username: "admin", Password: hashedPassword, Email: "encrypted@mail.com", Role: "ADMIN"}
		mockRepo.On("GetUserByUsername", ctx, "admin").Return(user, nil)

		resp, err := service.LoginUser(ctx, "admin", "secret123")

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "admin", resp.User.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		mockRepo.On("GetUserByUsername", ctx, "unknown").Return(domain.User{}, errors.New("no rows"))

		resp, err := service.LoginUser(ctx, "unknown", "secret123")

		assert.Error(t, err)
		assert.Empty(t, resp.Token)
		assert.Contains(t, err.Error(), "no user found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		hashedPassword, _ := utils.HashPassword("secret123")
		user := domain.User{ID: "1", Username: "admin", Password: hashedPassword, Email: "encrypted@mail.com", Role: "ADMIN"}
		mockRepo.On("GetUserByUsername", ctx, "admin").Return(user, nil)

		resp, err := service.LoginUser(ctx, "admin", "wrongpassword")

		assert.Error(t, err)
		assert.Empty(t, resp.Token)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: user found", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		expected := domain.User{ID: "1", Username: "admin", Role: "ADMIN"}
		mockRepo.On("GetUserByUsername", ctx, "admin").Return(expected, nil)

		user, err := service.GetUserByUsername(ctx, "admin")

		assert.NoError(t, err)
		assert.Equal(t, expected, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error: user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewIdentityService(mockRepo)

		mockRepo.On("GetUserByUsername", ctx, "unknown").Return(domain.User{}, errors.New("not found"))

		user, err := service.GetUserByUsername(ctx, "unknown")

		assert.Error(t, err)
		assert.Empty(t, user)
		mockRepo.AssertExpectations(t)
	})
}
