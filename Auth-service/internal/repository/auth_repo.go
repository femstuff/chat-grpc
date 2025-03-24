package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"chat-grpc/Auth-service/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *AuthRepo) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	hashPass, err := hashPassword(password)
	if err != nil {
		return 0, err
	}

	var id int64
	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`

	err = a.db.QueryRow(query, name, email, hashPass, role.StringRole()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	return id, nil
}

func (a *AuthRepo) Login(username, pass string) (string, error) {
	var hashPassword string
	query := `SELECT password_hash FROM users WHERE email = $1`
	err := a.db.QueryRow(query, username).Scan(&hashPassword)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(pass)); err != nil {
		return "", errors.New("incorrect login or password")
	}

	return "refresh_token", nil
}

func (a *AuthRepo) GetUser(id int64) (*entity.User, error) {
	var user entity.User
	var roleStr string

	query := `SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	err := a.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roleStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	user.Role = entity.ParseRole(roleStr)

	return &user, nil
}

func (a *AuthRepo) GetList() ([]*entity.User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM users`
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	var roleStr string

	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &roleStr, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
		user.Role = entity.ParseRole(roleStr)
	}

	return users, nil
}

func (a *AuthRepo) UpdateUser(id int64, name, email string) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	_, err := a.db.Exec(query, name, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := a.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	var roleStr string

	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE name = $1`
	err := a.db.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roleStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	user.Role = entity.ParseRole(roleStr)

	return &user, nil
}
