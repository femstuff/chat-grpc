package entity

import "time"

type Message struct {
	ID        int
	ChatID    int
	Sender    int
	Content   string
	CreatedAt time.Time
	UpdateAt  time.Time
}
