package main

import (
	"net"
	"time"

	"chat-grpc/Auth-service/internal/handler"
	"chat-grpc/Auth-service/internal/repository"
	"chat-grpc/Auth-service/internal/usecase"
	"chat-grpc/Auth-service/jwt"
	"chat-grpc/pkg"
	"chat-grpc/pkg/logger"
	"chat-grpc/proto_gen"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to init logger", zap.Error(err))
	}
	defer log.Sync()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to serve gRPC server", zap.Error(err))
	}

	db, err := pkg.NewDb(log)
	if err != nil {
		log.Fatal("Database conn failed", zap.Error(err))
	}
	defer db.Close()

	repo := repository.NewAuthRepository(db, log)
	jwt := jwt.NewJWTService("key", 15*time.Minute)
	usecase := usecase.NewAuthService(repo, jwt, log)
	handler := handler.NewAuthHandler(usecase, log)

	grpcServer := grpc.NewServer()
	proto_gen.RegisterAuthServiceServer(grpcServer, handler)

	log.Info("gRPC Auth Service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve: %v", zap.Error(err))
	}
}
