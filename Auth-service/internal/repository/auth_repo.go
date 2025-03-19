package repository

import (
	"errors"
	"log"
	"time"

	"chat-grpc/Auth-service/internal/entity"
)

type AuthRepo struct {
	users map[int64]*entity.User
	ID    int64
}

func NewAuthRepository() *AuthRepo {
	return &AuthRepo{
		users: make(map[int64]*entity.User),
		ID:    1,
	}
}

func (a *AuthRepo) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	id := a.ID
	a.users[id] = &entity.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return id, nil
}

func (a *AuthRepo) Login(username, pass string) (string, error) {
	for _, user := range a.users {
		if username == user.Name && pass == user.Password {
			return "refresh_token", nil
		}
	}

	return "", errors.New("incorrect login or password")
}

func (a *AuthRepo) GetUser(id int64) (*entity.User, error) {
	user, exists := a.users[id]
	log.Print("repo layer\n")
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *AuthRepo) GetList() ([]*entity.User, error) {
	var users []*entity.User
	for _, user := range a.users {
		users = append(users, user)
	}

	return users, nil
}

func (a *AuthRepo) UpdateUser(user *entity.User) error {
	_, exists := a.users[user.ID]
	if !exists {
		return errors.New("user not found")
	}

	a.users[user.ID] = user
	return nil
}

func (a *AuthRepo) DeleteUser(id int64) error {
	_, exists := a.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(a.users, id)
	return nil
}
