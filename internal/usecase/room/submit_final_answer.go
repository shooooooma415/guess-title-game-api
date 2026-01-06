package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// SubmitFinalAnswerUseCase submits the final answer with game data and transitions to checking
type SubmitFinalAnswerUseCase struct {
	roomRepo room.Repository
}

// NewSubmitFinalAnswerUseCase creates a new SubmitFinalAnswerUseCase
func NewSubmitFinalAnswerUseCase(roomRepo room.Repository) *SubmitFinalAnswerUseCase {
	return &SubmitFinalAnswerUseCase{
		roomRepo: roomRepo,
	}
}

// SubmitFinalAnswerInput represents input for submitting final answer
type SubmitFinalAnswerInput struct {
	RoomID          string
	Answer          string
	OriginalEmojis  []string
	DisplayedEmojis []string
	DummyIndex      int
	DummyEmoji      string
}

// Execute submits the final answer and transitions to checking phase
func (uc *SubmitFinalAnswerUseCase) Execute(ctx context.Context, input SubmitFinalAnswerInput) error {
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

	// Validate and set answer
	answer, err := room.NewAnswer(input.Answer)
	if err != nil {
		return err
	}
	foundRoom.SetAnswer(answer)

	// Validate and set game data
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

	foundRoom.SetGameData(
		originalEmojis,
		displayedEmojis,
		dummyIndex,
		dummyEmoji,
	)

	// Change status to checking
	foundRoom.ChangeStatus(room.StatusChecking)

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
