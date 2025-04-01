package repository

import (
	"database/sql"
	"fmt"
	"time"

	"chat-grpc/Chat-service/internal/entity"
	"go.uber.org/zap"
)

type ChatRepo interface {
	CreateChat(usernames []string) (int64, error)
	DeleteChat(id int64) error
	SendMessage(chatID int64, from, text string, timestamp time.Time) error
	GetMessages(chatID int64) ([]entity.Message, error)
}

type chatRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewChatRepository(db *sql.DB, log *zap.Logger) ChatRepo {
	return &chatRepository{db: db, log: log}
}

func (r *chatRepository) CreateChat(usernames []string) (int64, error) {
	r.log.Info("Creating chat")

	var chatID int64
	query := `INSERT INTO chats (created_at, updated_at) VALUES (NOW(), NOW()) RETURNING id`
	err := r.db.QueryRow(query).Scan(&chatID)
	if err != nil {
		r.log.Error("Failed to create chat", zap.Error(err))
		return 0, err
	}

	for _, username := range usernames {
		var userID int64
		err := r.db.QueryRow("SELECT id FROM users WHERE name = $1", username).Scan(&userID)
		if err != nil {
			r.log.Error("User not found", zap.String("username", username), zap.Error(err))
			return 0, fmt.Errorf("user %s not found: %w", username, err)
		}

		_, err = r.db.Exec("INSERT INTO chat_users (chat_id, user_id) VALUES ($1, $2)", chatID, userID)
		if err != nil {
			r.log.Error("Failed to add user to chat", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID), zap.Error(err))
			return 0, err
		}
	}

	r.log.Info("Chat created successfully", zap.Int64("chat_id", chatID))
	return chatID, nil
}

func (r *chatRepository) DeleteChat(id int64) error {
	r.log.Info("Deleting chat", zap.Int64("chat_id", id))

	query := "DELETE FROM chats WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		r.log.Error("Failed to delete chat", zap.Error(err))
		return err
	}

	r.log.Info("Chat deleted", zap.Int64("chat_id", id))
	return nil
}

func (r *chatRepository) SendMessage(chatID int64, username, text string, timestamp time.Time) error {
	r.log.Info("Sending message", zap.Int64("chat_id", chatID), zap.String("username", username))

	var userID int64
	err := r.db.QueryRow("SELECT id FROM users WHERE name = $1", username).Scan(&userID)
	if err != nil {
		r.log.Error("User not found", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("user %s not found: %w", username, err)
	}

	query := `INSERT INTO messages (chat_id, user_id, text, timestamp) VALUES ($1, $2, $3, $4)`
	_, err = r.db.Exec(query, chatID, userID, text, timestamp)
	if err != nil {
		r.log.Error("Failed to send message", zap.Error(err))
		return err
	}

	r.log.Info("Message sent successfully", zap.Int64("chat_id", chatID), zap.String("username", username))
	return nil
}

func (r *chatRepository) GetMessages(chatID int64) ([]entity.Message, error) {
	query := "SELECT user_id, text, timestamp FROM messages WHERE chat_id = $1 ORDER BY timestamp"
	rows, err := r.db.Query(query, chatID)
	if err != nil {
		r.log.Error("Failed to get messages", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		err := rows.Scan(&msg.Sender, &msg.Content, &msg.CreatedAt)
		if err != nil {
			r.log.Error("Failed to scan message", zap.Error(err))
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
