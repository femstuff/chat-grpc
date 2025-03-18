package service

import (
	"chat-grpc/internal/entity"
	"chat-grpc/internal/usecase"
	"chat-grpc/proto"
	"errors"
)

type ChatService struct {
	useCase *usecase.ChatUseCase
}

func NewChatService(useCase *usecase.ChatUseCase) *ChatService {
	return &ChatService{useCase: useCase}
}

func (cs *ChatService) CreateChat(req *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	chatId, err := cs.useCase.Create(req.Name, req.Users, entity.TypeChat(req.Type))
	if err != nil {
		return nil, errors.New("fail with create chat")
	}
	return &proto.CreateChatResponse{ChatId: int64(chatId)}, nil
}

func (cs *ChatService) DeleteChat(req *proto.DeleteChatRequest) (*proto.DeleteChatResponse, error) {
	err := cs.useCase.Delete(req.ChatId)
	if err != nil {
		return nil, errors.New("fail with delete chat")
	}
	return &proto.DeleteChatResponse{}, nil
}

func (cs *ChatService) SendMessage(req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	err := cs.useCase.SendMessage(req.ChatId, req.Text, req.SenderId)
	if err != nil {
		return nil, errors.New("fail with send message")
	}
	return &proto.SendMessageResponse{}, nil
}

func (cs *ChatService) ConnectToChat(req *proto.ConnectChatRequest) (*proto.ConnectChatResponse, error) {
	err := cs.useCase.Connect(req.ChatId, req.UserId)
	if err != nil {
		return nil, errors.New("error with connect to chat")
	}
	return &proto.ConnectChatResponse{}, nil
}
