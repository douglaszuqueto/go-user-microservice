package util

import (
	"os"

	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"
)

// GetStorageType GetStorageType
func GetStorageType() storage.UserStorage {
	storageType := os.Getenv("APP_STORAGE")

	var db storage.UserStorage

	switch storageType {
	case "memory":
		db = storage.NewUserMemoryStorage()
	case "postgres":
		db = storage.NewUserPostgresStorage()
	default:
		panic("unknow storage type: " + storageType)
	}

	return db
}
