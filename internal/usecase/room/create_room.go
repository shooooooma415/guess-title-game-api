package room

import (
	"context"
	"math/rand"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/theme"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
)

// CreateRoomOutput represents the output after creating a room
type CreateRoomOutput struct {
	RoomID   string
	UserID   string
	RoomCode string
	Theme    string
	Hint     string
}

// CreateRoomUseCase handles the logic for creating a room
type CreateRoomUseCase struct {
	userRepo        user.Repository
	roomRepo        room.Repository
	themeRepo       theme.Repository
	participantRepo participant.Repository
}

// NewCreateRoomUseCase creates a new CreateRoomUseCase
func NewCreateRoomUseCase(
	userRepo user.Repository,
	roomRepo room.Repository,
	themeRepo theme.Repository,
	participantRepo participant.Repository,
) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		userRepo:        userRepo,
		roomRepo:        roomRepo,
		themeRepo:       themeRepo,
		participantRepo: participantRepo,
	}
}

// Execute creates a new room
func (uc *CreateRoomUseCase) Execute(ctx context.Context) (*CreateRoomOutput, error) {
	// Get a random theme
	themes, err := uc.themeRepo.FindAll(ctx)
	if err != nil || len(themes) == 0 {
		return nil, err
	}
	selectedTheme := themes[rand.Intn(len(themes))]

	// Create host user
	hostUserID := user.NewUserID()
	hostUserName, _ := user.NewUserName("Host")
	hostUser := user.NewUser(hostUserID, hostUserName)
	if err := uc.userRepo.Save(ctx, hostUser); err != nil {
		return nil, err
	}

	// Create room
	roomID := room.NewRoomID()
	roomCode := room.NewRoomCode()
	themeID, _ := room.NewThemeIDFromString(selectedTheme.ID().String())
	hostID, _ := room.NewHostUserIDFromString(hostUserID.String())

	newRoom := room.NewRoom(roomID, roomCode, themeID, hostID)
	if err := uc.roomRepo.Save(ctx, newRoom); err != nil {
		return nil, err
	}

	// Create host participant
	participantID := participant.NewParticipantID()
	participantRoomID, _ := participant.NewRoomIDFromString(roomID.String())
	participantUserID, _ := participant.NewUserIDFromString(hostUserID.String())

	hostParticipant := participant.NewParticipant(
		participantID,
		participantRoomID,
		participantUserID,
		participant.RoleHost,
	)
	hostParticipant.SetAsLeader()

	if err := uc.participantRepo.Save(ctx, hostParticipant); err != nil {
		return nil, err
	}

	return &CreateRoomOutput{
		RoomID:   roomID.String(),
		UserID:   hostUserID.String(),
		RoomCode: roomCode.String(),
		Theme:    selectedTheme.Title().String(),
		Hint:     selectedTheme.Hint().String(),
	}, nil
}
