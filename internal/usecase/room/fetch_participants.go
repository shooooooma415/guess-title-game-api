package room

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
)

// FetchRoomParticipantsUseCase fetches all participants in a room with user information
type FetchRoomParticipantsUseCase struct {
	participantRepo participant.Repository
	userRepo        user.Repository
}

// NewFetchRoomParticipantsUseCase creates a new FetchRoomParticipantsUseCase
func NewFetchRoomParticipantsUseCase(
	participantRepo participant.Repository,
	userRepo user.Repository,
) *FetchRoomParticipantsUseCase {
	return &FetchRoomParticipantsUseCase{
		participantRepo: participantRepo,
		userRepo:        userRepo,
	}
}

// FetchRoomParticipantsInput represents input for fetching room participants
type FetchRoomParticipantsInput struct {
	RoomID string
}

// ParticipantInfo represents participant information with user details
type ParticipantInfo struct {
	UserID   string
	UserName string
	Role     string
	IsLeader bool
}

// FetchRoomParticipantsOutput represents output for fetching room participants
type FetchRoomParticipantsOutput struct {
	Participants []ParticipantInfo
}

// Execute fetches all participants in a room with their user information
func (uc *FetchRoomParticipantsUseCase) Execute(ctx context.Context, input FetchRoomParticipantsInput) (*FetchRoomParticipantsOutput, error) {
	// Validate input
	if input.RoomID == "" {
		return nil, errors.New("room ID is required")
	}

	// Convert room ID
	participantRoomID, err := participant.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return nil, err
	}

	// Fetch participants
	participants, err := uc.participantRepo.FindByRoomID(ctx, participantRoomID)
	if err != nil {
		return nil, err
	}

	// Build output with user information
	participantInfoList := []ParticipantInfo{}
	for _, p := range participants {
		// Fetch user info
		userID, err := user.NewUserIDFromString(p.UserID().String())
		userName := "Unknown"
		if err == nil {
			u, err := uc.userRepo.FindByID(ctx, userID)
			if err == nil {
				userName = u.Name().String()
			}
		}

		participantInfoList = append(participantInfoList, ParticipantInfo{
			UserID:   p.UserID().String(),
			UserName: userName,
			Role:     p.Role().String(),
			IsLeader: p.IsLeader(),
		})
	}

	return &FetchRoomParticipantsOutput{
		Participants: participantInfoList,
	}, nil
}
