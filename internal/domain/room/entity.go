package room

import "time"

// Room represents a room entity
type Room struct {
	id              RoomID
	code            RoomCode
	themeID         ThemeID
	topic           *Topic
	answer          *Answer
	status          RoomStatus
	hostUserID      HostUserID
	createdAt       time.Time
	startedAt       *time.Time
	// Game data fields
	originalEmojis  *EmojiList
	displayedEmojis *EmojiList
	dummyIndex      *DummyIndex
	dummyEmoji      *DummyEmoji
	assignments     *Assignments
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

func (r *Room) Topic() *Topic {
	return r.topic
}

func (r *Room) Answer() *Answer {
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
func (r *Room) OriginalEmojis() *EmojiList {
	return r.originalEmojis
}

func (r *Room) DisplayedEmojis() *EmojiList {
	return r.displayedEmojis
}

func (r *Room) DummyIndex() *DummyIndex {
	return r.dummyIndex
}

func (r *Room) DummyEmoji() *DummyEmoji {
	return r.dummyEmoji
}

func (r *Room) Assignments() *Assignments {
	return r.assignments
}

// SetTopic sets the topic for the room
func (r *Room) SetTopic(topic Topic) error {
	if r.status != StatusSettingTopic {
		return ErrInvalidStatusTransition
	}
	r.topic = &topic
	return nil
}

// SetGameData sets the game data (emojis, dummy info)
func (r *Room) SetGameData(originalEmojis, displayedEmojis EmojiList, dummyIndex DummyIndex, dummyEmoji DummyEmoji) error {
	r.originalEmojis = &originalEmojis
	r.displayedEmojis = &displayedEmojis
	r.dummyIndex = &dummyIndex
	r.dummyEmoji = &dummyEmoji
	return nil
}

// SetAnswer sets the answer for the room
func (r *Room) SetAnswer(answer Answer) error {
	r.answer = &answer
	return nil
}

// SetAssignments sets the emoji assignments
func (r *Room) SetAssignments(assignments Assignments) error {
	r.assignments = &assignments
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

// ChangeStatus changes the room status with validation
func (r *Room) ChangeStatus(status RoomStatus) error {
	if !r.status.CanTransitionTo(status) {
		return ErrInvalidStatusTransition
	}
	r.status = status
	return nil
}

// SetStatus sets the room status without validation (for repository reconstruction)
func (r *Room) SetStatus(status RoomStatus) {
	r.status = status
}

// SetTopicUnchecked sets the topic without validation (for repository reconstruction)
func (r *Room) SetTopicUnchecked(topic *Topic) {
	r.topic = topic
}
