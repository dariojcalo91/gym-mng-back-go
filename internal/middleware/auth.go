package middleware

import (
	"context"
	"strings"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey int

const authInfoKey ctxKey = iota

type AuthInfo struct {
	UserID string
	GymID  string
	Role   string
}

// publicMethods are the unique RPCs that do not require a token — before login
// or register, the client cannot have one yet.
var publicMethods = map[string]bool{
	"/gym.AuthService/Login":        true,
	"/gym.AuthService/Logout":       true,
	"/gym.UserService/RegisterUser": true,
}

func AuthInterceptor(jwtManager *utils.JWTManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}
		tokenString := strings.TrimPrefix(values[0], "Bearer ")

		claims, err := jwtManager.Validate(tokenString)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, authInfoKey, AuthInfo{
			UserID: claims.UserID,
			GymID:  claims.GymID,
			Role:   claims.Role,
		})
		return handler(ctx, req)
	}
}

// FromContext exposes the identity of the caller. Handlers use it instead of trusting any gym_id
// the client sends in the body.
func FromContext(ctx context.Context) (AuthInfo, bool) {
	info, ok := ctx.Value(authInfoKey).(AuthInfo)
	return info, ok
}
