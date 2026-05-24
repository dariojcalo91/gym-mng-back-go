package main

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"

	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var errNotFound = errors.New("not found")

const bufSize = 1024 * 1024

type e2eClients struct {
	gym    pb.GymServiceClient
	user   pb.UserServiceClient
	auth   pb.AuthServiceClient
	mock   *sv.MockRepository
	userMC *sv.MockUserRepo
}

func setupE2EServer(t *testing.T) (*e2eClients, func()) {
	if os.Getenv("SKIP_E2E") == "true" {
		t.Skip("Skipping e2e test: SKIP_E2E=true")
	}

	os.Setenv("AES_MASTER_KEY", "abcdefghijklmnopqrstuvwxyz123456")
	os.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")

	mockRepo := new(sv.MockRepository)
	service := sv.NewService(mockRepo)

	mockUserRepo := new(sv.MockUserRepo)
	identityService := sv.NewIdentityService(mockUserRepo)

	handler := &grpcHandler{
		service:         service,
		identityService: identityService,
	}

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterGymServiceServer(s, handler)
	pb.RegisterUserServiceServer(s, handler)
	pb.RegisterAuthServiceServer(s, handler)

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	clients := &e2eClients{
		gym:    pb.NewGymServiceClient(conn),
		user:   pb.NewUserServiceClient(conn),
		auth:   pb.NewAuthServiceClient(conn),
		mock:   mockRepo,
		userMC: mockUserRepo,
	}

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return clients, cleanup
}

func TestE2E_CreateMember_Success(t *testing.T) {
	clients, cleanup := setupE2EServer(t)
	defer cleanup()

	clients.mock.On("SaveMember", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	resp, err := clients.gym.CreateMember(ctx, &pb.MemberRequest{
		Name:  "Alex",
		Email: "alex@test.com",
		Plan:  "Gold",
	})

	require.NoError(t, err)
	require.Equal(t, "Success", resp.Status)
	clients.mock.AssertExpectations(t)
}

func TestE2E_CreateMember_ValidationError(t *testing.T) {
	clients, cleanup := setupE2EServer(t)
	defer cleanup()

	ctx := context.Background()
	_, err := clients.gym.CreateMember(ctx, &pb.MemberRequest{
		Name:  "Al",
		Email: "alex@test.com",
		Plan:  "Gold",
	})

	require.Error(t, err)
	clients.mock.AssertNotCalled(t, "SaveMember", mock.Anything, mock.Anything)
}

func TestE2E_RegisterUser_Success(t *testing.T) {
	clients, cleanup := setupE2EServer(t)
	defer cleanup()

	clients.userMC.On("SaveUser", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	resp, err := clients.user.RegisterUser(ctx, &pb.RegisterUserRequest{
		Username: "admin",
		Password: "secret123",
		Email:    "admin@test.com",
		Role:     "ADMIN",
	})

	require.NoError(t, err)
	require.Equal(t, "User registered successfully", resp.Message)
	clients.userMC.AssertExpectations(t)
}

func TestE2E_GetMemberStatus_NotFound(t *testing.T) {
	clients, cleanup := setupE2EServer(t)
	defer cleanup()

	clients.mock.On("GetMemberByID", mock.Anything, "bad-id").Return(nil, errNotFound)

	ctx := context.Background()
	_, err := clients.gym.GetMemberStatus(ctx, &pb.IdRequest{
		Id: "bad-id",
	})

	require.Error(t, err)
	clients.mock.AssertExpectations(t)
}

func TestE2E_Logout_Success(t *testing.T) {
	clients, cleanup := setupE2EServer(t)
	defer cleanup()

	ctx := context.Background()
	resp, err := clients.auth.Logout(ctx, &pb.LogoutRequest{Token: "some-token"})

	require.NoError(t, err)
	require.Equal(t, "Logged out successfully", resp.Message)
}
