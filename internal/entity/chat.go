package entity

import "time"

type TypeChat int

const (
	PrivateChat TypeChat = iota
	PublicChat
)

type Chat struct {
	ID        int64
	Name      string
	Users     []int64
	Type      TypeChat
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (t TypeChat) StringType() string {
	switch t {
	case PrivateChat:
		return "private chat"
	case PublicChat:
		return "public chat"
	default:
		return "unknown"
	}
}
