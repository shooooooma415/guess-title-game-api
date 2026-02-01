package user

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
)

// JoinRoomInput represents the input for joining a room
type JoinRoomInput struct {
	RoomCode string
	UserName string
}

// JoinRoomOutput represents the output after joining a room
type JoinRoomOutput struct {
	RoomID   string
	UserID   string
	IsLeader bool
}

// JoinRoomUseCase handles the logic for a user joining a room
type JoinRoomUseCase struct {
	userRepo        user.Repository
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewJoinRoomUseCase creates a new JoinRoomUseCase
func NewJoinRoomUseCase(
	userRepo user.Repository,
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *JoinRoomUseCase {
	return &JoinRoomUseCase{
		userRepo:        userRepo,
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute joins a user to a room
func (uc *JoinRoomUseCase) Execute(ctx context.Context, input JoinRoomInput) (*JoinRoomOutput, error) {
	// Validate input
	if input.RoomCode == "" {
		return nil, errors.New("room code is required")
	}
	if input.UserName == "" {
		return nil, errors.New("user name is required")
	}

	// Find room by code
	roomCode, err := room.NewRoomCodeFromString(input.RoomCode)
	if err != nil {
		return nil, err
	}

	foundRoom, err := uc.roomRepo.FindByCode(ctx, roomCode)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Create new user
	userID := user.NewUserID()
	userName, err := user.NewUserName(input.UserName)
	if err != nil {
		return nil, err
	}

	newUser := user.NewUser(userID, userName)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	// Check if this is the first non-host participant (first joiner becomes Leader)
	participantRoomID, _ := participant.NewRoomIDFromString(foundRoom.ID().String())
	existingParticipants, err := uc.participantRepo.FindByRoomID(ctx, participantRoomID)
	if err != nil {
		existingParticipants = []*participant.Participant{}
	}
	// Count only non-host participants to determine if this is the first joiner
	nonHostCount := 0
	for _, p := range existingParticipants {
		if p.Role() != participant.RoleHost {
			nonHostCount++
		}
	}
	isLeader := nonHostCount == 0

	// Create participant
	participantID := participant.NewParticipantID()
	participantUserID, _ := participant.NewUserIDFromString(userID.String())
	role := participant.RolePlayer

	newParticipant := participant.NewParticipant(
		participantID,
		participantRoomID,
		participantUserID,
		role,
	)

	if isLeader {
		newParticipant.SetAsLeader()
	}

	if err := uc.participantRepo.Save(ctx, newParticipant); err != nil {
		return nil, err
	}

	return &JoinRoomOutput{
		RoomID:   foundRoom.ID().String(),
		UserID:   userID.String(),
		IsLeader: isLeader,
	}, nil
}
