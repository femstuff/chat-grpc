package handler

import (
	"context"
	"errors"
	"time"

	"chat-grpc/Chat-service/internal/usecase"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService struct {
	useCase usecase.ChatUseCaseInterface
	proto_gen.UnimplementedChatServiceServer
	log *zap.Logger
}

func NewChatService(useCase usecase.ChatUseCaseInterface, log *zap.Logger) *ChatService {
	return &ChatService{useCase: useCase, log: log}
}

func (cs *ChatService) Create(ctx context.Context, req *proto_gen.CreateRequest) (*proto_gen.CreateResponse, error) {
	chatId, err := cs.useCase.Create(req.Usernames)
	if err != nil {
		cs.log.Error("failed to create chat", zap.Error(err))
		return nil, errors.New("failed to create chat")
	}

	return &proto_gen.CreateResponse{Id: chatId}, nil
}

func (cs *ChatService) Delete(ctx context.Context, req *proto_gen.DeleteRequest) (*proto_gen.ChatEmpty, error) {
	err := cs.useCase.Delete(req.Id)
	if err != nil {
		cs.log.Error("failed to delete chat", zap.Error(err))
		return nil, errors.New("failed to delete chat")
	}
	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) SendMessage(ctx context.Context, req *proto_gen.SendMessageRequest) (*proto_gen.ChatEmpty, error) {
	err := cs.useCase.SendMessage(req.From, req.Text, req.Timestamp.AsTime())
	if err != nil {
		cs.log.Error("failed to send message", zap.Error(err))
		return nil, errors.New("failed to send message")
	}

	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) Connect(req *proto_gen.ConnectRequest, stream proto_gen.ChatService_ConnectServer) error {
	cs.log.Info("Connecting to chat")

	for {
		msg := &proto_gen.Message{
			From:      "test",
			Text:      "test message",
			Timestamp: timestamppb.Now(),
		}
		if err := stream.Send(msg); err != nil {
			cs.log.Error("failed to stream send msg", zap.Error(err))
			return err
		}
		time.Sleep(2 * time.Second)
	}
}
