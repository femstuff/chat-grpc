package usecase

import (
	"errors"
	"log"
	"time"

	"chat-grpc/Auth-service/internal"
	"chat-grpc/Auth-service/internal/entity"
	"chat-grpc/Auth-service/internal/repository"
)

type AuthService struct {
	repo       *repository.AuthRepo
	jwtService *internal.JWTService
}

func NewAuthService(repo *repository.AuthRepo, jwtService *internal.JWTService) *AuthService {
	return &AuthService{repo: repo, jwtService: jwtService}
}

func (s *AuthService) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	return s.repo.CreateUser(name, email, password, role)
}

func (s *AuthService) Login(username, pass string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil || user.Password != pass {
		return "", errors.New("incorrect login or password")
	}

	refreshToken, err := s.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

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
	log.Print("uc layer\n")
	return s.repo.GetList()
}

func (s *AuthService) UpdateUser(id int64, name, email string) error {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return err
	}

	user.Name = name
	user.Email = email
	user.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateUser(user)
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
