package usecase

import (
	"context"

	"chat-grpc/Saga-orchestrator/internal/client"
	"chat-grpc/Saga-orchestrator/internal/repository"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
)

type Message struct {
	ID     int64
	ChatID int64
	Text   string
	Emails []string
}

type SagaService struct {
	repo     *repository.SagaRepo
	notifier client.NotificationClient
	log      *zap.Logger
}

func NewSagaService(repo *repository.SagaRepo, notifier client.NotificationClient, log *zap.Logger) *SagaService {
	return &SagaService{
		repo:     repo,
		notifier: notifier,
		log:      log,
	}
}

func (s *SagaService) SendMessageWithNotification(ctx context.Context, msg *Message, chatID, userID int64) error {
	s.log.Info("Saving message", zap.Int64("message_id", msg.ID))

	err := s.repo.SaveMessage(msg.ID, chatID, msg.Text, userID)
	if err != nil {
		s.log.Error("Failed to save message", zap.Error(err))
		return err
	}

	s.log.Info("Sending notifications", zap.Int("email_count", len(msg.Emails)))

	for _, email := range msg.Emails {
		emailReq := &proto_gen.SendEmailRequest{
			To:      email,
			Subject: "",
			Body:    msg.Text,
		}
		err := s.notifier.SendEmail(ctx, emailReq)
		if err != nil {
			s.log.Error("Notification failed", zap.String("email", email), zap.Error(err))
			rollbackErr := s.repo.DeleteMessage(msg.ID)
			if rollbackErr != nil {
				s.log.Error("Failed to rollback message", zap.Error(rollbackErr))
			}
			return err
		}
	}

	s.log.Info("Saga completed successfully", zap.Int64("message_id", msg.ID))
	return nil
}

func (s *SagaService) GetUserIdFromEmail(ctx context.Context, email string) (int64, error) {
	userID, err := s.repo.GetUserIdFromEmail(email)
	if err != nil {
		s.log.Error("Failed to get user_id", zap.String("email", email), zap.Error(err))
		return 0, err
	}
	return userID, nil
}
