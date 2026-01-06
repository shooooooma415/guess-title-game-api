package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// FetchRoomUseCase fetches room information
type FetchRoomUseCase struct {
	roomRepo room.Repository
}

// NewFetchRoomUseCase creates a new FetchRoomUseCase
func NewFetchRoomUseCase(roomRepo room.Repository) *FetchRoomUseCase {
	return &FetchRoomUseCase{
		roomRepo: roomRepo,
	}
}

// FetchRoomInput represents input for fetching room
type FetchRoomInput struct {
	RoomID string
}

// FetchRoomOutput represents output for fetching room
type FetchRoomOutput struct {
	Room *room.Room
}

// Execute fetches room information
func (uc *FetchRoomUseCase) Execute(ctx context.Context, input FetchRoomInput) (*FetchRoomOutput, error) {
	// Validate input
	if input.RoomID == "" {
		return nil, errors.New("room ID is required")
	}

	// Convert room ID
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return nil, err
	}

	// Fetch room
	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	return &FetchRoomOutput{
		Room: foundRoom,
	}, nil
}
