package interceptor

import (
	"context"

	proto "chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthClientInterface interface {
	GetChatUsersEmails(ctx context.Context, chatID int64) ([]string, error)
	GetChatUsers(ctx context.Context, chatID int64) ([]int64, error)
}

type AuthClient struct {
	client proto.AuthServiceClient
	log    *zap.Logger
}

func NewAuthClient(conn *grpc.ClientConn, log *zap.Logger) *AuthClient {
	return &AuthClient{
		client: proto.NewAuthServiceClient(conn),
		log:    log,
	}
}

func (a *AuthClient) CheckToken(ctx context.Context, req *proto.CheckTokenRequest) error {
	_, err := a.client.CheckToken(ctx, req)
	return err
}

func (a *AuthClient) GetChatUsersEmails(ctx context.Context, chatID int64) ([]string, error) {
	req := &proto.GetChatUsersEmailsRequest{ChatId: chatID}
	res, err := a.client.GetChatUsersEmails(ctx, req)
	if err != nil {
		a.log.Error("failed to get emails from auth service", zap.Error(err))
		return nil, err
	}
	return res.Emails, nil
}

func (a *AuthClient) GetChatUsers(ctx context.Context, chatID int64) ([]int64, error) {
	req := &proto.GetChatUsersRequest{ChatId: chatID}
	res, err := a.client.GetChatUsers(ctx, req)
	if err != nil {
		a.log.Error("failed to get chat users from auth service", zap.Error(err))
		return nil, err
	}
	return res.UserIds, nil
}
