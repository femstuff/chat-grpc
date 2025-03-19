package handler

import (
	"testing"

	"chat-grpc/internal/entity"
	"chat-grpc/internal/usecase/mocks"
	"chat-grpc/proto"

	"github.com/stretchr/testify/assert"
)

func TestCreateChat(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		CreateFunc: func(name string, users []int64, chatType entity.TypeChat) (int, error) {
			return 1, nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto.CreateChatRequest{
		Name:  "test",
		Users: []int64{1, 2},
		Type:  proto.ChatType(entity.PrivateChat),
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

	req := &proto.DeleteChatRequest{
		ChatId: 1,
	}
	resp, err := service.DeleteChat(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestSendMessage(t *testing.T) {
	mockUseCase := &mocks.MockChatUseCase{
		SendMessageFunc: func(chatID int64, contentMsg string, sender int64) error {
			return nil
		},
	}
	service := NewChatService(mockUseCase)

	req := &proto.SendMessageRequest{
		ChatId:   1,
		Text:     "test msg",
		SenderId: 1,
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

	req := &proto.ConnectChatRequest{
		ChatId: 1,
		UserId: 1,
	}
	ressp, err := service.ConnectToChat(req)

	assert.NoError(t, err)
	assert.NotNil(t, ressp)
}
