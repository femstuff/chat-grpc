package interceptor

import (
	"context"
	"errors"
	"strings"

	"chat-grpc/Auth-service/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	jwtService *jwt.JWTService
	log        *zap.Logger
}

func NewAuthInterceptor(jwtService *jwt.JWTService, log *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtService: jwtService,
		log:        log,
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		a.log.Info("Request", zap.String("req method", info.FullMethod))

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			a.log.Error("Failed with metadata in ctx")
			return nil, errors.New("Failed with metadata in ctx")
		}

		authHead := md.Get("authorization")
		if len(authHead) == 0 {
			a.log.Error("Authorization token is required")
			return nil, errors.New("Authorization token is required")
		}

		token := strings.TrimPrefix(authHead[0], "Bearer ")
		claims, err := a.jwtService.VerifyToken(token)
		if err != nil {
			a.log.Error("Invalid token")
			return nil, errors.New("Invalid token")
		}

		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		a.log.Info("User auth", zap.Int64("id", claims.UserID), zap.String("role", claims.Role))

		return handler(ctx, req)
	}
}
