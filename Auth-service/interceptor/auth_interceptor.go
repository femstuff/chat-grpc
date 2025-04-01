package interceptor

import (
	"context"
	"strings"

	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authClient *AuthClient
	log        *zap.Logger
}

func NewAuthInterceptor(authClient *AuthClient, log *zap.Logger) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient, log: log}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is missing")
		}

		tokenStr := authHeaders[0]
		a.log.Info("Received authorization token", zap.String("token", tokenStr))
		str := strings.Split(tokenStr, "Bearer ")
		a.log.Info("Received authorization token", zap.String("token", str[1]))

		err := a.authClient.CheckToken(ctx, &proto_gen.CheckTokenRequest{Token: str[1]})
		if err != nil {
			a.log.Info("Authorization failed:", zap.Error(err))
			return nil, status.Error(codes.PermissionDenied, "authorization failed")
		}

		return handler(ctx, req)
	}
}
