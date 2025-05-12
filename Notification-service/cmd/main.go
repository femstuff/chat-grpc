package main

import (
	"net"

	"chat-grpc/Auth-service/interceptor"
	"chat-grpc/Notification-service/internal/broker"
	"chat-grpc/Notification-service/internal/handler"
	"chat-grpc/Notification-service/internal/usecase"
	"chat-grpc/pkg/config"
	"chat-grpc/pkg/logger"
	"chat-grpc/proto_gen"

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

	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		log.Fatal("Failed to connect to NATS", zap.Error(err))
	}

	authConn, err := grpc.NewClient(cfg.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	defer authConn.Close()

	authClient := interceptor.NewAuthClient(authConn, log)

	emailSender := handler.NewStubEmailSender(log)
	notifier := usecase.NewNotifier(authClient, emailSender, log)

	consumer := broker.NewNatsConsumer(nc, "chat.*", log, notifier.Notify)
	if err := consumer.Start(); err != nil {
		log.Fatal("Failed to start NATS consumer", zap.Error(err))
	}

	lis, err := net.Listen("tcp", ":"+cfg.NotificationPort)
	if err != nil {
		log.Fatal("Failed to start listener", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	notifHandler := handler.NewNotificationHandler(emailSender, log)
	proto_gen.RegisterNotificationServiceServer(grpcServer, notifHandler)

	log.Info("Notification gRPC server listening", zap.String("port", cfg.NotificationPort))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}
