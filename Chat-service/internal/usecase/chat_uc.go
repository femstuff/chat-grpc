package usecase

import (
	"errors"
	"time"

	"chat-grpc/Chat-service/internal/entity"
	"chat-grpc/Chat-service/internal/repository"
	"go.uber.org/zap"
)

type ChatUseCaseInterface interface {
	Create(usernames []string) (int64, error)
	Delete(chatID int64) error
	SendMessage(chatID int64, from, text string, timestamp time.Time) error
	GetMessages(chatID int64) ([]entity.Message, error)
}

type ChatUseCase struct {
	repo repository.ChatRepo
	log  *zap.Logger
}

func NewChatUseCase(repo repository.ChatRepo, log *zap.Logger) *ChatUseCase {
	return &ChatUseCase{repo: repo, log: log}
}

func (uc *ChatUseCase) Create(usernames []string) (int64, error) {
	if len(usernames) == 0 {
		return 0, errors.New("usernames list is empty")
	}

	return uc.repo.CreateChat(usernames)
}

func (uc *ChatUseCase) Delete(chatID int64) error {
	if chatID == 0 {
		return errors.New("invalid chat ID")
	}

	return uc.repo.DeleteChat(chatID)
}

func (uc *ChatUseCase) SendMessage(chatID int64, from, text string, timestamp time.Time) error {
	if chatID == 0 || from == "" || text == "" {
		return errors.New("invalid message parameters")
	}

	return uc.repo.SendMessage(chatID, from, text, timestamp)
}

func (uc *ChatUseCase) GetMessages(chatID int64) ([]entity.Message, error) {
	if chatID == 0 {
		return nil, errors.New("invalid chat ID")
	}

	return uc.repo.GetMessages(chatID)
}
