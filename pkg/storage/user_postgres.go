package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// UserPostgresStorage UserPostgresStorage
type UserPostgresStorage struct {
	db *sql.DB
}

// NewUserPostgresStorage NewUserPostgresStorage
func NewUserPostgresStorage() *UserPostgresStorage {
	log.Println("Storage: UserPostgresStorage")

	db := &UserPostgresStorage{}
	db.connect()

	return db
}

// ListUser ListUser
func (s *UserPostgresStorage) ListUser() ([]User, error) {
	var l []User

	query := `
		SELECT
			id,
			username,
			state,
			created_at,
			updated_at
		FROM
			public.user
		ORDER BY
			created_at
		`

	rows, err := s.db.Query(query)
	if err != nil {
		return l, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.State,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return l, err
		}

		l = append(l, u)
	}

	return l, nil
}

// GetUser GetUser
func (s *UserPostgresStorage) GetUser(id string) (User, error) {
	var u User

	query := `
		SELECT
			id,
			username,
			state,
			created_at,
			updated_at
		FROM
			public.user
		WHERE
			id = $1
			`

	err := s.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.State,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, HandlePSQLError(err)
	}

	return u, nil
}

// CreateUser CreateUser
func (s *UserPostgresStorage) CreateUser(u User) (string, error) {
	query := `
		INSERT INTO public.user 
			(username, state) 
		VALUES 
			($1, $2)
		RETURNING id`

	id, err := doInsert(s.db, query, u.Username, u.State)
	return id, err
}

// UpdateUser UpdateUser
func (s *UserPostgresStorage) UpdateUser(u User) error {
	query := `
		UPDATE
			public.user
		SET
			username = $1,
			state = $2
		WHERE
			id = $3`

	return doUpdate(s.db, query, u.Username, u.State, u.ID)
}

// DeleteUser DeleteUser
func (s *UserPostgresStorage) DeleteUser(id string) error {
	return doDelete(s.db, "user", id)
}

func (s *UserPostgresStorage) connect() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbDatabase := os.Getenv("DB_DATABASE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbDatabase,
	)

	var err error
	s.db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	s.db.SetMaxIdleConns(5)
	s.db.SetMaxOpenConns(5)
	s.db.SetConnMaxLifetime(5 * time.Minute)

	err = s.db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}
