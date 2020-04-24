package storage

import (
	"time"
)

// UserStorage UserStorage
type UserStorage interface {
	ListUser() ([]User, error)
	GetUser(id string) (User, error)
	CreateUser(u User) error
	UpdateUser(u User) error
	DeleteUser(id string) error
}

// User User
type User struct {
	ID        string
	Username  string
	Email     string
	State     uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}
