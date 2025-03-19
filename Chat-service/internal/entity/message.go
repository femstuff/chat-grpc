package entity

import "time"

type Message struct {
	ID        int64
	ChatID    int64
	Sender    int64
	Content   string
	CreatedAt time.Time
	UpdateAt  time.Time
}
