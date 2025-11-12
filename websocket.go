package main

import (
	"bukidlink/db"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket message types
const (
	WSMessageTypeText          = "message"
	WSMessageTypeTyping        = "typing"
	WSMessageTypeStopTyping    = "stop_typing"
	WSMessageTypePresence      = "presence"
	WSMessageTypeRead          = "read"
	WSMessageTypeMessageEdit   = "message_edit"
	WSMessageTypeMessageDelete = "message_delete"
	WSMessageTypeError         = "error"
	WSMessageTypePing          = "ping"
	WSMessageTypePong          = "pong"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: In production, validate origin properly
		return true
	},
}

// WSMessage represents a WebSocket message envelope
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// WSMessagePayload represents the payload for a new message
type WSMessagePayload struct {
	ConversationID string  `json:"conversation_id"`
	MessageID      string  `json:"message_id,omitempty"` // For edits/deletes
	SenderID       string  `json:"sender_id"`
	Content        string  `json:"content"`
	MessageType    string  `json:"message_type,omitempty"`
	AttachmentURL  *string `json:"attachment_url,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
}

// WSTypingPayload represents typing indicator payload
type WSTypingPayload struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	Username       string `json:"username,omitempty"`
}

// WSPresencePayload represents presence update payload
type WSPresencePayload struct {
	UserID string `json:"user_id"`
	Status string `json:"status"` // online, offline, away
}

// WSReadPayload represents read receipt payload
type WSReadPayload struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	Timestamp      string `json:"timestamp"`
}

// WSErrorPayload represents error message payload
type WSErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID            string
	UserID        string
	Conn          *websocket.Conn
	Hub           *Hub
	Send          chan []byte
	Conversations map[string]bool // Set of conversation IDs user is part of
	mu            sync.RWMutex
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients mapped by user ID
	Clients map[string]map[*Client]bool

	// Broadcast channel for messages
	Broadcast chan *BroadcastMessage

	// Register client
	Register chan *Client

	// Unregister client
	Unregister chan *Client

	mu sync.RWMutex
}

// BroadcastMessage represents a message to broadcast
type BroadcastMessage struct {
	ConversationID string
	UserIDs        []string // If specified, only send to these users
	ExcludeUserID  string   // Exclude this user (e.g., sender)
	Data           []byte
}

// Global hub instance
var hub *Hub

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Broadcast:  make(chan *BroadcastMessage, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = make(map[*Client]bool)
			}
			h.Clients[client.UserID][client] = true
			h.mu.Unlock()

			// Update presence to online
			db.UpsertUserPresence(client.UserID, "online")

			// Broadcast presence update to all users
			h.broadcastPresence(client.UserID, "online")

			log.Printf("Client connected: %s (user: %s)", client.ID, client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.UserID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					// If no more clients for this user, update presence
					if len(clients) == 0 {
						delete(h.Clients, client.UserID)
						h.mu.Unlock()

						// Update presence to offline
						db.UpsertUserPresence(client.UserID, "offline")
						h.broadcastPresence(client.UserID, "offline")

						log.Printf("Client disconnected: %s (user: %s) - last connection", client.ID, client.UserID)
					} else {
						h.mu.Unlock()
						log.Printf("Client disconnected: %s (user: %s) - other connections active", client.ID, client.UserID)
					}
				} else {
					h.mu.Unlock()
				}
			} else {
				h.mu.Unlock()
			}

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// broadcastMessage sends a message to relevant clients
func (h *Hub) broadcastMessage(msg *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Get all participants in the conversation if conversation ID is specified
	var targetUserIDs []string
	if msg.ConversationID != "" {
		conv, err := db.GetConversationByID(msg.ConversationID)
		if err != nil {
			log.Printf("Error getting conversation participants: %v", err)
			return
		}
		targetUserIDs = conv.Participants
	} else if len(msg.UserIDs) > 0 {
		targetUserIDs = msg.UserIDs
	}

	// Send to target users
	for _, userID := range targetUserIDs {
		if userID == msg.ExcludeUserID {
			continue
		}

		if clients, ok := h.Clients[userID]; ok {
			for client := range clients {
				select {
				case client.Send <- msg.Data:
				default:
					// Client's send buffer is full, disconnect
					h.mu.RUnlock()
					h.Unregister <- client
					h.mu.RLock()
				}
			}
		}
	}
}

// broadcastPresence sends presence update to all connected users
func (h *Hub) broadcastPresence(userID, status string) {
	payload := WSPresencePayload{
		UserID: userID,
		Status: status,
	}

	payloadBytes, _ := json.Marshal(payload)
	message := WSMessage{
		Type:    WSMessageTypePresence,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(message)

	// Broadcast to all connected users
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.Clients {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
				// Skip if send buffer is full
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse message
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			c.sendError("Invalid message format")
			continue
		}

		// Handle message based on type
		c.handleMessage(&wsMsg)
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(msg *WSMessage) {
	switch msg.Type {
	case WSMessageTypeText:
		c.handleNewMessage(msg.Payload)

	case WSMessageTypeTyping:
		c.handleTyping(msg.Payload)

	case WSMessageTypeStopTyping:
		c.handleStopTyping(msg.Payload)

	case WSMessageTypePresence:
		c.handlePresenceUpdate(msg.Payload)

	case WSMessageTypeRead:
		c.handleReadReceipt(msg.Payload)

	case WSMessageTypeMessageEdit:
		c.handleMessageEdit(msg.Payload)

	case WSMessageTypeMessageDelete:
		c.handleMessageDelete(msg.Payload)

	case WSMessageTypePong:
		// Client responded to ping, do nothing

	default:
		c.sendError("Unknown message type: " + msg.Type)
	}
}

// handleNewMessage processes a new message
func (c *Client) handleNewMessage(payload json.RawMessage) {
	var msgPayload WSMessagePayload
	if err := json.Unmarshal(payload, &msgPayload); err != nil {
		c.sendError("Invalid message payload")
		return
	}

	// Verify sender matches client user
	if msgPayload.SenderID != c.UserID {
		c.sendError("Sender ID does not match authenticated user")
		return
	}

	// Default message type
	if msgPayload.MessageType == "" {
		msgPayload.MessageType = "text"
	}

	// Insert message into database
	message, err := db.InsertMessage(
		msgPayload.ConversationID,
		msgPayload.SenderID,
		msgPayload.Content,
		msgPayload.MessageType,
		msgPayload.AttachmentURL,
	)
	if err != nil {
		c.sendError("Failed to save message: " + err.Error())
		return
	}

	// Prepare broadcast payload
	broadcastPayload := WSMessagePayload{
		ConversationID: message.ConversationID,
		MessageID:      message.ID,
		SenderID:       message.SenderID,
		Content:        message.Content,
		MessageType:    message.MessageType,
		AttachmentURL:  message.AttachmentURL,
		CreatedAt:      message.CreatedAt.Format(time.RFC3339),
	}

	payloadBytes, _ := json.Marshal(broadcastPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeText,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	// Broadcast to conversation participants
	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: msgPayload.ConversationID,
		Data:           data,
	}
}

// handleTyping broadcasts typing indicator
func (c *Client) handleTyping(payload json.RawMessage) {
	var typingPayload WSTypingPayload
	if err := json.Unmarshal(payload, &typingPayload); err != nil {
		c.sendError("Invalid typing payload")
		return
	}

	typingPayload.UserID = c.UserID

	payloadBytes, _ := json.Marshal(typingPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeTyping,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	// Broadcast to conversation participants except sender
	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: typingPayload.ConversationID,
		ExcludeUserID:  c.UserID,
		Data:           data,
	}
}

// handleStopTyping broadcasts stop typing indicator
func (c *Client) handleStopTyping(payload json.RawMessage) {
	var typingPayload WSTypingPayload
	if err := json.Unmarshal(payload, &typingPayload); err != nil {
		c.sendError("Invalid stop typing payload")
		return
	}

	typingPayload.UserID = c.UserID

	payloadBytes, _ := json.Marshal(typingPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeStopTyping,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	// Broadcast to conversation participants except sender
	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: typingPayload.ConversationID,
		ExcludeUserID:  c.UserID,
		Data:           data,
	}
}

// handlePresenceUpdate updates user presence
func (c *Client) handlePresenceUpdate(payload json.RawMessage) {
	var presencePayload WSPresencePayload
	if err := json.Unmarshal(payload, &presencePayload); err != nil {
		c.sendError("Invalid presence payload")
		return
	}

	// Verify user can only update their own presence
	if presencePayload.UserID != c.UserID {
		c.sendError("Cannot update another user's presence")
		return
	}

	// Update presence in database
	if err := db.UpsertUserPresence(c.UserID, presencePayload.Status); err != nil {
		c.sendError("Failed to update presence")
		return
	}

	// Broadcast presence update
	c.Hub.broadcastPresence(c.UserID, presencePayload.Status)
}

// handleReadReceipt updates read receipt
func (c *Client) handleReadReceipt(payload json.RawMessage) {
	var readPayload WSReadPayload
	if err := json.Unmarshal(payload, &readPayload); err != nil {
		c.sendError("Invalid read receipt payload")
		return
	}

	// Verify user matches
	if readPayload.UserID != c.UserID {
		c.sendError("Cannot update read receipt for another user")
		return
	}

	// Update last read timestamp
	if err := db.UpdateLastReadAt(readPayload.ConversationID, c.UserID); err != nil {
		c.sendError("Failed to update read receipt")
		return
	}

	// Broadcast read receipt to conversation participants
	readPayload.Timestamp = time.Now().Format(time.RFC3339)
	payloadBytes, _ := json.Marshal(readPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeRead,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: readPayload.ConversationID,
		ExcludeUserID:  c.UserID,
		Data:           data,
	}
}

// handleMessageEdit processes message edit
func (c *Client) handleMessageEdit(payload json.RawMessage) {
	var editPayload WSMessagePayload
	if err := json.Unmarshal(payload, &editPayload); err != nil {
		c.sendError("Invalid edit payload")
		return
	}

	// Update message in database
	if err := db.UpdateMessage(editPayload.MessageID, editPayload.Content); err != nil {
		c.sendError("Failed to edit message")
		return
	}

	// Broadcast edit to conversation participants
	payloadBytes, _ := json.Marshal(editPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeMessageEdit,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: editPayload.ConversationID,
		Data:           data,
	}
}

// handleMessageDelete processes message deletion
func (c *Client) handleMessageDelete(payload json.RawMessage) {
	var deletePayload WSMessagePayload
	if err := json.Unmarshal(payload, &deletePayload); err != nil {
		c.sendError("Invalid delete payload")
		return
	}

	// Delete message in database
	if err := db.DeleteMessage(deletePayload.MessageID); err != nil {
		c.sendError("Failed to delete message")
		return
	}

	// Broadcast deletion to conversation participants
	payloadBytes, _ := json.Marshal(deletePayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeMessageDelete,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	c.Hub.Broadcast <- &BroadcastMessage{
		ConversationID: deletePayload.ConversationID,
		Data:           data,
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(message string) {
	errorPayload := WSErrorPayload{
		Message: message,
	}

	payloadBytes, _ := json.Marshal(errorPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeError,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	select {
	case c.Send <- data:
	default:
		// Send buffer full, skip error message
	}
}

// HandleWebSocket handles WebSocket connection upgrade
func HandleWebSocket(c *gin.Context) {
	// Get user ID from query parameter
	// TODO: In production, get user ID from JWT token in auth middleware
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Log connection attempt with client info
	log.Printf("WebSocket connection attempt from %s (user_id: %s)", c.ClientIP(), userID)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create client
	client := &Client{
		ID:            generateClientID(),
		UserID:        userID,
		Conn:          conn,
		Hub:           hub,
		Send:          make(chan []byte, 256),
		Conversations: make(map[string]bool),
	}

	// Register client
	hub.Register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return uuid.New().String()
}

// InitWebSocket initializes the WebSocket hub
func InitWebSocket() {
	hub = NewHub()
	go hub.Run()
	log.Println("WebSocket hub initialized")
}
