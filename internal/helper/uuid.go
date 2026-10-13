package helper

import (
	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.Must(uuid.NewV7()).String()
}
