package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/storage/postgres"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto" // Import generated code from the .proto file
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedGymServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedAuthServiceServer
	service         *sv.GymService
	identityService *sv.IdentityService
}

// Implement methods we defined in .proto file (handlres)
func (h *grpcHandler) CreateMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	// 1. convert from gRPC to Domain
	m := &domain.Member{
		Name:  req.Name,
		Email: req.Email,
		Plan:  req.Plan,
	}

	// 2. Call the business logic
	err := h.service.RegisterMember(ctx, m)
	if err != nil {
		return nil, err
	}

	return &pb.MemberResponse{
		Status:  "Success",
		Message: "Member registered successfully",
	}, nil
}

func (h *grpcHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	// 1. convert from gRPC to Domain
	u := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

	// 2. Call the business logic
	err := h.identityService.RegisterUser(ctx, u)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserId:  u.ID,
		Message: "User registered successfully",
	}, nil
}

func (h *grpcHandler) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	// 1. Call the business logic
	user, err := h.identityService.GetUserByUsername(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:       user.ID,
		Username: user.Username,
		// Email:    user.Email,
		Role: user.Role,
	}, nil
}

func (h *grpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := h.identityService.LoginUser(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	userResponse := &pb.User{
		Id:       user.User.ID,
		Username: user.User.Username,
		Role:     user.User.Role,
	}

	return &pb.LoginResponse{
		Token: user.Token,
		User:  userResponse,
	}, nil
}

func main() {
	// Load environment variables from .env file if exists
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 1. DB connection
	dbURL := os.Getenv("DATABASE_URL")
	// for testing purposes, we could use the URL from container
	if dbURL == "" {
		dbURL = "postgres://gopher:gympassword@localhost:5432/gym_management?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer pool.Close() // Ensure the connection is closed when main exits (reserver word for defer)

	// 2. Initialize layers (Hexagonal Architecture)
	repo := postgres.NewStorage(pool) // repo: output adapter
	service := sv.NewService(repo)    // service: core + Worker Pool of emails (3 workers automatically started in NewService)

	identityRepo := postgres.NewUserRepo(pool)             // repo: output adapter for identity
	identityService := sv.NewIdentityService(identityRepo) // service: core for identity management

	handler := &grpcHandler{service: service, identityService: identityService} // handler: input adapter (gRPC) that translates gRPC calls to service methods

	// 3. Start Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGymServiceServer(s, handler)
	pb.RegisterUserServiceServer(s, handler)
	pb.RegisterAuthServiceServer(s, handler)

	log.Println("🚀 gRPC server and worker pools running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
