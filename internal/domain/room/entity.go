package room

import "time"

// Room represents a room entity
type Room struct {
	id              RoomID
	code            RoomCode
	themeID         ThemeID
	topic           string
	answer          string
	status          RoomStatus
	hostUserID      HostUserID
	createdAt       time.Time
	startedAt       *time.Time
	// Game data fields
	originalEmojis  []string
	displayedEmojis []string
	dummyIndex      *int
	dummyEmoji      string
	assignments     []string
}

// NewRoom creates a new Room
func NewRoom(
	id RoomID,
	code RoomCode,
	themeID ThemeID,
	hostUserID HostUserID,
) *Room {
	return &Room{
		id:         id,
		code:       code,
		themeID:    themeID,
		hostUserID: hostUserID,
		status:     StatusWaiting,
		createdAt:  time.Now(),
	}
}

// Getters
func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) Code() RoomCode {
	return r.code
}

func (r *Room) ThemeID() ThemeID {
	return r.themeID
}

func (r *Room) Topic() string {
	return r.topic
}

func (r *Room) Answer() string {
	return r.answer
}

func (r *Room) Status() RoomStatus {
	return r.status
}

func (r *Room) HostUserID() HostUserID {
	return r.hostUserID
}

func (r *Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) StartedAt() *time.Time {
	return r.startedAt
}

// Game data getters
func (r *Room) OriginalEmojis() []string {
	return r.originalEmojis
}

func (r *Room) DisplayedEmojis() []string {
	return r.displayedEmojis
}

func (r *Room) DummyIndex() *int {
	return r.dummyIndex
}

func (r *Room) DummyEmoji() string {
	return r.dummyEmoji
}

func (r *Room) Assignments() []string {
	return r.assignments
}

// SetTopic sets the topic for the room
func (r *Room) SetTopic(topic string) error {
	if r.status != StatusSettingTopic {
		return ErrInvalidStatusTransition
	}
	r.topic = topic
	return nil
}

// SetGameData sets the game data (emojis, dummy info)
func (r *Room) SetGameData(originalEmojis, displayedEmojis []string, dummyIndex int, dummyEmoji string) error {
	r.originalEmojis = originalEmojis
	r.displayedEmojis = displayedEmojis
	r.dummyIndex = &dummyIndex
	r.dummyEmoji = dummyEmoji
	return nil
}

// SetAnswer sets the answer for the room
func (r *Room) SetAnswer(answer string) error {
	r.answer = answer
	return nil
}

// SetAssignments sets the emoji assignments
func (r *Room) SetAssignments(assignments []string) error {
	r.assignments = assignments
	return nil
}

// Start starts the room discussion
func (r *Room) Start() error {
	if r.status != StatusWaiting {
		return ErrInvalidStatusTransition
	}
	now := time.Now()
	r.startedAt = &now
	r.status = StatusSettingTopic
	return nil
}

// ChangeStatus changes the room status
func (r *Room) ChangeStatus(status RoomStatus) error {
	// TODO: Add status transition validation logic
	r.status = status
	return nil
}
