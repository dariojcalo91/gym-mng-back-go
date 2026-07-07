package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/middleware"
	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/storage/postgres"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system environment")
	}

	if err := run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("could not connect to DB: %w", err)
	}
	defer pool.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO(repo layer): estos constructores todavía no existen — próxima sesión.
	memberRepo := postgres.NewMemberRepo(pool)
	checkInRepo := postgres.NewCheckInRepo(pool)
	memberService := sv.NewMemberService(memberRepo, checkInRepo)

	identityRepo := postgres.NewIdentityRepo(pool)
	identityService := sv.NewIdentityService(identityRepo)

	emailPool := sv.NewEmailWorkerPool(10)
	emailPool.Start(ctx, 3)

	expiringRepo := postgres.NewReminderRepo(pool)
	reminderService := sv.NewReminderService(expiringRepo, emailPool)
	go runDailyReminders(ctx, reminderService)

	jwtManager := utils.NewJWTManager()
	handler := &grpcHandler{members: memberService, identityService: identityService}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor,
			loggingInterceptor,
			middleware.AuthInterceptor(jwtManager),
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
		emailPool.Shutdown()
	}()

	slog.Info("gRPC server and worker pools running", "port", "50051")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("error serving: %w", err)
	}

	return nil
}
