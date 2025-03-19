package main

import (
	"log"
	"net"

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
	usecase := usecase.NewAuthService(repo)
	handler := handler.NewAuthHandler(usecase)

	grpcServer := grpc.NewServer()
	proto_gen.RegisterAuthServiceServer(grpcServer, handler)

	log.Println("gRPC Auth Service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
