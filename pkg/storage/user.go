package storage

import (
	"os"
	"time"
)

// UserStorage UserStorage
type UserStorage interface {
	ListUser() ([]User, error)
	GetUser(id string) (User, error)
	CreateUser(u User) (string, error)
	UpdateUser(u User) error
	DeleteUser(id string) error
}

// User User
type User struct {
	ID        string
	Username  string
	State     uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetStorageType GetStorageType
func GetStorageType() UserStorage {
	storageType := os.Getenv("APP_STORAGE")

	var db UserStorage

	switch storageType {
	case "memory":
		db = NewUserMemoryStorage()
	case "postgres":
		db = NewUserPostgresStorage()
	default:
		panic("unknow storage type: " + storageType)
	}

	return db
}
