package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// StartDiscussionUseCase starts the discussion phase by setting game data and changing room status
type StartDiscussionUseCase struct {
	roomRepo room.Repository
}

// NewStartDiscussionUseCase creates a new StartDiscussionUseCase
func NewStartDiscussionUseCase(roomRepo room.Repository) *StartDiscussionUseCase {
	return &StartDiscussionUseCase{
		roomRepo: roomRepo,
	}
}

// StartDiscussionInput represents input for starting discussion
type StartDiscussionInput struct {
	RoomID          string
	OriginalEmojis  []string
	DisplayedEmojis []string
	DummyIndex      int
	DummyEmoji      string
}

// Execute starts the discussion phase
func (uc *StartDiscussionUseCase) Execute(ctx context.Context, input StartDiscussionInput) error {
	// Validate input
	if input.RoomID == "" {
		return errors.New("room ID is required")
	}

	// Convert room ID
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return err
	}

	// Find room
	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	// Validate and create game data
	originalEmojis := room.NewEmojiList(input.OriginalEmojis)
	displayedEmojis := room.NewEmojiList(input.DisplayedEmojis)

	dummyIndex, err := room.NewDummyIndex(input.DummyIndex)
	if err != nil {
		return err
	}

	dummyEmoji, err := room.NewDummyEmoji(input.DummyEmoji)
	if err != nil {
		return err
	}

	// Set game data
	foundRoom.SetGameData(
		originalEmojis,
		displayedEmojis,
		dummyIndex,
		dummyEmoji,
	)

	// Save room with game data
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	// Change status to discussing
	foundRoom.ChangeStatus(room.StatusDiscussing)

	// Save room with new status
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
