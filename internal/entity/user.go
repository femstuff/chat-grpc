package entity

import (
	"time"
)

type Role int

const (
	UserRole Role = iota
	AdminRole
)

type User struct {
	ID        int
	Name      string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Role) StringRole() string {
	switch r {
	case UserRole:
		return "user"
	case AdminRole:
		return "admin"
	default:
		return "unknown role"
	}
}
