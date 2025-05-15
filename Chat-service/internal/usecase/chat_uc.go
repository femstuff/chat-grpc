package usecase

import (
	"context"
	"errors"
	"time"

	"chat-grpc/Chat-service/internal/broker"
	"chat-grpc/Chat-service/internal/repository"
	"chat-grpc/proto_gen"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatUseCaseInterface interface {
	Create(usernames []string) (int64, error)
	Delete(chatID int64) error
	SendMessage(chatID int64, from, text string, timestamp time.Time) error
	GetChatHistory(ctx context.Context, chatID int64) ([]*proto_gen.Message, error)
	Subscribe(subject string, handler func(*proto_gen.Message)) (*nats.Subscription, error)
}

type ChatUseCase struct {
	repo   repository.ChatRepo
	log    *zap.Logger
	broker broker.Broker
}

func NewChatUseCase(repo repository.ChatRepo, log *zap.Logger, broker broker.Broker) *ChatUseCase {
	return &ChatUseCase{repo: repo, log: log, broker: broker}
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

	timestamp = time.Now().Local()
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

	return uc.broker.Publish(msg)
}

func (uc *ChatUseCase) GetChatHistory(ctx context.Context, chatID int64) ([]*proto_gen.Message, error) {
	if chatID == 0 {
		return nil, errors.New("invalid chat ID")
	}

	return uc.repo.GetMessagesByChatID(ctx, chatID)
}

func (uc *ChatUseCase) Subscribe(subject string, handler func(*proto_gen.Message)) (*nats.Subscription, error) {
	return uc.broker.Subscribe(subject, handler)
}
