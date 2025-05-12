package repository

import (
	"database/sql"
	"fmt"
)

type SagaRepo struct {
	db      *sql.DB
	dbUsers *sql.DB
}

func NewSagaRepo(db *sql.DB, dbUsers *sql.DB) *SagaRepo {
	return &SagaRepo{db: db, dbUsers: dbUsers}
}

func (r *SagaRepo) DeleteMessage(messageID int64) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, messageID)
	return err
}

func (r *SagaRepo) SaveMessage(messageID int64, chatID int64, text string, userID int64) error {
	query := `INSERT INTO messages (id, chat_id, text, user_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, messageID, chatID, text, userID)
	return err
}

func (r *SagaRepo) GetUserIdFromEmail(email string) (int64, error) {
	var userID int64
	query := `SELECT id FROM users WHERE email = $1`
	err := r.dbUsers.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user_id for email %s: %w", email, err)
	}
	return userID, nil
}
