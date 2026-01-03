package room

import "time"

// Room represents a room entity
type Room struct {
	id         RoomID
	code       RoomCode
	themeID    ThemeID
	topic      string
	answer     string
	status     RoomStatus
	hostUserID HostUserID
	createdAt  time.Time
	startedAt  *time.Time
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

// SetTopic sets the topic and answer for the room
func (r *Room) SetTopic(topic, answer string) error {
	if r.status != StatusSettingTopic {
		return ErrInvalidStatusTransition
	}
	r.topic = topic
	r.answer = answer
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
