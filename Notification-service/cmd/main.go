package main

import (
	"chat-grpc/Auth-service/interceptor"
	"chat-grpc/Notification-service/internal/broker"
	"chat-grpc/Notification-service/internal/handler"
	"chat-grpc/Notification-service/internal/usecase"
	"chat-grpc/pkg/config"
	"chat-grpc/pkg/logger"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.LoadConfig()
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to init logger", zap.Error(err))
	}
	defer log.Sync()

	nc, _ := nats.Connect(cfg.NatsUrl)
	grpcConn, err := grpc.NewClient(cfg.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	defer grpcConn.Close()

	emailSender := handler.NewEmailSender(cfg.SmtpUser, cfg.SmtpPass, cfg.SmtpHost, cfg.SmtpPort, log)

	authClient := interceptor.NewAuthClient(grpcConn, log)
	notifier := usecase.NewNotifier(authClient, emailSender, log)

	consumer := broker.NewNatsConsumer(nc, "chat.*", log, notifier.Notify)
	if err := consumer.Start(); err != nil {
		log.Fatal("Failed to start consumer", zap.Error(err))
	}

	select {}
}
