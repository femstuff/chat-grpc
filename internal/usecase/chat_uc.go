package usecase

import (
	"errors"
	"log"
	"time"

	"chat-grpc/internal/entity"
	"chat-grpc/internal/repository"
)

type ChatUseCaseInterface interface {
	Create(name string, users []int64, chatType entity.TypeChat) (int, error)
	Delete(chatID int64) error
	SendMessage(chatID int64, contentMsg string, sender int64) error
	Connect(chatID, userID int64) error
}

type ChatUseCase struct {
	repo repository.ChatRepo
}

func NewChatUseCase(repo repository.ChatRepo) *ChatUseCase {
	return &ChatUseCase{repo: repo}
}

func (uc *ChatUseCase) Create(name string, users []int64, chatType entity.TypeChat) (int, error) {
	log.Printf("uc layer\n")
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

func (uc *ChatUseCase) SendMessage(chatID int64, contentMsg string, sender int64) error {
	if chatID == 0 || len(contentMsg) == 0 || sender == 0 {
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
