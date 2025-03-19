package entity

import (
	"errors"
	"time"
)

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

func StringType(s string) (TypeChat, error) {
	switch s {
	case "private":
		return PrivateChat, nil
	case "public":
		return PublicChat, nil
	default:
		return 0, errors.New("invalid type chat")
	}
}
