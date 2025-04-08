package handler

import (
	"context"

	"chat-grpc/Auth-service/internal/entity"
	"chat-grpc/Auth-service/internal/usecase"
	"chat-grpc/proto_gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	proto_gen.UnimplementedAuthServiceServer
	usecase *usecase.AuthService
	log     *zap.Logger
}

func NewAuthHandler(uc *usecase.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{usecase: uc, log: log}
}

func (h *AuthHandler) Create(ctx context.Context, req *proto_gen.CreateUserRequest) (*proto_gen.CreateUserResponse, error) {
	id, err := h.usecase.CreateUser(req.Name, req.Email, req.Password, entity.Role(req.Role))
	if err != nil {
		return nil, err
	}

	return &proto_gen.CreateUserResponse{Id: id}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *proto_gen.LoginRequest) (*proto_gen.LoginResponse, error) {
	refreshToken, err := h.usecase.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto_gen.LoginResponse{RefreshToken: refreshToken}, nil
}

func (h *AuthHandler) GetRefreshToken(ctx context.Context, req *proto_gen.RefreshTokenRequest) (*proto_gen.RefreshTokenResponse, error) {
	refreshToken, err := h.usecase.GetNewRefreshToken(req.OldRefreshToken)
	if err != nil {
		return nil, err
	}

	return &proto_gen.RefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (h *AuthHandler) GetAccessToken(ctx context.Context, req *proto_gen.AccessTokenRequest) (*proto_gen.AccessTokenResponse, error) {
	accessToken, err := h.usecase.GetNewAccessToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &proto_gen.AccessTokenResponse{AccessToken: accessToken}, nil
}

func (h *AuthHandler) Get(ctx context.Context, req *proto_gen.GetUserRequest) (*proto_gen.GetUserResponse, error) {
	user, err := h.usecase.GetUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto_gen.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      proto_gen.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *AuthHandler) GetList(ctx context.Context, req *proto_gen.AuthEmpty) (*proto_gen.GetListResponse, error) {
	users, err := h.usecase.GetUsers()
	if err != nil {
		return nil, err
	}

	var usersProto []*proto_gen.GetUserResponse
	for _, user := range users {
		usersProto = append(usersProto, &proto_gen.GetUserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      proto_gen.Role(user.Role),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &proto_gen.GetListResponse{Users: usersProto}, nil
}

func (h *AuthHandler) Update(ctx context.Context, req *proto_gen.UpdateUserRequest) (*proto_gen.AuthEmpty, error) {
	err := h.usecase.UpdateUser(req.Id, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	return &proto_gen.AuthEmpty{}, nil
}

func (h *AuthHandler) Delete(ctx context.Context, req *proto_gen.DeleteUserRequest) (*proto_gen.AuthEmpty, error) {
	err := h.usecase.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &proto_gen.AuthEmpty{}, nil
}

func (h *AuthHandler) Check(ctx context.Context, req *proto_gen.CheckAccessRequest) (*proto_gen.AuthEmpty, error) {
	err := h.usecase.CheckToken(req.EndpointAddress)
	if err != nil {
		return nil, err
	}

	return &proto_gen.AuthEmpty{}, nil
}

func (h *AuthHandler) CheckToken(ctx context.Context, req *proto_gen.CheckTokenRequest) (*proto_gen.AuthEmpty, error) {
	err := h.usecase.CheckToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &proto_gen.AuthEmpty{}, nil
}
