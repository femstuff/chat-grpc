package mocks

import (
	"log"

	"chat-grpc/internal/entity"
)

type MockChatUseCase struct {
	CreateFunc      func(name string, users []int64, chatType entity.TypeChat) (int, error)
	DeleteFunc      func(chatID int64) error
	SendMessageFunc func(chatID int64, contentMsg string, sender int64) error
	ConnectFunc     func(chatID, userID int64) error
}

func (m *MockChatUseCase) Create(name string, users []int64, chatType entity.TypeChat) (int, error) {
	log.Print("mock layer\n")
	return m.CreateFunc(name, users, chatType)
}

func (m *MockChatUseCase) Delete(chatID int64) error {
	log.Print("mock layer del\n")
	return m.DeleteFunc(chatID)
}

func (m *MockChatUseCase) SendMessage(chatID int64, contentMsg string, sender int64) error {
	return m.SendMessageFunc(chatID, contentMsg, sender)
}

func (m *MockChatUseCase) Connect(chatID, userID int64) error {
	return m.ConnectFunc(chatID, userID)
}
