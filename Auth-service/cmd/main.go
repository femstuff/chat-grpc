package main

import (
	"net"
	"time"

	"chat-grpc/Auth-service/internal/handler"
	"chat-grpc/Auth-service/internal/repository"
	"chat-grpc/Auth-service/internal/usecase"
	"chat-grpc/Auth-service/jwt"
	"chat-grpc/pkg"
	"chat-grpc/pkg/config"
	"chat-grpc/pkg/logger"
	"chat-grpc/proto_gen"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to init logger", zap.Error(err))
	}
	defer log.Sync()

	listener, err := net.Listen("tcp", ":"+cfg.ServerPortAuth)
	if err != nil {
		log.Fatal("Failed to serve gRPC server", zap.Error(err))
	}

	dbUser, err := pkg.NewDbUsers(log)
	if err != nil {
		log.Fatal("Database conn failed", zap.Error(err))
	}
	defer dbUser.Close()
	dbAuth, err := pkg.NewDbChat(log)
	if err != nil {
		log.Fatal("database chat conn failes", zap.Error(err))
	}

	repo := repository.NewAuthRepository(dbUser, dbAuth, log)
	jwt := jwt.NewJWTService("key", 15*time.Minute, log)
	usecase := usecase.NewAuthService(repo, jwt, log)
	handler := handler.NewAuthHandler(usecase, log)

	grpcServer := grpc.NewServer()
	proto_gen.RegisterAuthServiceServer(grpcServer, handler)

	log.Info("gRPC Auth Service is running on ", zap.String("port", cfg.ServerPortAuth))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve: %v", zap.Error(err))
	}
}
