package participant

import (
	"errors"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

// ParticipantID represents a participant identifier
type ParticipantID struct {
	value string
}

func NewParticipantID() ParticipantID {
	return ParticipantID{value: utils.GenerateUUID()}
}

func NewParticipantIDFromString(value string) (ParticipantID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return ParticipantID{}, err
	}
	return ParticipantID{value: value}, nil
}

func (id ParticipantID) String() string {
	return id.value
}

// RoomID represents a room identifier (reference to room domain)
type RoomID struct {
	value string
}

func NewRoomIDFromString(value string) (RoomID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return RoomID{}, err
	}
	return RoomID{value: value}, nil
}

func (id RoomID) String() string {
	return id.value
}

// UserID represents a user identifier (reference to user domain)
type UserID struct {
	value string
}

func NewUserIDFromString(value string) (UserID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return UserID{}, err
	}
	return UserID{value: value}, nil
}

func (id UserID) String() string {
	return id.value
}

// ParticipantRole represents the role of a participant
type ParticipantRole int

const (
	RoleHost ParticipantRole = iota
	RolePlayer
)

func (r ParticipantRole) String() string {
	switch r {
	case RoleHost:
		return "host"
	case RolePlayer:
		return "player"
	default:
		return "unknown"
	}
}

func NewParticipantRoleFromString(value string) (ParticipantRole, error) {
	switch value {
	case "host":
		return RoleHost, nil
	case "player":
		return RolePlayer, nil
	default:
		return 0, errors.New("invalid participant role")
	}
}
