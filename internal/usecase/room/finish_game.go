package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// FinishGameInput represents the input for finishing a game
type FinishGameInput struct {
	RoomID string
	UserID string
}

// FinishGameUseCase handles the logic for finishing a game
type FinishGameUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewFinishGameUseCase creates a new FinishGameUseCase
func NewFinishGameUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *FinishGameUseCase {
	return &FinishGameUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute finishes a game
func (uc *FinishGameUseCase) Execute(ctx context.Context, input FinishGameInput) error {
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
		return errors.New("only host can finish the game")
	}

	// Change status to finished
	if err := foundRoom.ChangeStatus(room.StatusFinished); err != nil {
		return err
	}

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
