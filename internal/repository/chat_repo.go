package repository

import "chat-grpc/internal/entity"

type ChatRepo interface {
	CreateChat(chat *entity.Chat) (int, error)
	DeleteChat(id int64) error
}
