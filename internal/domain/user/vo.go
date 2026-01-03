package user

import (
	"errors"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

// UserID represents a user identifier
type UserID struct {
	value string
}

// NewUserID creates a new UserID
func NewUserID() UserID {
	return UserID{value: utils.GenerateUUID()}
}

// NewUserIDFromString creates a UserID from a string
func NewUserIDFromString(value string) (UserID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return UserID{}, err
	}
	return UserID{value: value}, nil
}

// String returns the string representation of UserID
func (id UserID) String() string {
	return id.value
}

// Equals checks if two UserIDs are equal
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// UserName represents a user name
type UserName struct {
	value string
}

// NewUserName creates a new UserName
func NewUserName(value string) (UserName, error) {
	if value == "" {
		return UserName{}, errors.New("user name cannot be empty")
	}
	if len(value) > 255 {
		return UserName{}, errors.New("user name is too long")
	}
	return UserName{value: value}, nil
}

// String returns the string representation of UserName
func (n UserName) String() string {
	return n.value
}

// Equals checks if two UserNames are equal
func (n UserName) Equals(other UserName) bool {
	return n.value == other.value
}
