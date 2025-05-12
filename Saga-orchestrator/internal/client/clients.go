package client

import (
	"context"

	proto "chat-grpc/proto_gen"
	"google.golang.org/grpc"
)

type NotificationClient interface {
	SendEmail(ctx context.Context, req *proto.SendEmailRequest) error
}

type notificationClient struct {
	client proto.NotificationServiceClient
}

func NewNotificationClient(conn *grpc.ClientConn) NotificationClient {
	return &notificationClient{client: proto.NewNotificationServiceClient(conn)}
}

func (n *notificationClient) SendEmail(ctx context.Context, req *proto.SendEmailRequest) error {
	_, err := n.client.SendEmail(ctx, req)
	return err
}
