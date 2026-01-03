package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
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

// Placeholder handlers - implement these based on your business logic
func (h *Handler) handleClientConnected(client *Client, payload interface{}) {
	// TODO: Implement
	log.Printf("Client connected: %v", payload)
}

func (h *Handler) handleFetchParticipants(client *Client) {
	// TODO: Implement
	log.Printf("Fetch participants")
}

func (h *Handler) handleSubmitTopic(client *Client, payload interface{}) {
	// TODO: Implement
	log.Printf("Submit topic: %v", payload)
}

func (h *Handler) handleAnswering(client *Client, payload interface{}) {
	// TODO: Implement
	log.Printf("Answering: %v", payload)
}
