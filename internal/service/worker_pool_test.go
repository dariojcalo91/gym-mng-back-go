package service

import (
	"context"
	"testing"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkerPool_ProcessesEmails(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service.StartWorkers(ctx)

	mockRepo.On("SaveMember", ctx, mock.Anything).Return(nil)

	member := &domain.Member{Name: "Test", Email: "email-worker@test.com", Plan: "Gold"}
	err := service.RegisterMember(ctx, member)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	member2 := &domain.Member{Name: "Test", Email: "email-worker2@test.com", Plan: "Gold"}
	err = service.RegisterMember(ctx, member2)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	cancel()
	service.Shutdown()

	mockRepo.AssertExpectations(t)
}

func TestWorkerPool_ShutdownWaitsForWorkers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())

	service.StartWorkers(ctx)

	mockRepo.On("SaveMember", ctx, mock.Anything).Return(nil)

	member := &domain.Member{Name: "Test", Email: "shutdown-test@test.com", Plan: "Gold"}
	err := service.RegisterMember(ctx, member)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	cancel()

	done := make(chan struct{})
	go func() {
		service.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown timed out - workers did not complete")
	}

	mockRepo.AssertExpectations(t)
}

func TestWorkerPool_ContextCancelStopsWorkers(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())

	service.emailChannel = make(chan string, 100)

	service.StartWorkers(ctx)

	mockRepo.On("SaveMember", ctx, mock.Anything).Return(nil)

	for range 10 {
		member := &domain.Member{Name: "Test", Email: "cancel-test@test.com", Plan: "Gold"}
		_ = service.RegisterMember(ctx, member)
	}

	cancel()

	done := make(chan struct{})
	go func() {
		service.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Shutdown timed out after context cancel")
	}

	mockRepo.AssertExpectations(t)
}

func TestWorkerPool_NonBlockingSend(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	service.emailChannel = make(chan string, 0)

	ctx := context.Background()
	mockRepo.On("SaveMember", ctx, mock.Anything).Return(nil)

	member := &domain.Member{Name: "Test", Email: "nonblock@test.com", Plan: "Gold"}
	err := service.RegisterMember(ctx, member)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
