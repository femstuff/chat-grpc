package usecase

import (
	"chat-grpc/internal/entity"
	"chat-grpc/internal/repository"
	"time"
)

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

func (uc *ChatUseCase) SendMessage(chatID int64, contentMsg string, sender int64) error {
	msg := &entity.Message{
		ChatID:    chatID,
		Sender:    sender,
		Content:   contentMsg,
		CreatedAt: time.Now().UTC(),
	}

	return uc.repo.SendMessage(msg)
}

func (uc *ChatUseCase) Connect(chatID, userID int64) error {
	return uc.repo.ConnectChat(chatID, userID)
}
