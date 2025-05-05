package usecase

import (
	"context"

	"chat-grpc/Auth-service/interceptor"
	"chat-grpc/Notification-service/internal/handler"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
)

type AuthClientInterface interface {
	GetChatUsers(ctx context.Context, chatID int64) ([]int64, error)
	GetChatUsersEmails(ctx context.Context, chatID int64) ([]string, error)
}

type Notifier struct {
	authClient  interceptor.AuthClientInterface
	emailSender *handler.EmailSender
	log         *zap.Logger
}

func NewNotifier(authClient interceptor.AuthClientInterface, emailSender *handler.EmailSender, log *zap.Logger) *Notifier {
	return &Notifier{
		authClient:  authClient,
		emailSender: emailSender,
		log:         log,
	}
}

func (n *Notifier) Notify(ctx context.Context, msg *proto_gen.Message) {
	n.log.Info("Processing message", zap.String("message", msg.Text))

	emails, err := n.authClient.GetChatUsersEmails(ctx, msg.ChatId)
	if err != nil {
		n.log.Error("Failed to get emails from auth service", zap.Error(err))
		return
	}

	n.log.Info("Sending notifications", zap.Int("emailCount", len(emails)))
	for _, email := range emails {
		err := n.emailSender.Send(email, "New message in chat", msg.Text)
		if err != nil {
			n.log.Error("Failed to send email", zap.Error(err))
			return
		} else {
			n.log.Info("Email sent", zap.String("email", email))
		}
	}
}
