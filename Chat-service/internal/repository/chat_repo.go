package repository

import (
	"chat-grpc/Chat-service/internal/entity"
)

type ChatRepo interface {
	CreateChat(chat *entity.Chat) (int, error)
	DeleteChat(id int64) error
}
