package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	// ErrAlreadyExists ErrAlreadyExists
	ErrAlreadyExists = errors.New("object already exists")

	// ErrDoesNotExist ErrDoesNotExist
	ErrDoesNotExist = errors.New("object does not exist")
)

func doInsert(db *sql.DB, query string, args ...interface{}) (string, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var id string

	err = stmt.QueryRow(args...).Scan(&id)
	if err != nil {
		return "", HandlePSQLError(err)
	}

	return id, nil
}

func doUpdate(db *sql.DB, query string, args ...interface{}) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("db.Prepare", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println("stmt.Exec", err)
		return err
	}

	return nil
}

func doDelete(db *sql.DB, table string, id string) error {
	query := fmt.Sprintf(`DELETE FROM public.%s WHERE id = $1`, table)

	stmt, err := db.Prepare(query)
	if err != nil {
		return HandlePSQLError(err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(id)
	if err != nil {
		return HandlePSQLError(err)
	}

	rows, _ := r.RowsAffected()
	if rows == 0 {
		return HandlePSQLError(ErrDoesNotExist)
	}

	return nil
}

// HandlePSQLError handlePSQLError
// https://github.com/brocaar/lora-app-server/blob/master/internal/storage/errors.go#L41
func HandlePSQLError(err error) error {
	if err == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	switch err := err.(type) {
	case *pq.Error:
		switch err.Code.Name() {
		case "unique_violation":
			return errors.Wrap(ErrAlreadyExists, err.Constraint)
		case "foreign_key_violation":
			return ErrDoesNotExist
		default:
			log.Println(err.Code.Name())
			return err
		}
	}

	return err
}
