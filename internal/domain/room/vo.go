package room

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidStatus           = errors.New("invalid room status")
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
		return 0, ErrInvalidStatus
	}
}

// CanTransitionTo checks if the current status can transition to the target status
func (s RoomStatus) CanTransitionTo(target RoomStatus) bool {
	validTransitions := map[RoomStatus][]RoomStatus{
		StatusWaiting:      {StatusSettingTopic},
		StatusSettingTopic: {StatusDiscussing},
		StatusDiscussing:   {StatusAnswering},
		StatusAnswering:    {StatusChecking},
		StatusChecking:     {StatusFinished},
		StatusFinished:     {},
	}

	allowedTargets, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedTargets {
		if allowed == target {
			return true
		}
	}
	return false
}

// Topic represents a room topic
type Topic struct {
	value string
}

func NewTopic(value string) (Topic, error) {
	if value == "" {
		return Topic{}, errors.New("topic cannot be empty")
	}
	return Topic{value: value}, nil
}

func (t Topic) String() string {
	return t.value
}

func (t Topic) IsEmpty() bool {
	return t.value == ""
}

// Answer represents a room answer
type Answer struct {
	value string
}

func NewAnswer(value string) (Answer, error) {
	if value == "" {
		return Answer{}, errors.New("answer cannot be empty")
	}
	return Answer{value: value}, nil
}

func (a Answer) String() string {
	return a.value
}

func (a Answer) IsEmpty() bool {
	return a.value == ""
}

// EmojiList represents a list of emojis
type EmojiList struct {
	value []string
}

func NewEmojiList(emojis []string) EmojiList {
	return EmojiList{value: emojis}
}

func (e EmojiList) Values() []string {
	return e.value
}

func (e EmojiList) IsEmpty() bool {
	return len(e.value) == 0
}

func (e EmojiList) Count() int {
	return len(e.value)
}

// DummyIndex represents the index of the dummy emoji
type DummyIndex struct {
	value int
}

func NewDummyIndex(value int) (DummyIndex, error) {
	if value < 0 {
		return DummyIndex{}, errors.New("dummy index must be non-negative")
	}
	return DummyIndex{value: value}, nil
}

func (d DummyIndex) Value() int {
	return d.value
}

// DummyEmoji represents the dummy emoji
type DummyEmoji struct {
	value string
}

func NewDummyEmoji(value string) (DummyEmoji, error) {
	if value == "" {
		return DummyEmoji{}, errors.New("dummy emoji cannot be empty")
	}
	return DummyEmoji{value: value}, nil
}

func (d DummyEmoji) String() string {
	return d.value
}

// Assignments represents emoji assignments
type Assignments struct {
	value []string
}

func NewAssignments(assignments []string) Assignments {
	return Assignments{value: assignments}
}

func (a Assignments) Values() []string {
	return a.value
}

func (a Assignments) IsEmpty() bool {
	return len(a.value) == 0
}

func (a Assignments) Count() int {
	return len(a.value)
}
