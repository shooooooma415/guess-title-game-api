package participant

import "time"

// Participant represents a participant in a room
type Participant struct {
	id       ParticipantID
	roomID   RoomID
	userID   UserID
	role     ParticipantRole
	isLeader bool
	joinedAt time.Time
}

// NewParticipant creates a new Participant
func NewParticipant(
	id ParticipantID,
	roomID RoomID,
	userID UserID,
	role ParticipantRole,
) *Participant {
	return &Participant{
		id:       id,
		roomID:   roomID,
		userID:   userID,
		role:     role,
		isLeader: false,
		joinedAt: time.Now(),
	}
}

// Getters
func (p *Participant) ID() ParticipantID {
	return p.id
}

func (p *Participant) RoomID() RoomID {
	return p.roomID
}

func (p *Participant) UserID() UserID {
	return p.userID
}

func (p *Participant) Role() ParticipantRole {
	return p.role
}

func (p *Participant) IsLeader() bool {
	return p.isLeader
}

func (p *Participant) JoinedAt() time.Time {
	return p.joinedAt
}

// SetAsLeader sets this participant as the leader
func (p *Participant) SetAsLeader() {
	p.isLeader = true
}

// RemoveLeader removes the leader status
func (p *Participant) RemoveLeader() {
	p.isLeader = false
}
