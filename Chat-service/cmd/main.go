package main

import (
	"net"

	"chat-grpc/Auth-service/interceptor"
	"chat-grpc/Chat-service/internal/broker"
	"chat-grpc/Chat-service/internal/handler"
	"chat-grpc/Chat-service/internal/repository"
	"chat-grpc/Chat-service/internal/usecase"
	"chat-grpc/pkg"
	"chat-grpc/pkg/config"
	"chat-grpc/pkg/logger"
	"chat-grpc/proto_gen"
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

	db, err := pkg.NewDbChat(log)
	if err != nil {
		log.Fatal("Database conn failed", zap.Error(err))
	}
	defer db.Close()

	dbUsers, err := pkg.NewDbUsers(log)
	if err != nil {
		log.Fatal("Database conn users failed", zap.Error(err))
	}
	defer dbUsers.Close()

	broker, err := broker.NewNatsBroker(cfg.NatsUrl)
	if err != nil {
		log.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer broker.Close()

	chatRepo := repository.NewChatRepository(db, dbUsers, log)
	chatUseCase := usecase.NewChatUseCase(chatRepo, log, broker)
	chatHandler := handler.NewChatService(chatUseCase, log)

	listener, err := net.Listen("tcp", ":"+cfg.ServerPortChat)
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	conn, err := grpc.NewClient(cfg.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	defer conn.Close()

	log.Info("Successfully connected to auth service")

	authClient := interceptor.NewAuthClient(conn)
	authInterceptor := interceptor.NewAuthInterceptor(authClient, log)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	proto_gen.RegisterChatServiceServer(grpcServer, chatHandler)

	log.Info("Chat Service is running on ", zap.String("port", cfg.ServerPortChat))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve", zap.Error(err))
	}
}
