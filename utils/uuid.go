package utils

import (
	"errors"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// ValidateUUID validates if a string is a valid UUID
func ValidateUUID(id string) error {
	if id == "" {
		return errors.New("uuid cannot be empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid uuid format")
	}

	return nil
}

// ParseUUID parses a string to UUID and returns error if invalid
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// MustParseUUID parses a string to UUID and panics if invalid
// Use this only when you are certain the input is valid
func MustParseUUID(id string) uuid.UUID {
	return uuid.MustParse(id)
}
