package user

import "context"

// Repository defines the interface for user persistence
type Repository interface {
	// Save persists a user
	Save(ctx context.Context, user *User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id UserID) (*User, error)

	// Delete removes a user
	Delete(ctx context.Context, id UserID) error
}
