package utils

import (
	"time"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
