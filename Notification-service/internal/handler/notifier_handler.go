package handler

import (
	"context"

	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationHandler struct {
	proto_gen.UnimplementedNotificationServiceServer
	emailSender EmailSender
	log         *zap.Logger
}

func NewNotificationHandler(emailSender EmailSender, log *zap.Logger) *NotificationHandler {
	return &NotificationHandler{emailSender: emailSender, log: log}
}

func (h *NotificationHandler) SendEmail(ctx context.Context, req *proto_gen.SendEmailRequest) (*emptypb.Empty, error) {
	err := h.emailSender.Send(req.To, req.Subject, req.Body)
	if err != nil {
		h.log.Error("Failed to send email", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to send email: %v", err)
	}
	return &emptypb.Empty{}, nil
}
