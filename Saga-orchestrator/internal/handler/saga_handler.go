package handler

import (
	"context"

	"chat-grpc/Saga-orchestrator/internal/usecase"
	pb "chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SagaHandler struct {
	pb.UnimplementedSagaServiceServer
	sagaService *usecase.SagaService
	log         *zap.Logger
}

func NewSagaHandler(service *usecase.SagaService, log *zap.Logger) *SagaHandler {
	return &SagaHandler{sagaService: service, log: log}
}

func (h *SagaHandler) StartSaga(ctx context.Context, req *pb.StartSagaRequest) (*emptypb.Empty, error) {
	msg := &usecase.Message{
		ID:     req.GetMessageId(),
		Text:   req.GetText(),
		Emails: req.GetEmails(),
		ChatID: req.GetChatId(),
	}

	userID, err := h.sagaService.GetUserIdFromEmail(ctx, msg.Emails[0])

	err = h.sagaService.SendMessageWithNotification(ctx, msg, req.GetChatId(), userID)
	if err != nil {
		h.log.Error("Saga failed", zap.Error(err))
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
