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
	ID        string    `db:"id"`
	Username  string    `db:"username" gdb_i:"-" gdb_u:"-"`
	Password  string    `db:"password" gdb_i:"-" gdb_u:"-"`
	State     uint32    `db:"state" gdb_i:"-" gdb_u:"-"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
