package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "net/url"
	"strings"
	"time"

	"chat-grpc/proto_gen"
	"github.com/otiai10/opengraph/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatRepo interface {
	CreateChat(usernames []string) (int64, error)
	DeleteChat(id int64) error
	SendMessage(chatID int64, from, text string, timestamp time.Time) error
	GetMessagesByChatID(ctx context.Context, chatID int64) ([]*proto_gen.Message, error)
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
	// text = https://github.com/ split : [0]
	checkUrl := strings.Split(text, "://")
	if checkUrl[0] == "https" {
		//resp, err := http.Get(text)
		//if err != nil {
		//	r.log.Error("Failed to fetch URL", zap.String("url", text), zap.Error(err))
		//	return err
		//}
		//defer resp.Body.Close()

		//parsedUrl, err := url.Parse(text)
		//if err != nil {
		//	r.log.Error("Failed to fetch URL", zap.String("url", text), zap.Error(err))
		//	return err
		//}
		//
		//article, err := readability.FromReader(resp.Body, parsedUrl)
		//if err != nil {
		//	r.log.Error("Failed to extract metadata from URL", zap.String("url", text), zap.Error(err))
		//	return err
		//}
		//
		//shortContent := article.TextContent
		//if len(article.TextContent) >= 200 {
		//	shortContent = shortContent[:200] + "..."
		//}
		//
		//text += fmt.Sprintf("\n\nSite name: %s\n\nTitle: %s\n\nDescription: %s\n\nImage: %s\n\nFavicon: %s\n\nURL: [%s]",
		//	article.SiteName, article.Title, shortContent, article.Image, article.Favicon, text)
		ogp, err := opengraph.Fetch(text)
		if err != nil {
			fmt.Println("Ошибка при получении Open Graph данных:", err)
			return err
		}

		text += fmt.Sprintf("\n\nTitle: %s\n\nDescription: %s\n\nImage: %v\n\nURL: %s",
			ogp.Title, ogp.Description, ogp.Image, ogp.URL)
	}

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

func (r *chatRepository) GetMessagesByChatID(ctx context.Context, chatID int64) ([]*proto_gen.Message, error) {
	query := `SELECT u.name, m.text, m.timestamp 
			  FROM messages m
			  JOIN users u ON m.user_id = u.id
			  WHERE m.chat_id = $1
			  ORDER BY m.timestamp ASC`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		r.log.Error("Failed to fetch chat history", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var messages []*proto_gen.Message
	for rows.Next() {
		var username, text string
		var timestamp time.Time

		if err := rows.Scan(&username, &text, &timestamp); err != nil {
			r.log.Error("Failed to scan message row", zap.Error(err))
			return nil, err
		}

		messages = append(messages, &proto_gen.Message{
			ChatId:    chatID,
			From:      username,
			Text:      text,
			Timestamp: timestamppb.New(timestamp),
		})
	}

	if err := rows.Err(); err != nil {
		r.log.Error("Row iteration error", zap.Error(err))
		return nil, err
	}

	return messages, nil
}
