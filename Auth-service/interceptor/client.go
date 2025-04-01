package interceptor

import (
	"context"

	proto "chat-grpc/proto_gen"
	"google.golang.org/grpc"
)

type AuthClient struct {
	client proto.AuthServiceClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		client: proto.NewAuthServiceClient(conn),
	}
}

func (a *AuthClient) CheckToken(ctx context.Context, req *proto.CheckTokenRequest) error {
	_, err := a.client.CheckToken(ctx, req)
	return err
}
