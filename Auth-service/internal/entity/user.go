package entity

import "time"

type Role int

const (
	UserRole Role = iota
	AdminRole
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
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

func ParseRole(roleStr string) Role {
	switch roleStr {
	case "admin":
		return AdminRole
	case "user":
		return UserRole
	default:
		return UserRole
	}
}
