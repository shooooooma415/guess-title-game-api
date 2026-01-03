package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// SubmitAnswerInput represents the input for submitting an answer
type SubmitAnswerInput struct {
	RoomID string
	UserID string
	Answer string
}

// SubmitAnswerUseCase handles the logic for submitting an answer
type SubmitAnswerUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewSubmitAnswerUseCase creates a new SubmitAnswerUseCase
func NewSubmitAnswerUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *SubmitAnswerUseCase {
	return &SubmitAnswerUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute submits an answer
func (uc *SubmitAnswerUseCase) Execute(ctx context.Context, input SubmitAnswerInput) error {
	// Find room
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return err
	}

	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	// Verify user is leader
	participantRoomID, _ := participant.NewRoomIDFromString(input.RoomID)
	participantUserID, _ := participant.NewUserIDFromString(input.UserID)

	foundParticipant, err := uc.participantRepo.FindByRoomAndUser(ctx, participantRoomID, participantUserID)
	if err != nil {
		return errors.New("participant not found")
	}

	if !foundParticipant.IsLeader() {
		return errors.New("only leader can submit answer")
	}

	// Set answer
	answer, err := room.NewAnswer(input.Answer)
	if err != nil {
		return err
	}
	if err := foundRoom.SetAnswer(answer); err != nil {
		return err
	}

	// Change status to checking
	if err := foundRoom.ChangeStatus(room.StatusChecking); err != nil {
		return err
	}

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
