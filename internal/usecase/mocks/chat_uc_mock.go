package mocks

import (
	"chat-grpc/internal/entity"
)

type MockChatUseCase struct {
	CreateFunc      func(name string, users []int64, chatType entity.TypeChat) (int, error)
	DeleteFunc      func(chatID int64) error
	SendMessageFunc func(sender, text, timestamp string) error
	ConnectFunc     func(chatID, userID int64) error
}

func (m *MockChatUseCase) Create(name string, users []int64, chatType entity.TypeChat) (int, error) {
	return m.CreateFunc(name, users, chatType)
}

func (m *MockChatUseCase) Delete(chatID int64) error {
	return m.DeleteFunc(chatID)
}

func (m *MockChatUseCase) SendMessage(sender, text string, timestamp string) error {
	return m.SendMessageFunc(sender, text, timestamp)
}

func (m *MockChatUseCase) Connect(chatID, userID int64) error {
	return m.ConnectFunc(chatID, userID)
}
