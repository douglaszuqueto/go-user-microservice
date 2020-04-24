package util

import (
	"github.com/google/uuid"
)

// GenerateID GenerateID
func GenerateID() string {
	id := uuid.New()

	return id.String()
}
