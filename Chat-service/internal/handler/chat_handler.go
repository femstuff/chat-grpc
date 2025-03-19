package handler

import (
	"errors"

	"chat-grpc/Chat-service/internal/entity"
	"chat-grpc/Chat-service/internal/usecase"
	"chat-grpc/proto_gen"
)

type ChatService struct {
	useCase usecase.ChatUseCaseInterface
}

func NewChatService(useCase usecase.ChatUseCaseInterface) *ChatService {
	return &ChatService{useCase: useCase}
}

func (cs *ChatService) CreateChat(req *proto_gen.CreateChatRequest) (*proto_gen.CreateChatResponse, error) {
	chatType, err := entity.StringType(req.Type)
	if err != nil {
		return nil, errors.New("invalid chat type")
	}

	chatId, err := cs.useCase.Create(req.Name, req.Users, chatType)
	if err != nil {
		return nil, errors.New("fail with create chat")
	}
	return &proto_gen.CreateChatResponse{ChatId: int64(chatId)}, nil
}

func (cs *ChatService) DeleteChat(req *proto_gen.DeleteChatRequest) (*proto_gen.ChatEmpty, error) {
	err := cs.useCase.Delete(req.ChatId)
	if err != nil {
		return nil, errors.New("fail with delete chat")
	}
	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) SendMessage(req *proto_gen.SendMessageRequest) (*proto_gen.ChatEmpty, error) {
	err := cs.useCase.SendMessage(req.Sender, req.Text, req.Timestamp)
	if err != nil {
		return nil, errors.New("fail with send message")
	}

	return &proto_gen.ChatEmpty{}, nil
}

func (cs *ChatService) ConnectToChat(req *proto_gen.ConnectChatRequest) (*proto_gen.Message, error) {
	err := cs.useCase.Connect(req.ChatId, req.UserId)
	if err != nil {
		return nil, errors.New("error with connect to chat")
	}
	return &proto_gen.Message{}, nil
}
