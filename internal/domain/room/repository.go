package room

import "context"

// Repository defines the interface for room persistence
type Repository interface {
	// Save persists a room
	Save(ctx context.Context, room *Room) error

	// FindByID retrieves a room by ID
	FindByID(ctx context.Context, id RoomID) (*Room, error)

	// FindByCode retrieves a room by code
	FindByCode(ctx context.Context, code RoomCode) (*Room, error)

	// Delete removes a room
	Delete(ctx context.Context, id RoomID) error
}
