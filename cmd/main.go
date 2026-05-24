package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/storage/postgres"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type grpcHandler struct {
	pb.UnimplementedGymServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedAuthServiceServer
	service         *sv.GymService
	identityService *sv.IdentityService
}

func (h *grpcHandler) CreateMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	m := &domain.Member{
		Name:  req.Name,
		Email: req.Email,
		Plan:  req.Plan,
	}

	err := h.service.RegisterMember(ctx, m)
	if err != nil {
		return nil, err
	}

	return &pb.MemberResponse{
		Status:  "Success",
		Message: "Member registered successfully",
	}, nil
}

func (h *grpcHandler) GetMemberStatus(ctx context.Context, req *pb.IdRequest) (*pb.MemberResponse, error) {
	member, err := h.service.GetMemberStatus(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.MemberResponse{
		Id:     member.ID,
		Name:   member.Name,
		Status: member.Status,
	}, nil
}

func (h *grpcHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	u := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

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
	user, err := h.identityService.GetUserByUsername(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:       user.ID,
		Username: user.Username,
		Role:     user.Role,
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

func (h *grpcHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{
		Message: "Logged out successfully",
	}, nil
}

func loggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	slog.Info("gRPC request", "method", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		slog.Error("gRPC request failed", "method", info.FullMethod, "error", err)
	}
	return resp, err
}

func recoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("gRPC handler panic", "method", info.FullMethod, "panic", r)
		}
	}()
	return handler(ctx, req)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn(".env file not found, using system environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL environment variable is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		slog.Error("could not connect to DB", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewStorage(pool)
	service := sv.NewService(repo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service.StartWorkers(ctx)

	identityRepo := postgres.NewUserRepo(pool)
	identityService := sv.NewIdentityService(identityRepo)

	handler := &grpcHandler{service: service, identityService: identityService}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("error listening", "error", err)
		os.Exit(1)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			loggingInterceptor,
		),
	)

	pb.RegisterGymServiceServer(s, handler)
	pb.RegisterUserServiceServer(s, handler)
	pb.RegisterAuthServiceServer(s, handler)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-stop
		slog.Info("received signal, shutting down", "signal", sig)
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		cancel()
		s.GracefulStop()
		service.Shutdown()
	}()

	slog.Info("gRPC server and worker pools running", "port", "50051")
	if err := s.Serve(lis); err != nil {
		slog.Error("error serving", "error", err)
		os.Exit(1)
	}
}
