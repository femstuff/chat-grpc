package repository

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"chat-grpc/Auth-service/interceptor"
	"chat-grpc/Auth-service/internal/entity"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	dbUser *sql.DB
	dbAuth *sql.DB
	log    *zap.Logger
	client *interceptor.AuthClient
}

func NewAuthRepository(dbUser *sql.DB, dbAuth *sql.DB, log *zap.Logger) *AuthRepo {
	return &AuthRepo{dbUser: dbUser, dbAuth: dbAuth, log: log}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (a *AuthRepo) SaveRefreshToken(userID int64, token string) error {
	hashedToken := hashToken(token)

	_, err := a.dbUser.Exec(`
		INSERT INTO refresh_tokens (user_id, token) 
		VALUES ($1, $2) 
		ON CONFLICT (user_id) DO UPDATE SET token = $2`,
		userID, hashedToken,
	)

	if err != nil {
		a.log.Error("Failed to save refresh token", zap.Error(err))
	}

	a.log.Info("Success save refresh token")
	return err
}

func hashPassword(password string, log *zap.Logger) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed with hash password", zap.Error(err))
		return "", err
	}

	return string(hash), nil
}

func (a *AuthRepo) CreateUser(name, email, password string, role entity.Role) (int64, error) {
	a.log.Info("Creating user", zap.String("name", name), zap.String("email", email), zap.String("role", role.StringRole()))

	var id int64
	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`

	err := a.dbUser.QueryRow(query, name, email, password, role.StringRole()).Scan(&id)
	if err != nil {
		a.log.Error("Failed to insert new user in dbUser", zap.Error(err))
		return 0, err
	}

	a.log.Info("Successful create user", zap.Int64("userid:", id))
	return id, nil
}

func (a *AuthRepo) Login(username, pass string) (string, error) {
	a.log.Info("Login attempt", zap.String("username", username))

	var hashPassword string
	query := `SELECT password_hash FROM users WHERE name = $1`
	err := a.dbUser.QueryRow(query, username).Scan(&hashPassword)
	if err != nil {
		a.log.Error("Failed to scan user", zap.Error(err))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(pass)); err != nil {
		a.log.Error("Auth fail", zap.String("username", username))
		return "", errors.New("incorrect login or password")
	}

	a.log.Info("Success login", zap.String("username", username))
	return "refresh_token", nil
}

func (a *AuthRepo) GetUser(id int64) (*entity.User, error) {
	var user entity.User
	var roleStr string

	query := `SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	err := a.dbUser.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roleStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		a.log.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	user.Role = entity.ParseRole(roleStr)

	a.log.Info("User found", zap.Int64("id", id))
	return &user, nil
}

func (a *AuthRepo) GetList() ([]*entity.User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM users`
	rows, err := a.dbUser.Query(query)
	if err != nil {
		a.log.Error("Failed to get list users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	var roleStr string

	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &roleStr, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			a.log.Error("Failed with scan user", zap.Error(err))
			return nil, err
		}

		users = append(users, user)
		user.Role = entity.ParseRole(roleStr)
	}

	a.log.Info("Get list users", zap.Int("count", len(users)))
	return users, nil
}

func (a *AuthRepo) UpdateUser(id int64, name, email string) error {
	a.log.Info("Update user", zap.Int64("id", id))

	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	_, err := a.dbUser.Exec(query, name, email, id)
	if err != nil {
		a.log.Error("Failed with update user", zap.Error(err))
		return err
	}

	a.log.Info("Success update user", zap.Int64("id", id))
	return nil
}

func (a *AuthRepo) DeleteUser(id int64) error {
	a.log.Info("Delete user", zap.Int64("id", id))

	query := `DELETE FROM users WHERE id = $1`
	_, err := a.dbUser.Exec(query, id)
	if err != nil {
		a.log.Error("Failed with delete user", zap.Error(err))
		return err
	}

	a.log.Info("Success delete user", zap.Int64("id", id))
	return nil
}

func (a *AuthRepo) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	var roleStr string

	query := `SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE name = $1`
	err := a.dbUser.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &roleStr, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		a.log.Error("Failed to scan user", zap.Error(err))
		return nil, err
	}

	user.Role = entity.ParseRole(roleStr)

	a.log.Info("Success get user by username", zap.String("username", username))
	return &user, nil
}

func (a *AuthRepo) GetUserByUsernameAndValidatePassword(email, pass string) (*entity.User, error) {
	a.log.Info("Login attempt", zap.String("email", email))

	var user entity.User
	var roleStr string

	query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`
	err := a.dbUser.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &roleStr)
	if err != nil {
		a.log.Error("Failed to scan user", zap.Error(err))
		return nil, errors.New("incorrect login or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		a.log.Error("Auth fail", zap.String("email", email))
		return nil, errors.New("incorrect login or password")
	}

	user.Role = entity.ParseRole(roleStr)

	return &user, nil
}

func (a *AuthRepo) CheckRefreshToken(userID int64, token string) error {
	hashedToken := hashToken(token)

	var dbToken string
	err := a.dbUser.QueryRow(`SELECT token FROM refresh_tokens WHERE user_id = $1`, userID).Scan(&dbToken)
	if err != nil {
		a.log.Warn("Refresh token not found", zap.Error(err))
		return errors.New("invalid refresh token")
	}

	if hashedToken != dbToken {
		a.log.Warn("Refresh token mismatch")
		return errors.New("invalid refresh token")
	}

	return nil
}

func (a *AuthRepo) DeleteRefreshToken(userID int64) error {
	_, err := a.dbUser.Exec(`DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		a.log.Error("Failed to delete refresh token", zap.Error(err))
		return err
	}

	a.log.Info("Refresh token deleted", zap.Int64("userID", userID))
	return nil
}

func (a *AuthRepo) GetEmailByUserID(userID int64) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = $1`
	err := a.dbUser.QueryRow(query, userID).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("error fetching email")
	}
	return email, nil
}

func (a *AuthRepo) GetChatUsersEmails(chatID int64) ([]string, error) {
	var usersId []int64
	var emails []string
	query := `SELECT user_id FROM chat_users WHERE chat_id = $1`
	rowsUsers, err := a.dbAuth.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed ti exec query: %w", err)
	}
	defer rowsUsers.Close()

	for rowsUsers.Next() {
		var userId int64
		if err := rowsUsers.Scan(&userId); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usersId = append(usersId, userId)
	}

	a.log.Info("users len", zap.Int("len", len(usersId)))

	for _, userId := range usersId {
		email, err := a.GetEmailByUserID(userId)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	fmt.Println("emails:", emails)
	return emails, nil
}
