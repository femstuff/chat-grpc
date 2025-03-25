package usecase

import (
	"errors"

	"chat-grpc/Auth-service/internal"
	"chat-grpc/Auth-service/internal/entity"
	"chat-grpc/Auth-service/internal/repository"
	"chat-grpc/Auth-service/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.AuthRepo
	jwtService *jwt.JWTService
	log        *zap.Logger
}

func NewAuthService(repo *repository.AuthRepo, jwtService *jwt.JWTService, log *zap.Logger) *AuthService {
	return &AuthService{repo: repo, jwtService: jwtService, log: log}
}

func (s *AuthService) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	if err := internal.Validate(name, email, password); err != nil {
		return 0, err
	}

	return s.repo.CreateUser(name, email, password, role)
}

func (s *AuthService) Login(username, pass string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		s.log.Error("Failed to get user by username", zap.Error(err))
		return "", errors.New("incorrect login or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		s.log.Warn("Incorrect password attempt", zap.String("username", username))
		return "", errors.New("incorrect login or password")
	}

	refreshToken, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		s.log.Error("Failed to generate token", zap.Error(err))
		return "", err
	}

	s.log.Info("User logged in successfully", zap.String("username", username))
	return refreshToken, nil
}

func (s *AuthService) GetNewRefreshToken(oldToken string) (string, error) {
	claims, err := s.jwtService.VerifyToken(oldToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	newRefreshToken, err := s.jwtService.GenerateToken(claims.UserID, entity.ParseRole(claims.Role))
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil
}

func (s *AuthService) GetNewAccessToken(refreshToken string) (string, error) {
	claims, err := s.jwtService.VerifyToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	accessToken, err := s.jwtService.GenerateToken(claims.UserID, entity.ParseRole(claims.Role))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) GetUser(id int64) (*entity.User, error) {
	return s.repo.GetUser(id)
}

func (s *AuthService) GetUsers() ([]*entity.User, error) {
	return s.repo.GetList()
}

func (s *AuthService) UpdateUser(id int64, name, email string) error {
	if name == "" || email == "" {
		return errors.New("field name and/or email cannot be empty")
	}

	return s.repo.UpdateUser(id, name, email)
}

func (s *AuthService) DeleteUser(id int64) error {
	return s.repo.DeleteUser(id)
}

func (s *AuthService) CheckToken(endpoint string) error {
	if endpoint == "" {
		return errors.New("empty endpoint")
	}

	return nil
}
