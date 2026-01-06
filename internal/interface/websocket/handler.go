package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	roomID string
	userID string
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[string]map[*Client]bool // roomID -> clients
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	RoomID  string
	Message []byte
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.roomID] == nil {
				h.clients[client.roomID] = make(map[*Client]bool)
			}
			h.clients[client.roomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.roomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.roomID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[message.RoomID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.send <- message.Message:
				default:
					close(client.send)
					h.mu.Lock()
					delete(h.clients[message.RoomID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(roomID string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- BroadcastMessage{
		RoomID:  roomID,
		Message: data,
	}
}

// Handler handles WebSocket connections
type Handler struct {
	hub                       *Hub
	timer                     *Timer
	fetchRoomUseCase          *roomUseCase.FetchRoomUseCase
	fetchParticipantsUseCase  *roomUseCase.FetchRoomParticipantsUseCase
	startDiscussionUseCase    *roomUseCase.StartDiscussionUseCase
	submitFinalAnswerUseCase  *roomUseCase.SubmitFinalAnswerUseCase
}

// NewHandler creates a new WebSocket handler
func NewHandler(
	hub *Hub,
	timer *Timer,
	fetchRoomUseCase *roomUseCase.FetchRoomUseCase,
	fetchParticipantsUseCase *roomUseCase.FetchRoomParticipantsUseCase,
	startDiscussionUseCase *roomUseCase.StartDiscussionUseCase,
	submitFinalAnswerUseCase *roomUseCase.SubmitFinalAnswerUseCase,
) *Handler {
	return &Handler{
		hub:                       hub,
		timer:                     timer,
		fetchRoomUseCase:          fetchRoomUseCase,
		fetchParticipantsUseCase:  fetchParticipantsUseCase,
		startDiscussionUseCase:    startDiscussionUseCase,
		submitFinalAnswerUseCase:  submitFinalAnswerUseCase,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	if roomID == "" {
		return echo.NewHTTPError(400, "room_id is required")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return err
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		roomID: roomID,
	}

	h.hub.register <- client

	go h.writePump(client)
	go h.readPump(client)

	return nil
}

// readPump reads messages from the WebSocket connection
func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		h.handleMessage(client, message)
	}
}

// writePump writes messages to the WebSocket connection
func (h *Handler) writePump(client *Client) {
	defer func() {
		client.conn.Close()
	}()

	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			return
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (h *Handler) handleMessage(client *Client, data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		h.sendError(client, "INVALID_MESSAGE", "Invalid message format")
		return
	}

	switch msg.Type {
	case MessageTypeClientConnected:
		h.handleClientConnected(client, msg.Payload)

	case MessageTypeFetchParticipants:
		h.handleFetchParticipants(client)

	case MessageTypeSubmitTopic:
		h.handleSubmitTopic(client, msg.Payload)

	case MessageTypeAnswering:
		h.handleAnswering(client, msg.Payload)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
		h.sendError(client, "UNKNOWN_MESSAGE_TYPE", "Unknown message type: "+string(msg.Type))
	}
}

// handleClientConnected handles CLIENT_CONNECTED message
func (h *Handler) handleClientConnected(client *Client, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	var data ClientConnectedPayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		log.Printf("Error unmarshaling CLIENT_CONNECTED payload: %v", err)
		h.sendError(client, "INVALID_PAYLOAD", "Invalid CLIENT_CONNECTED payload")
		return
	}

	client.userID = data.UserID

	// Send initial room state to the connected client
	h.sendInitialRoomState(client)

	// Broadcast participant update
	h.broadcastParticipantUpdate(client.roomID)
}

// handleFetchParticipants handles FETCH_PARTICIPANTS message
func (h *Handler) handleFetchParticipants(client *Client) {
	h.broadcastParticipantUpdate(client.roomID)
}

// handleSubmitTopic handles SUBMIT_TOPIC message
func (h *Handler) handleSubmitTopic(client *Client, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	var data SubmitTopicPayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		log.Printf("Error unmarshaling SUBMIT_TOPIC payload: %v", err)
		h.sendError(client, "INVALID_PAYLOAD", "Invalid SUBMIT_TOPIC payload")
		return
	}

	ctx := context.Background()

	// Execute use case to start discussion
	input := roomUseCase.StartDiscussionInput{
		RoomID:          client.roomID,
		OriginalEmojis:  data.OriginalEmojis,
		DisplayedEmojis: data.DisplayedEmojis,
		DummyIndex:      data.DummyIndex,
		DummyEmoji:      data.DummyEmoji,
	}

	if err := h.startDiscussionUseCase.Execute(ctx, input); err != nil {
		log.Printf("Error starting discussion: %v", err)
		h.sendError(client, "START_DISCUSSION_ERROR", err.Error())
		return
	}

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: client.roomID,
	})
	if err != nil {
		log.Printf("Error fetching room: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Broadcast state update
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

	h.hub.Broadcast(client.roomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: "discussing",
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

	// Start timer after 5 seconds delay
	h.timer.StartTimer(client.roomID)
}

// handleAnswering handles ANSWERING message
func (h *Handler) handleAnswering(client *Client, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	var data AnsweringPayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		log.Printf("Error unmarshaling ANSWERING payload: %v", err)
		h.sendError(client, "INVALID_PAYLOAD", "Invalid ANSWERING payload")
		return
	}

	ctx := context.Background()

	// Execute use case to submit final answer
	input := roomUseCase.SubmitFinalAnswerInput{
		RoomID:          client.roomID,
		Answer:          data.Answer,
		OriginalEmojis:  data.OriginalEmojis,
		DisplayedEmojis: data.DisplayedEmojis,
		DummyIndex:      data.DummyIndex,
		DummyEmoji:      data.DummyEmoji,
	}

	if err := h.submitFinalAnswerUseCase.Execute(ctx, input); err != nil {
		log.Printf("Error submitting final answer: %v", err)
		h.sendError(client, "SUBMIT_ANSWER_ERROR", err.Error())
		return
	}

	// Fetch room for broadcasting
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: client.roomID,
	})
	if err != nil {
		log.Printf("Error fetching room: %v", err)
		return
	}
	foundRoom := roomOutput.Room

	// Broadcast state update to checking
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

	h.hub.Broadcast(client.roomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: "checking",
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

	// Stop timer when transitioning to checking phase
	h.timer.StopTimer(client.roomID)
}

