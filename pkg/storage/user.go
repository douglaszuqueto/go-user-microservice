package storage

import (
	"errors"
	"sync"
	"time"
)

var db = sync.Map{}
var rw = sync.Mutex{}

// User User
type User struct {
	ID        string
	Username  string
	Email     string
	State     uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListUser ListUser
func ListUser() ([]User, error) {
	rw.Lock()
	defer rw.Unlock()

	var l []User
	db.Range(func(key interface{}, value interface{}) bool {
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
func GetUser(id string) (User, error) {
	rw.Lock()
	defer rw.Unlock()

	var l User

	value, ok := db.Load(id)
	if !ok {
		return l, errors.New("not found")
	}

	l = value.(User)

	return l, nil
}

// CreateUser CreateUser
func CreateUser(u User) error {
	rw.Lock()
	defer rw.Unlock()

	db.Store(u.ID, u)

	return nil
}

// UpdateUser UpdateUser
func UpdateUser(u User) error {
	rw.Lock()
	defer rw.Unlock()

	db.Store(u.ID, u)

	return nil
}

// DeleteUser DeleteUser
func DeleteUser(id string) error {
	rw.Lock()
	defer rw.Unlock()

	db.Delete(id)

	return nil
}
