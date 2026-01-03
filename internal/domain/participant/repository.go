package participant

import "context"

// Repository defines the interface for participant persistence
type Repository interface {
	// Save persists a participant
	Save(ctx context.Context, participant *Participant) error

	// FindByID retrieves a participant by ID
	FindByID(ctx context.Context, id ParticipantID) (*Participant, error)

	// FindByRoomID retrieves all participants in a room
	FindByRoomID(ctx context.Context, roomID RoomID) ([]*Participant, error)

	// FindByRoomAndUser retrieves a specific participant by room and user
	FindByRoomAndUser(ctx context.Context, roomID RoomID, userID UserID) (*Participant, error)

	// Delete removes a participant
	Delete(ctx context.Context, roomID RoomID, userID UserID) error
}
