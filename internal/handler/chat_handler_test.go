package handler

import (
	"testing"

	"chat-grpc/internal/entity"
	"chat-grpc/internal/usecase/mocks"
	"chat-grpc/proto_gen"

	"github.com/stretchr/testify/assert"
)

func TestCreateChat(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		CreateFunc: func(name string, users []int64, chatType entity.TypeChat) (int, error) {
			return 1, nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto_gen.CreateChatRequest{
		Name:  "test",
		Users: []int64{1, 2},
		Type:  "public",
	}
	resp, err := service.CreateChat(req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.ChatId)
}

func TestDeleteChat(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		DeleteFunc: func(chatID int64) error {
			return nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto_gen.DeleteChatRequest{
		ChatId: 1,
	}
	resp, err := service.DeleteChat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestSendMessage(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		SendMessageFunc: func(sender, text, timestamps string) error {
			return nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto_gen.SendMessageRequest{
		Sender:    "user",
		Text:      "test msg",
		Timestamp: "2025-01-01T16:22:16Z",
	}
	resp, err := service.SendMessage(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestConnectChat(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		ConnectFunc: func(chatID, userID int64) error {
			return nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto_gen.ConnectChatRequest{
		ChatId: 1,
		UserId: 1,
	}
	resp, err := service.ConnectToChat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
