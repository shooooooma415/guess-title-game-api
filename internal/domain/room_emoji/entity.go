package room_emoji

// RoomEmoji represents an emoji reaction in a room
type RoomEmoji struct {
	id            RoomEmojiID
	roomID        RoomID
	participantID *ParticipantID
	emoji         Emoji
}

// NewRoomEmoji creates a new RoomEmoji
func NewRoomEmoji(
	id RoomEmojiID,
	roomID RoomID,
	participantID *ParticipantID,
	emoji Emoji,
) *RoomEmoji {
	return &RoomEmoji{
		id:            id,
		roomID:        roomID,
		participantID: participantID,
		emoji:         emoji,
	}
}

// Getters
func (re *RoomEmoji) ID() RoomEmojiID {
	return re.id
}

func (re *RoomEmoji) RoomID() RoomID {
	return re.roomID
}

func (re *RoomEmoji) ParticipantID() *ParticipantID {
	return re.participantID
}

func (re *RoomEmoji) Emoji() Emoji {
	return re.emoji
}
