package websocket

import (
	"context"
	"log"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/event"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/theme"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

// SetupEventHandlers sets up event handlers for WebSocket broadcasting
func (h *Handler) SetupEventHandlers(eventPublisher event.Publisher) {
	// Subscribe to GameStartedEvent
	eventPublisher.Subscribe("GameStarted", func(evt event.Event) {
		gameStartedEvt, ok := evt.(*event.GameStartedEvent)
		if !ok {
			log.Printf("Invalid event type for GameStarted")
			return
		}

		h.handleGameStartedEvent(gameStartedEvt)
	})

	// Subscribe to DiscussionSkippedEvent
	eventPublisher.Subscribe("DiscussionSkipped", func(evt event.Event) {
		discussionSkippedEvt, ok := evt.(*event.DiscussionSkippedEvent)
		if !ok {
			log.Printf("Invalid event type for DiscussionSkipped")
			return
		}

		h.handleDiscussionSkippedEvent(discussionSkippedEvt)
	})

	// Subscribe to AnswerSubmittedEvent
	eventPublisher.Subscribe("AnswerSubmitted", func(evt event.Event) {
		answerSubmittedEvt, ok := evt.(*event.AnswerSubmittedEvent)
		if !ok {
			log.Printf("Invalid event type for AnswerSubmitted")
			return
		}

		h.handleAnswerSubmittedEvent(answerSubmittedEvt)
	})

	// Subscribe to GameFinishedEvent
	eventPublisher.Subscribe("GameFinished", func(evt event.Event) {
		gameFinishedEvt, ok := evt.(*event.GameFinishedEvent)
		if !ok {
			log.Printf("Invalid event type for GameFinished")
			return
		}

		h.handleGameFinishedEvent(gameFinishedEvt)
	})
}

// handleGameStartedEvent handles GameStartedEvent and broadcasts STATE_UPDATE
func (h *Handler) handleGameStartedEvent(evt *event.GameStartedEvent) {
	ctx := context.Background()

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: evt.RoomID,
	})
	if err != nil {
		log.Printf("Error fetching room for GameStartedEvent: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Broadcast STATE_UPDATE with setting_topic status
	h.hub.Broadcast(evt.RoomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: foundRoom.Status().String(), // "setting_topic"
			Data:      &StateUpdateDataPayload{},
		},
	})
}

// handleDiscussionSkippedEvent handles DiscussionSkippedEvent and broadcasts STATE_UPDATE
func (h *Handler) handleDiscussionSkippedEvent(evt *event.DiscussionSkippedEvent) {
	ctx := context.Background()

	// Stop timer
	h.timer.StopTimer(evt.RoomID)

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: evt.RoomID,
	})
	if err != nil {
		log.Printf("Error fetching room for DiscussionSkippedEvent: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Build state data payload
	topicStr := ""
	if foundRoom.Topic() != nil {
		topicStr = foundRoom.Topic().String()
	}

	var dummyIdxPtr *int
	if foundRoom.DummyIndex() != nil {
		val := foundRoom.DummyIndex().Value()
		dummyIdxPtr = &val
	}

	dummyEmojiStr := ""
	if foundRoom.DummyEmoji() != nil {
		dummyEmojiStr = foundRoom.DummyEmoji().String()
	}

	displayedEmojisSlice := []string{}
	if foundRoom.DisplayedEmojis() != nil {
		displayedEmojisSlice = foundRoom.DisplayedEmojis().Values()
	}

	originalEmojisSlice := []string{}
	if foundRoom.OriginalEmojis() != nil {
		originalEmojisSlice = foundRoom.OriginalEmojis().Values()
	}

	assignmentsSlice := []string{}
	if foundRoom.Assignments() != nil {
		assignmentsSlice = foundRoom.Assignments().Values()
	}

	// Broadcast STATE_UPDATE with answering status
	h.hub.Broadcast(evt.RoomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: foundRoom.Status().String(), // "answering"
			Data: &StateUpdateDataPayload{
				Topic:           topicStr,
				DisplayedEmojis: displayedEmojisSlice,
				OriginalEmojis:  originalEmojisSlice,
				DummyIndex:      dummyIdxPtr,
				DummyEmoji:      dummyEmojiStr,
				Assignments:     assignmentsSlice,
			},
		},
	})
}

