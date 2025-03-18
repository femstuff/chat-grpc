package service

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
