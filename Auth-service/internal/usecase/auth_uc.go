package usecase

import (
	"errors"
	"log"
	"time"

	"chat-grpc/Auth-service/internal/entity"
	"chat-grpc/Auth-service/internal/repository"
)

type AuthService struct {
	repo *repository.AuthRepo
}

func NewAuthService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	return s.repo.CreateUser(name, email, password, role)
}

func (s *AuthService) Login(username, pass string) (string, error) {
	return s.repo.Login(username, pass)
}

func (s *AuthService) GetNewRefreshToken(oldToken string) (string, error) {
	if oldToken == "refresh_token" {
		return "new_refresh_token", nil
	}

	return "", errors.New("invalid refresh token")
}

func (s *AuthService) GetNewAccessToken(refreshToken string) (string, error) {
	if refreshToken == "new_refresh_token" {
		return "access_token", nil
	}

	return "", errors.New("invalid access token")
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
