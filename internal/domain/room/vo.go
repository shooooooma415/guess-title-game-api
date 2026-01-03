package room

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// RoomID represents a room identifier
type RoomID struct {
	value string
}

func NewRoomID() RoomID {
	return RoomID{value: utils.GenerateUUID()}
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

func (id RoomID) Equals(other RoomID) bool {
	return id.value == other.value
}

// RoomCode represents a room join code
type RoomCode struct {
	value string
}

func NewRoomCode() RoomCode {
	// Generate 6-digit code
	code := rand.Intn(900000) + 100000
	return RoomCode{value: fmt.Sprintf("%06d", code)}
}

func NewRoomCodeFromString(value string) (RoomCode, error) {
	if len(value) != 6 {
		return RoomCode{}, errors.New("room code must be 6 characters")
	}
	return RoomCode{value: value}, nil
}

func (c RoomCode) String() string {
	return c.value
}

// ThemeID represents a theme identifier
type ThemeID struct {
	value string
}

func NewThemeIDFromString(value string) (ThemeID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return ThemeID{}, err
	}
	return ThemeID{value: value}, nil
}

func (id ThemeID) String() string {
	return id.value
}

// HostUserID represents the host user identifier
type HostUserID struct {
	value string
}

func NewHostUserIDFromString(value string) (HostUserID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return HostUserID{}, err
	}
	return HostUserID{value: value}, nil
}

func (id HostUserID) String() string {
	return id.value
}

// RoomStatus represents the current status of a room
type RoomStatus int

const (
	StatusWaiting RoomStatus = iota
	StatusSettingTopic
	StatusDiscussing
	StatusAnswering
	StatusChecking
	StatusFinished
)

func (s RoomStatus) String() string {
	switch s {
	case StatusWaiting:
		return "waiting"
	case StatusSettingTopic:
		return "setting_topic"
	case StatusDiscussing:
		return "discussing"
	case StatusAnswering:
		return "answering"
	case StatusChecking:
		return "checking"
	case StatusFinished:
		return "finished"
	default:
		return "unknown"
	}
}

func NewRoomStatusFromString(value string) (RoomStatus, error) {
	switch value {
	case "waiting":
		return StatusWaiting, nil
	case "setting_topic":
		return StatusSettingTopic, nil
	case "discussing":
		return StatusDiscussing, nil
	case "answering":
		return StatusAnswering, nil
	case "checking":
		return StatusChecking, nil
	case "finished":
		return StatusFinished, nil
	default:
		return 0, errors.New("invalid room status")
	}
}