// sendError sends an error message to a specific client
func (h *Handler) sendError(client *Client, code string, message string) {
	errorMsg := Message{
		Type: MessageTypeError,
		Payload: ErrorPayload{
			Code:    code,
			Message: message,
		},
	}
	data, err := json.Marshal(errorMsg)
	if err != nil {
		log.Printf("Error marshaling error message: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		log.Printf("Failed to send error message to client")
	}
}

// broadcastParticipantUpdate broadcasts participant list to all clients in a room
func (h *Handler) broadcastParticipantUpdate(roomID string) {
	ctx := context.Background()

	// Execute use case to fetch participants
	input := roomUseCase.FetchRoomParticipantsInput{
		RoomID: roomID,
	}

	output, err := h.fetchParticipantsUseCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error fetching participants: %v", err)
		return
	}

	// Convert to WebSocket payload format
	participantDataList := []ParticipantData{}
	for _, p := range output.Participants {
		participantDataList = append(participantDataList, ParticipantData{
			UserID:   p.UserID,
			UserName: p.UserName,
			Role:     p.Role,
			IsLeader: p.IsLeader,
		})
	}

	h.hub.Broadcast(roomID, Message{
		Type: MessageTypeParticipantUpdate,
		Payload: ParticipantUpdatePayload{
			Participants: participantDataList,
		},
	})
}

// sendInitialRoomState sends the current room state to a newly connected client
func (h *Handler) sendInitialRoomState(client *Client) {
	ctx := context.Background()

	// Fetch room using UseCase
	roomOutput, err := h.fetchRoomUseCase.Execute(ctx, roomUseCase.FetchRoomInput{
		RoomID: client.roomID,
	})
	if err != nil {
		log.Printf("Error fetching room for initial state: %v", err)
		h.sendError(client, "FETCH_ROOM_ERROR", "Failed to fetch room state")
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

	// Create state update message with current room status
	stateMsg := Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: foundRoom.Status().String(),
			Data: &StateUpdateDataPayload{
				Topic:           topicStr,
				DisplayedEmojis: displayedEmojisSlice,
				OriginalEmojis:  originalEmojisSlice,
				DummyIndex:      dummyIdxPtr,
				DummyEmoji:      dummyEmojiStr,
				Assignments:     assignmentsSlice,
			},
		},
	}

	// Send state to the specific client
	data, err := json.Marshal(stateMsg)
	if err != nil {
		log.Printf("Error marshaling initial state message: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		log.Printf("Failed to send initial state to client")
	}
}
