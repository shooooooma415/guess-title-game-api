package room_emoji

import "context"

// Repository defines the interface for room emoji persistence
type Repository interface {
	// Save persists a room emoji
	Save(ctx context.Context, emoji *RoomEmoji) error

	// FindByID retrieves a room emoji by ID
	FindByID(ctx context.Context, id RoomEmojiID) (*RoomEmoji, error)

	// FindByRoomID retrieves all emojis in a room
	FindByRoomID(ctx context.Context, roomID RoomID) ([]*RoomEmoji, error)

	// Delete removes a room emoji
	Delete(ctx context.Context, id RoomEmojiID) error
}
