package room_emoji

import (
	"errors"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

// RoomEmojiID represents a room emoji identifier
type RoomEmojiID struct {
	value string
}

func NewRoomEmojiID() RoomEmojiID {
	return RoomEmojiID{value: utils.GenerateUUID()}
}

func NewRoomEmojiIDFromString(value string) (RoomEmojiID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return RoomEmojiID{}, err
	}
	return RoomEmojiID{value: value}, nil
}

func (id RoomEmojiID) String() string {
	return id.value
}

// RoomID represents a room identifier
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

// ParticipantID represents a participant identifier
type ParticipantID struct {
	value string
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

// Emoji represents an emoji value
type Emoji struct {
	value string
}

func NewEmoji(value string) (Emoji, error) {
	if value == "" {
		return Emoji{}, errors.New("emoji cannot be empty")
	}
	if len(value) > 50 {
		return Emoji{}, errors.New("emoji is too long")
	}
	return Emoji{value: value}, nil
}

func (e Emoji) String() string {
	return e.value
}
