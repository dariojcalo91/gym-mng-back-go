package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/gym"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/storage/postgres"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto" // Import generated code from the .proto file
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedGymServiceServer
	service *gym.GymService // Logic core
}

// Implement methods we defined in .proto file
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

func main() {
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
	repo := postgres.NewStorage(pool)         // repo: output adapter
	service := gym.NewService(repo)           // service: core + Worker Pool of emails (3 workers automatically started in NewService)
	handler := &grpcHandler{service: service} // handler: input adapter (gRPC) that translates gRPC calls to service methods

	// 3. Start Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGymServiceServer(s, handler)

	log.Println("🚀 gRPC server and worker pools running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
