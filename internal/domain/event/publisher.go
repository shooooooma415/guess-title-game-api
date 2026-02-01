package event

// Publisher defines the interface for publishing domain events
type Publisher interface {
	Publish(event Event)
	Subscribe(eventType string, handler EventHandler)
}

// EventHandler is a function that handles an event
type EventHandler func(event Event)
