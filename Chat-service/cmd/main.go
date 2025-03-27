package main

import (
	"net"

	"chat-grpc/Chat-service/internal/handler"
	"chat-grpc/Chat-service/internal/repository"
	"chat-grpc/Chat-service/internal/usecase"
	"chat-grpc/pkg/logger"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to init logger", zap.Error(err))
	}
	defer log.Sync()

	chatRepo := repository.NewChatRepository(log)
	chatUseCase := usecase.NewChatUseCase(chatRepo, log)
	chatHandler := handler.NewChatService(chatUseCase, log)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Failed to listen: %v", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	proto_gen.RegisterChatServiceServer(grpcServer, chatHandler)

	log.Info("Chat Service is running on port 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve: %v", zap.Error(err))
	}
}
