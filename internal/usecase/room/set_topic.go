package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// SetTopicInput represents the input for setting a topic
type SetTopicInput struct {
	RoomID string
	UserID string
	Topic  string
	Emojis []string
}

// SetTopicUseCase handles the logic for setting a topic
type SetTopicUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewSetTopicUseCase creates a new SetTopicUseCase
func NewSetTopicUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *SetTopicUseCase {
	return &SetTopicUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute sets a topic for the room
func (uc *SetTopicUseCase) Execute(ctx context.Context, input SetTopicInput) error {
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
		return errors.New("only host can set topic")
	}

	// Set topic
	topic, err := room.NewTopic(input.Topic)
	if err != nil {
		return err
	}
	if err := foundRoom.SetTopic(topic); err != nil {
		return err
	}

	// Save room
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		return err
	}

	return nil
}
