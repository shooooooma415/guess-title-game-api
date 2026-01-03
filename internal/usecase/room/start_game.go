package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// StartGameInput represents the input for starting a game
type StartGameInput struct {
	RoomID string
	UserID string
}

// StartGameUseCase handles the logic for starting a game
type StartGameUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewStartGameUseCase creates a new StartGameUseCase
func NewStartGameUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *StartGameUseCase {
	return &StartGameUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute starts a game
func (uc *StartGameUseCase) Execute(ctx context.Context, input StartGameInput) error {
	// Find room
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return err
	}

	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	// Verify user is host
	participantRoomID, _ := participant.NewRoomIDFromString(input.RoomID)
	participantUserID, _ := participant.NewUserIDFromString(input.UserID)

	foundParticipant, err := uc.participantRepo.FindByRoomAndUser(ctx, participantRoomID, participantUserID)
	if err != nil {
		return errors.New("participant not found")
	}

	if foundParticipant.Role() != participant.RoleHost {
		return errors.New("only host can start the game")
	}

	// Start the game
	if err := foundRoom.Start(); err != nil {
		return err
	}

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
