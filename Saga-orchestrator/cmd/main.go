package main

import (
	"net"

	"chat-grpc/Saga-orchestrator/internal/client"
	"chat-grpc/Saga-orchestrator/internal/handler"
	"chat-grpc/Saga-orchestrator/internal/repository"
	"chat-grpc/Saga-orchestrator/internal/usecase"
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
		log.Fatal("failed to initialize logger", zap.Error(err))
	}
	defer log.Sync()

	db, err := pkg.NewDbChat(log)
	if err != nil {
		log.Fatal("failed to connect to db", zap.Error(err))
	}
	defer db.Close()

	dbUsers, err := pkg.NewDbUsers(log)
	if err != nil {
		log.Fatal("failed to connect to users db", zap.Error(err))
	}
	defer dbUsers.Close()

	conn, err := grpc.NewClient(cfg.NotificationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to notification service", zap.Error(err))
	}
	defer conn.Close()

	notifClient := client.NewNotificationClient(conn)
	repo := repository.NewSagaRepo(db, dbUsers)
	sagaService := usecase.NewSagaService(repo, notifClient, log)
	sagaHandler := handler.NewSagaHandler(sagaService, log)

	listener, err := net.Listen("tcp", ":"+cfg.SagaPort)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	proto_gen.RegisterSagaServiceServer(grpcServer, sagaHandler)

	log.Info("Saga Orchestrator is running on port " + cfg.SagaPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("failed to serve gRPC", zap.Error(err))
	}
}
