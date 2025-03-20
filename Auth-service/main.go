package main

import (
	"log"
	"net"
	"time"

	"chat-grpc/Auth-service/internal"
	"chat-grpc/Auth-service/internal/handler"
	"chat-grpc/Auth-service/internal/repository"
	"chat-grpc/Auth-service/internal/usecase"
	"chat-grpc/proto_gen"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	repo := repository.NewAuthRepository()
	jwt := internal.NewJWTService("key", 15*time.Minute)
	usecase := usecase.NewAuthService(repo, jwt)
	handler := handler.NewAuthHandler(usecase)

	grpcServer := grpc.NewServer()
	proto_gen.RegisterAuthServiceServer(grpcServer, handler)

	log.Println("gRPC Auth Service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
