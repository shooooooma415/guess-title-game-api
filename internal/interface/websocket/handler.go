package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
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
	hub             *Hub
	timer           *Timer
	roomRepo        room.Repository
	participantRepo participant.Repository
	userRepo        user.Repository
}

// NewHandler creates a new WebSocket handler
func NewHandler(
	hub *Hub,
	timer *Timer,
	roomRepo room.Repository,
	participantRepo participant.Repository,
	userRepo user.Repository,
) *Handler {
	return &Handler{
		hub:             hub,
		timer:           timer,
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
		userRepo:        userRepo,
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
	}
}

// handleClientConnected handles CLIENT_CONNECTED message
func (h *Handler) handleClientConnected(client *Client, payload interface{}) {
	payloadBytes, _ := json.Marshal(payload)
	var data ClientConnectedPayload
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		log.Printf("Error unmarshaling CLIENT_CONNECTED payload: %v", err)
		return
	}

	client.userID = data.UserID

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
		return
	}

	ctx := context.Background()

	// Find room
	roomID, err := room.NewRoomIDFromString(client.roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		return
	}

	foundRoom, err := h.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		log.Printf("Room not found: %v", err)
		return
	}

	// Save game data
	originalEmojis := room.NewEmojiList(data.OriginalEmojis)
	displayedEmojis := room.NewEmojiList(data.DisplayedEmojis)
	dummyIndex, _ := room.NewDummyIndex(data.DummyIndex)
	dummyEmoji, _ := room.NewDummyEmoji(data.DummyEmoji)

	foundRoom.SetGameData(
		originalEmojis,
		displayedEmojis,
		dummyIndex,
		dummyEmoji,
	)

	if err := h.roomRepo.Save(ctx, foundRoom); err != nil {
		log.Printf("Error saving room: %v", err)
		return
	}

	// Change status to discussing
	foundRoom.ChangeStatus(room.StatusDiscussing)
	h.roomRepo.Save(ctx, foundRoom)

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
		return
	}

	ctx := context.Background()

	// Find room
	roomID, err := room.NewRoomIDFromString(client.roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		return
	}

	foundRoom, err := h.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		log.Printf("Room not found: %v", err)
		return
	}

	// Save all data
	answer, _ := room.NewAnswer(data.Answer)
	foundRoom.SetAnswer(answer)

	originalEmojis := room.NewEmojiList(data.OriginalEmojis)
	displayedEmojis := room.NewEmojiList(data.DisplayedEmojis)
	dummyIndex, _ := room.NewDummyIndex(data.DummyIndex)
	dummyEmoji, _ := room.NewDummyEmoji(data.DummyEmoji)

	foundRoom.SetGameData(
		originalEmojis,
		displayedEmojis,
		dummyIndex,
		dummyEmoji,
	)
	foundRoom.ChangeStatus(room.StatusChecking)

	if err := h.roomRepo.Save(ctx, foundRoom); err != nil {
		log.Printf("Error saving room: %v", err)
		return
	}

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

	displayedEmojisSlice2 := []string{}
	if foundRoom.DisplayedEmojis() != nil {
		displayedEmojisSlice2 = foundRoom.DisplayedEmojis().Values()
	}

	originalEmojisSlice2 := []string{}
	if foundRoom.OriginalEmojis() != nil {
		originalEmojisSlice2 = foundRoom.OriginalEmojis().Values()
	}

	assignmentsSlice2 := []string{}
	if foundRoom.Assignments() != nil {
		assignmentsSlice2 = foundRoom.Assignments().Values()
	}

	h.hub.Broadcast(client.roomID, Message{
		Type: MessageTypeStateUpdate,
		Payload: StateUpdatePayload{
			NextState: "checking",
			Data: &StateUpdateDataPayload{
				Topic:           topicStr,
				DisplayedEmojis: displayedEmojisSlice2,
				OriginalEmojis:  originalEmojisSlice2,
				DummyIndex:      dummyIdxPtr,
				DummyEmoji:      dummyEmojiStr,
				Assignments:     assignmentsSlice2,
			},
		},
	})
}

// broadcastParticipantUpdate broadcasts participant list to all clients in a room
func (h *Handler) broadcastParticipantUpdate(roomID string) {
	ctx := context.Background()

	participantRoomID, err := participant.NewRoomIDFromString(roomID)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		return
	}

	participants, err := h.participantRepo.FindByRoomID(ctx, participantRoomID)
	if err != nil {
		log.Printf("Error fetching participants: %v", err)
		return
	}

	participantDataList := []ParticipantData{}
	for _, p := range participants {
		// Fetch user info
		u, err := h.userRepo.FindByID(ctx, user.UserID{})
		userName := "Unknown"
		if err == nil {
			userName = u.Name().String()
		}

		participantDataList = append(participantDataList, ParticipantData{
			UserID:   p.UserID().String(),
			UserName: userName,
			Role:     p.Role().String(),
			IsLeader: p.IsLeader(),
		})
	}

	h.hub.Broadcast(roomID, Message{
		Type: MessageTypeParticipantUpdate,
		Payload: ParticipantUpdatePayload{
			Participants: participantDataList,
		},
	})
}
