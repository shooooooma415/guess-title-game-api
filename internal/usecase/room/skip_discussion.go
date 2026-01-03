package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// SkipDiscussionInput represents the input for skipping discussion
type SkipDiscussionInput struct {
	RoomID string
	UserID string
}

// SkipDiscussionUseCase handles the logic for skipping discussion
type SkipDiscussionUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewSkipDiscussionUseCase creates a new SkipDiscussionUseCase
func NewSkipDiscussionUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *SkipDiscussionUseCase {
	return &SkipDiscussionUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute skips discussion and moves to answering phase
func (uc *SkipDiscussionUseCase) Execute(ctx context.Context, input SkipDiscussionInput) error {
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
		return errors.New("only host can skip discussion")
	}

	// Verify dummy data exists
	if foundRoom.DummyIndex() == nil {
		return errors.New("dummy data is required before skipping discussion")
	}

	// Change status to answering
	if err := foundRoom.ChangeStatus(room.StatusAnswering); err != nil {
		return err
	}

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
