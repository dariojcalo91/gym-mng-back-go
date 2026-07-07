package main

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

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
