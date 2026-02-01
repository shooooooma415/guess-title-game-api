package websocket

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client -> Server
	MessageTypeClientConnected  MessageType = "CLIENT_CONNECTED"
	MessageTypeFetchParticipants MessageType = "FETCH_PARTICIPANTS"
	MessageTypeSubmitTopic       MessageType = "SUBMIT_TOPIC"
	MessageTypeAnswering         MessageType = "ANSWERING"

	// Server -> Client
	MessageTypeStateUpdate        MessageType = "STATE_UPDATE"
	MessageTypeParticipantUpdate  MessageType = "PARTICIPANT_UPDATE"
	MessageTypeTimerTick          MessageType = "TIMER_TICK"
	MessageTypeError              MessageType = "ERROR"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ClientConnectedPayload represents the payload for CLIENT_CONNECTED
type ClientConnectedPayload struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

// SubmitTopicPayload represents the payload for SUBMIT_TOPIC
type SubmitTopicPayload struct {
	DisplayedEmojis []string `json:"displayedEmojis"`
	OriginalEmojis  []string `json:"originalEmojis"`
	DummyIndex      int      `json:"dummyIndex"`
	DummyEmoji      string   `json:"dummyEmoji"`
}

// AnsweringPayload represents the payload for ANSWERING
type AnsweringPayload struct {
	Answer          string   `json:"answer"`
	DisplayedEmojis []string `json:"displayedEmojis"`
	OriginalEmojis  []string `json:"originalEmojis"`
	DummyIndex      int      `json:"dummyIndex"`
	DummyEmoji      string   `json:"dummyEmoji"`
}

// StateUpdatePayload represents the payload for STATE_UPDATE
type StateUpdatePayload struct {
	NextState string                `json:"nextState"`
	Data      *StateUpdateDataPayload `json:"data,omitempty"`
}

// StateUpdateDataPayload represents the data in STATE_UPDATE
type StateUpdateDataPayload struct {
	Theme           string   `json:"theme,omitempty"`
	Topic           string   `json:"topic,omitempty"`
	Answer          string   `json:"answer,omitempty"`
	DisplayedEmojis []string `json:"displayedEmojis,omitempty"`
	OriginalEmojis  []string `json:"originalEmojis,omitempty"`
	DummyIndex      *int     `json:"dummyIndex,omitempty"`
	DummyEmoji      string   `json:"dummyEmoji,omitempty"`
	Assignments     []string `json:"assignments,omitempty"`
}

// ParticipantData represents participant information
type ParticipantData struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	IsLeader bool   `json:"is_leader"`
}

// ParticipantUpdatePayload represents the payload for PARTICIPANT_UPDATE
type ParticipantUpdatePayload struct {
	Participants []ParticipantData `json:"participants"`
}

// TimerTickPayload represents the payload for TIMER_TICK
type TimerTickPayload struct {
	Time string `json:"time"`
}

// ErrorPayload represents the payload for ERROR
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
