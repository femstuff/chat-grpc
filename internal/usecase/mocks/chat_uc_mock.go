package mocks

import "chat-grpc/internal/entity"

type mockChatUC struct {
	createChatMock  func(name string, users []int64, chatType entity.TypeChat) (int, error)
	deleteChatMock  func(chatID int64) error
	sendMessageMock func(chatID int64, content string, sender int64) error
	connectChatMock func(chatID, userID int64) error
}

func (m *mockChatUC) Create(name string, users []int64, chatType entity.TypeChat) (int, error) {
	return m.createChatMock(name, users, chatType)
}

func (m *mockChatUC) Delete(chatID int64) error {
	return m.deleteChatMock(chatID)
}

func (m *mockChatUC) SendMessage(chatID int64, content string, sender int64) error {
	return m.sendMessageMock(chatID, content, sender)
}

func (m *mockChatUC) Connect(chatID, userID int64) error {
	return m.connectChatMock(chatID, userID)
}