// handleAnswerSubmittedEvent handles AnswerSubmittedEvent and broadcasts STATE_UPDATE
func (h *Handler) handleAnswerSubmittedEvent(evt *event.AnswerSubmittedEvent) {
	ctx := context.Background()

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: evt.RoomID,
	})
	if err != nil {
		log.Printf("Error fetching room for AnswerSubmittedEvent: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Fetch theme
	themeStr := ""
	themeID, err := theme.NewThemeIDFromString(foundRoom.ThemeID().String())
	if err == nil {
		themeObj, err := h.themeRepo.FindByID(ctx, themeID)
		if err == nil && themeObj != nil {
			themeStr = themeObj.Title().String()
		}
	}

	// Build state data payload
	topicStr := ""
	if foundRoom.Topic() != nil {
		topicStr = foundRoom.Topic().String()
	}

	var dummyIdxPtr *int
	if foundRoom.DummyIndex() != nil {
		val := foundRoom.DummyIndex().Value()
		dummyIdxPtr = &val
	}

	dummyEmojiStr := ""
	if foundRoom.DummyEmoji() != nil {
		dummyEmojiStr = foundRoom.DummyEmoji().String()
	}

	displayedEmojisSlice := []string{}
	if foundRoom.DisplayedEmojis() != nil {
		displayedEmojisSlice = foundRoom.DisplayedEmojis().Values()
	}

	originalEmojisSlice := []string{}
	if foundRoom.OriginalEmojis() != nil {
		originalEmojisSlice = foundRoom.OriginalEmojis().Values()
	}

	assignmentsSlice := []string{}
	if foundRoom.Assignments() != nil {
		assignmentsSlice = foundRoom.Assignments().Values()
	}

	answerStr := ""
	if foundRoom.Answer() != nil {
		answerStr = foundRoom.Answer().String()
	}

	// Broadcast STATE_UPDATE with checking status (include answer and theme)
	h.hub.Broadcast(evt.RoomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: foundRoom.Status().String(), // "checking"
			Data: &StateUpdateDataPayload{
				Theme:           themeStr,
				Topic:           topicStr,
				Answer:          answerStr,
				DisplayedEmojis: displayedEmojisSlice,
				OriginalEmojis:  originalEmojisSlice,
				DummyIndex:      dummyIdxPtr,
				DummyEmoji:      dummyEmojiStr,
				Assignments:     assignmentsSlice,
			},
		},
	})

	log.Printf("Answer submitted event broadcasted for room %s with answer: %s, theme: %s", evt.RoomID, answerStr, themeStr)
}

// handleGameFinishedEvent handles GameFinishedEvent and broadcasts STATE_UPDATE
func (h *Handler) handleGameFinishedEvent(evt *event.GameFinishedEvent) {
	ctx := context.Background()

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: evt.RoomID,
	})
	if err != nil {
		log.Printf("Error fetching room for GameFinishedEvent: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Fetch theme
	themeStr := ""
	themeID, err := theme.NewThemeIDFromString(foundRoom.ThemeID().String())
	if err == nil {
		themeObj, err := h.themeRepo.FindByID(ctx, themeID)
		if err == nil && themeObj != nil {
			themeStr = themeObj.Title().String()
		}
	}

	// Build state data payload
	topicStr := ""
	if foundRoom.Topic() != nil {
		topicStr = foundRoom.Topic().String()
	}

	var dummyIdxPtr *int
	if foundRoom.DummyIndex() != nil {
		val := foundRoom.DummyIndex().Value()
		dummyIdxPtr = &val
	}

	dummyEmojiStr := ""
	if foundRoom.DummyEmoji() != nil {
		dummyEmojiStr = foundRoom.DummyEmoji().String()
	}

	displayedEmojisSlice := []string{}
	if foundRoom.DisplayedEmojis() != nil {
		displayedEmojisSlice = foundRoom.DisplayedEmojis().Values()
	}

	originalEmojisSlice := []string{}
	if foundRoom.OriginalEmojis() != nil {
		originalEmojisSlice = foundRoom.OriginalEmojis().Values()
	}

	assignmentsSlice := []string{}
	if foundRoom.Assignments() != nil {
		assignmentsSlice = foundRoom.Assignments().Values()
	}

	answerStr := ""
	if foundRoom.Answer() != nil {
		answerStr = foundRoom.Answer().String()
	}

	// Broadcast STATE_UPDATE with finished status (include answer and theme)
	h.hub.Broadcast(evt.RoomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: foundRoom.Status().String(), // "finished"
			Data: &StateUpdateDataPayload{
				Theme:           themeStr,
				Topic:           topicStr,
				Answer:          answerStr,
				DisplayedEmojis: displayedEmojisSlice,
				OriginalEmojis:  originalEmojisSlice,
				DummyIndex:      dummyIdxPtr,
				DummyEmoji:      dummyEmojiStr,
				Assignments:     assignmentsSlice,
			},
		},
	})

	log.Printf("Game finished event broadcasted for room %s with answer: %s, theme: %s", evt.RoomID, answerStr, themeStr)
}
