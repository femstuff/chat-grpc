package repository

import (
	"errors"
	"time"

	"chat-grpc/Chat-service/internal/entity"
	"go.uber.org/zap"
)

type ChatRepo interface {
	CreateChat(usernames []string) (int64, error)
	DeleteChat(id int64) error
	SendMessage(from, text string, timestamp time.Time) error
}

type chatRepository struct {
	messages []entity.Message
	log      *zap.Logger
}

func NewChatRepository(log *zap.Logger) ChatRepo {
	return &chatRepository{log: log}
}

func (r *chatRepository) CreateChat(usernames []string) (int64, error) {
	r.log.Info("creating chat")

	chatID := int64(len(r.messages) + 1)

	r.log.Info("success create chat with ", zap.Int64("chat_id", chatID))
	return chatID, nil
}

func (r *chatRepository) DeleteChat(id int64) error {
	r.log.Info("delete chat with ", zap.Int64("chat_id", id))

	r.log.Info("success delete chat with ", zap.Int64("chat_id", id))
	return nil
}

func (r *chatRepository) SendMessage(from, text string, timestamp time.Time) error {
	if from == "" || text == "" {
		r.log.Error("invalid message param")
		return errors.New("invalid message parameters")
	}

	r.log.Info("new message ", zap.String("from", from))
	return nil
}
