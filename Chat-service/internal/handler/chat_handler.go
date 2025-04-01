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
	err := cs.useCase.SendMessage(req.ChatId, req.From, req.Text, req.Timestamp.AsTime())
	if err != nil {
		cs.log.Error("failed to send message", zap.Error(err))
		return nil, errors.New("failed to send message")
	}

	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) Connect(req *proto_gen.ConnectRequest, stream proto_gen.ChatService_ConnectServer) error {
	if req.ChatId == 0 {
		cs.log.Error("invalid chat ID", zap.Int64("chat_id", req.ChatId))
		return errors.New("invalid chat ID")
	}

	cs.log.Info("Connecting to chat", zap.Int64("chat_id", req.ChatId))

	for {
		messages, err := cs.useCase.GetMessages(req.ChatId)
		if err != nil {
			cs.log.Error("failed to get messages", zap.Error(err))
			return err
		}

		for _, msg := range messages {
			err := stream.Send(&proto_gen.Message{
				From:      string(msg.Sender),
				Text:      msg.Content,
				Timestamp: timestamppb.New(msg.CreatedAt),
			})
			if err != nil {
				cs.log.Error("failed to stream message", zap.Error(err))
				return err
			}
		}

		time.Sleep(2 * time.Second)
	}
}
