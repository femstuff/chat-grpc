package usecase

import (
	"errors"
	"time"

	"chat-grpc/Chat-service/internal/entity"
	"chat-grpc/Chat-service/internal/repository"
)

type ChatUseCaseInterface interface {
	Create(name string, users []int64, chatType entity.TypeChat) (int, error)
	Delete(chatID int64) error
	SendMessage(sender, text string, timestamp string) error
	Connect(chatID, userID int64) error
}

type ChatUseCase struct {
	repo repository.ChatRepo
}

func NewChatUseCase(repo repository.ChatRepo) *ChatUseCase {
	return &ChatUseCase{repo: repo}
}

func (uc *ChatUseCase) Create(name string, users []int64, chatType entity.TypeChat) (int, error) {
	chat := &entity.Chat{
		Name:      name,
		Users:     users,
		Type:      chatType,
		CreatedAt: time.Now().UTC(),
	}

	return uc.repo.CreateChat(chat)
}

func (uc *ChatUseCase) Delete(chatID int64) error {
	return uc.repo.DeleteChat(chatID)
}

func (uc *ChatUseCase) SendMessage(sender, text, timestamp string) error {
	if sender == "" || text == "" || timestamp == "" {
		return errors.New("invalid msg params")
	}
	
	return nil
}

func (uc *ChatUseCase) Connect(chatID, userID int64) error {
	if chatID == 0 || userID == 0 {
		return errors.New("invalid connect params")
	}

	return nil
}
