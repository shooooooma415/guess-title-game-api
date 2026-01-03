package websocket

import (
	"fmt"
	"sync"
	"time"
)

const (
	DiscussionDuration = 300 * time.Second // 5 minutes
	StartDelay         = 5 * time.Second   // 5 seconds
)

// Timer manages game timers for rooms
type Timer struct {
	hub        *Hub
	timers     map[string]*RoomTimer
	timerMutex sync.RWMutex
}

// RoomTimer represents a timer for a specific room
type RoomTimer struct {
	roomID    string
	remaining time.Duration
	ticker    *time.Ticker
	stopChan  chan bool
	stopped   bool
	mu        sync.Mutex
}

// NewTimer creates a new Timer
func NewTimer(hub *Hub) *Timer {
	return &Timer{
		hub:    hub,
		timers: make(map[string]*RoomTimer),
	}
}

// StartTimer starts a timer for a room with initial delay
func (t *Timer) StartTimer(roomID string) {
	t.timerMutex.Lock()
	defer t.timerMutex.Unlock()

	// Stop existing timer if any
	if existingTimer, exists := t.timers[roomID]; exists {
		existingTimer.Stop()
	}

	// Wait for start delay, then start the actual timer
	go func() {
		time.Sleep(StartDelay)

		roomTimer := &RoomTimer{
			roomID:    roomID,
			remaining: DiscussionDuration,
			ticker:    time.NewTicker(1 * time.Second),
			stopChan:  make(chan bool),
			stopped:   false,
		}

		t.timerMutex.Lock()
		t.timers[roomID] = roomTimer
		t.timerMutex.Unlock()

		roomTimer.Run(t.hub)
	}()
}

// StopTimer stops the timer for a room
func (t *Timer) StopTimer(roomID string) {
	t.timerMutex.Lock()
	defer t.timerMutex.Unlock()

	if timer, exists := t.timers[roomID]; exists {
		timer.Stop()
		delete(t.timers, roomID)
	}
}

// Run starts the room timer
func (rt *RoomTimer) Run(hub *Hub) {
	defer func() {
		rt.ticker.Stop()
	}()

	for {
		select {
		case <-rt.ticker.C:
			rt.mu.Lock()
			if rt.stopped {
				rt.mu.Unlock()
				return
			}

			rt.remaining -= 1 * time.Second

			if rt.remaining <= 0 {
				rt.mu.Unlock()
				// Timer finished - send final tick
				hub.Broadcast(rt.roomID, Message{
					Type: MessageTypeTimerTick,
					Payload: TimerTickPayload{
						Time: "00:00",
					},
				})
				return
			}

			// Send timer tick
			timeStr := formatTime(rt.remaining)
			rt.mu.Unlock()

			hub.Broadcast(rt.roomID, Message{
				Type: MessageTypeTimerTick,
				Payload: TimerTickPayload{
					Time: timeStr,
				},
			})

		case <-rt.stopChan:
			return
		}
	}
}

// Stop stops the room timer
func (rt *RoomTimer) Stop() {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if !rt.stopped {
		rt.stopped = true
		close(rt.stopChan)
	}
}

// formatTime formats duration as "MM:SS"
func formatTime(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
