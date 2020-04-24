package storage

import (
	"errors"
	"log"
	"sync"

	"github.com/douglaszuqueto/go-grpc-user/pkg/util"
)

// UserMemoryStorage UserMemoryStorage
type UserMemoryStorage struct {
	db sync.Map
}

// NewUserMemoryStorage NewUserMemoryStorage
func NewUserMemoryStorage() *UserMemoryStorage {
	log.Println("Storage: UserMemoryStorage")

	return &UserMemoryStorage{
		db: sync.Map{},
	}
}

// ListUser ListUser
func (s *UserMemoryStorage) ListUser() ([]User, error) {
	var l []User
	s.db.Range(func(key interface{}, value interface{}) bool {
		u, ok := value.(User)
		if !ok {
			return false
		}
		l = append(l, u)
		return true
	})

	return l, nil
}

// GetUser GetUser
func (s *UserMemoryStorage) GetUser(id string) (User, error) {
	var l User

	value, ok := s.db.Load(id)
	if !ok {
		return l, errors.New("not found")
	}

	l = value.(User)

	return l, nil
}

// CreateUser CreateUser
func (s *UserMemoryStorage) CreateUser(u User) error {
	u.ID = util.GenerateID()

	s.db.Store(u.ID, u)

	return nil
}

// UpdateUser UpdateUser
func (s *UserMemoryStorage) UpdateUser(u User) error {
	s.db.Store(u.ID, u)

	return nil
}

// DeleteUser DeleteUser
func (s *UserMemoryStorage) DeleteUser(id string) error {
	s.db.Delete(id)

	return nil
}
