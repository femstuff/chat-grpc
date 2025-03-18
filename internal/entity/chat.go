package entity

import "time"

type TypeChat int

const (
	PrivateChat TypeChat = iota
	PublicChat
)

type Chat struct {
	ID        int
	Name      string
	Users     []int
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
