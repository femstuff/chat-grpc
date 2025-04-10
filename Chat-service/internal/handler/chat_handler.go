package handler

import (
	"context"
	"errors"
	"fmt"

	"chat-grpc/Chat-service/internal/usecase"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
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

	cs.log.Info("Message successfully sent", zap.Int64("chat_id", req.ChatId), zap.String("from", req.From))

	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) GetMessages(ctx context.Context, req *proto_gen.GetMessagesRequest) (*proto_gen.GetMessagesResponse, error) {
	messages, err := cs.useCase.GetChatHistory(ctx, req.ChatId)
	if err != nil {
		cs.log.Error("failed to get messages", zap.Error(err))
		return nil, errors.New("failed to get messages")
	}

	return &proto_gen.GetMessagesResponse{Messages: messages}, nil
}

func (cs *ChatService) Connect(req *proto_gen.ConnectRequest, stream proto_gen.ChatService_ConnectServer) error {
	ctx := stream.Context()
	chatID := req.ChatId

	messages, err := cs.useCase.GetChatHistory(ctx, chatID)
	if err != nil {
		return fmt.Errorf("error loading chat history: %w", err)
	}

	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			return fmt.Errorf("error sending history message: %w", err)
		}
	}

	subject := fmt.Sprintf("chat.%d", chatID)
	sub, err := cs.useCase.Subscribe(subject, func(msg *proto_gen.Message) {
		if err := stream.Send(msg); err != nil {
			cs.log.Warn("failed to send message stream", zap.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS: %w", err)
	}
	defer sub.Unsubscribe()

	<-ctx.Done()
	return nil
}
