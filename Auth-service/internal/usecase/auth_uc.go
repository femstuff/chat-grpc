package usecase

import (
	"errors"

	"chat-grpc/Auth-service/internal"
	"chat-grpc/Auth-service/internal/entity"
	"chat-grpc/Auth-service/internal/repository"
	"chat-grpc/Auth-service/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       *repository.AuthRepo
	jwtService *jwt.JWTService
	log        *zap.Logger
}

func NewAuthService(repo *repository.AuthRepo, jwtService *jwt.JWTService, log *zap.Logger) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
		log:        log,
	}
}

func (s *AuthService) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	if err := internal.Validate(name, email, password); err != nil {
		return 0, err
	}

	existingUser, _ := s.repo.GetUserByUsername(name)
	if existingUser != nil {
		s.log.Warn("Attempt to create an already existing user", zap.String("name", name))
		return 0, errors.New("user with this name already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to hash password", zap.Error(err))
		return 0, errors.New("failed to secure password")
	}

	userID, err := s.repo.CreateUser(name, email, string(hashedPassword), role)
	if err != nil {
		s.log.Error("Failed to create user", zap.Error(err))
		return 0, err
	}

	s.log.Info("User created successfully", zap.Int64("userID", userID), zap.String("email", email))
	return userID, nil
}

func (s *AuthService) Login(email, pass string) (string, error) {
	user, err := s.repo.GetUserByUsernameAndValidatePassword(email, pass)
	if err != nil {
		s.log.Error("Login failed", zap.Error(err))
		return "", err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		s.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", err
	}

	if err := s.repo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		s.log.Error("Failed to save refresh token", zap.Error(err))
		return "", err
	}

	s.log.Info("User logged in successfully", zap.String("email", email))
	return refreshToken, nil
}

func (s *AuthService) GetNewAccessToken(refreshToken string) (string, error) {
	claims, err := s.jwtService.VerifyRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("Invalid refresh token", zap.Error(err))
		return "", errors.New("invalid refresh token")
	}

	if err := s.repo.CheckRefreshToken(claims.UserID, refreshToken); err != nil {
		return "", err
	}

	accessToken, err := s.jwtService.GenerateAccessToken(claims.UserID, entity.ParseRole(claims.Role))
	if err != nil {
		s.log.Error("Failed to generate access token", zap.Error(err))
		return "", err
	}

	s.log.Info("Access token refreshed", zap.Int64("userID", claims.UserID))
	return accessToken, nil
}

func (s *AuthService) GetNewRefreshToken(oldToken string) (string, error) {
	claims, err := s.jwtService.VerifyRefreshToken(oldToken)
	if err != nil {
		s.log.Warn("Invalid refresh token", zap.Error(err))
		return "", errors.New("invalid refresh token")
	}

	if err := s.repo.CheckRefreshToken(claims.UserID, oldToken); err != nil {
		return "", err
	}

	if err := s.repo.DeleteRefreshToken(claims.UserID); err != nil {
		return "", err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(claims.UserID, entity.ParseRole(claims.Role))
	if err != nil {
		s.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", err
	}

	if err := s.repo.SaveRefreshToken(claims.UserID, newRefreshToken); err != nil {
		s.log.Error("Failed to save new refresh token", zap.Error(err))
		return "", err
	}

	s.log.Info("Refresh token updated", zap.Int64("userID", claims.UserID))
	return newRefreshToken, nil
}

func (s *AuthService) GetUser(id int64) (*entity.User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		s.log.Warn("User not found", zap.Int64("userID", id))
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetUsers() ([]*entity.User, error) {
	users, err := s.repo.GetList()
	if err != nil {
		s.log.Error("Failed to retrieve users", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (s *AuthService) UpdateUser(id int64, name, email string) error {
	if name == "" || email == "" {
		return errors.New("field name and/or email cannot be empty")
	}

	err := s.repo.UpdateUser(id, name, email)
	if err != nil {
		s.log.Error("Failed to update user", zap.Int64("userID", id), zap.Error(err))
		return err
	}

	s.log.Info("User updated successfully", zap.Int64("userID", id))
	return nil
}

func (s *AuthService) DeleteUser(id int64) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		s.log.Error("Failed to delete user", zap.Int64("userID", id), zap.Error(err))
		return err
	}

	s.log.Info("User deleted successfully", zap.Int64("userID", id))
	return nil
}

func (s *AuthService) CheckToken(accessToken string) error {
	if accessToken == "" {
		return errors.New("empty token")
	}

	_, err := s.jwtService.VerifyAccessToken(accessToken)
	if err != nil {
		s.log.Warn("Invalid access token", zap.Error(err))
		return errors.New("invalid token")
	}

	s.log.Info("Access token is valid")
	return nil
}

func (s *AuthService) GetChatUsersEmails(chatID int64) ([]string, error) {
	return s.repo.GetChatUsersEmails(chatID)
}
