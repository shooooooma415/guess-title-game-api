package event

import "time"

// Event represents a domain event
type Event interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
}

// BaseEvent provides common event fields
type BaseEvent struct {
	eventType   string
	occurredAt  time.Time
	aggregateID string
}

func (e BaseEvent) EventType() string {
	return e.eventType
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

// GameStartedEvent is fired when a game starts (WAITING -> SETTING_TOPIC)
type GameStartedEvent struct {
	BaseEvent
	RoomID string
	Status string // "setting_topic"
}

func NewGameStartedEvent(roomID string) *GameStartedEvent {
	return &GameStartedEvent{
		BaseEvent: BaseEvent{
			eventType:   "GameStarted",
			occurredAt:  time.Now(),
			aggregateID: roomID,
		},
		RoomID: roomID,
		Status: "setting_topic",
	}
}

// DiscussionSkippedEvent is fired when discussion is skipped (DISCUSSING -> ANSWERING)
type DiscussionSkippedEvent struct {
	BaseEvent
	RoomID string
	Status string // "answering"
}

func NewDiscussionSkippedEvent(roomID string) *DiscussionSkippedEvent {
	return &DiscussionSkippedEvent{
		BaseEvent: BaseEvent{
			eventType:   "DiscussionSkipped",
			occurredAt:  time.Now(),
			aggregateID: roomID,
		},
		RoomID: roomID,
		Status: "answering",
	}
}

// AnswerSubmittedEvent is fired when answer is submitted (ANSWERING -> CHECKING)
type AnswerSubmittedEvent struct {
	BaseEvent
	RoomID string
	Status string // "checking"
}

func NewAnswerSubmittedEvent(roomID string) *AnswerSubmittedEvent {
	return &AnswerSubmittedEvent{
		BaseEvent: BaseEvent{
			eventType:   "AnswerSubmitted",
			occurredAt:  time.Now(),
			aggregateID: roomID,
		},
		RoomID: roomID,
		Status: "checking",
	}
}

// GameFinishedEvent is fired when a game finishes (CHECKING -> FINISHED)
type GameFinishedEvent struct {
	BaseEvent
	RoomID string
	Status string // "finished"
}

func NewGameFinishedEvent(roomID string) *GameFinishedEvent {
	return &GameFinishedEvent{
		BaseEvent: BaseEvent{
			eventType:   "GameFinished",
			occurredAt:  time.Now(),
			aggregateID: roomID,
		},
		RoomID: roomID,
		Status: "finished",
	}
}
