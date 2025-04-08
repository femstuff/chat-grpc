package usecase

import (
	"context"
	"errors"
	"sync"
	"time"

	"chat-grpc/Chat-service/internal/repository"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatUseCaseInterface interface {
	Create(usernames []string) (int64, error)
	Delete(chatID int64) error
	SendMessage(chatID int64, from, text string, timestamp time.Time) error
	GetChatHistory(ctx context.Context, chatID int64) ([]*proto_gen.Message, error)
	SubscribeToChat(chatID int64) <-chan *proto_gen.Message
	PublishMessage(msg *proto_gen.Message)
}

type ChatUseCase struct {
	repo      repository.ChatRepo
	log       *zap.Logger
	msgStream map[int64]chan *proto_gen.Message
	mu        sync.Mutex
}

func NewChatUseCase(repo repository.ChatRepo, log *zap.Logger) *ChatUseCase {
	return &ChatUseCase{repo: repo, log: log, msgStream: make(map[int64]chan *proto_gen.Message)}
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

	modText, err := uc.repo.SendMessage(chatID, from, text, timestamp)
	if err != nil {
		return err
	}

	msg := &proto_gen.Message{
		ChatId:    chatID,
		From:      from,
		Text:      modText,
		Timestamp: timestamppb.New(timestamp),
	}

	uc.PublishMessage(msg)

	return nil
}

func (uc *ChatUseCase) GetChatHistory(ctx context.Context, chatID int64) ([]*proto_gen.Message, error) {
	if chatID == 0 {
		return nil, errors.New("invalid chat ID")
	}

	return uc.repo.GetMessagesByChatID(ctx, chatID)
}

func (uc *ChatUseCase) SubscribeToChat(chatID int64) <-chan *proto_gen.Message {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, exists := uc.msgStream[chatID]; !exists {
		uc.msgStream[chatID] = make(chan *proto_gen.Message, 10)
	}

	return uc.msgStream[chatID]
}

func (uc *ChatUseCase) PublishMessage(msg *proto_gen.Message) {
	uc.mu.Lock()
	stream, exists := uc.msgStream[msg.ChatId]
	uc.mu.Unlock()

	if exists {
		select {
		case stream <- msg:
		default:
			uc.log.Warn("Канал подписки переполнен, сообщение потеряно", zap.Int64("chatID", msg.ChatId))
		}
	}
}
