package main

import (
	"bukidlink/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test user IDs
const (
	testWSUserJohn   = "d30869ec-fb97-46d8-85a3-82608c01f803"
	testWSUserDaniel = "c6554794-849f-4338-87c5-6db2e2f76514"
)

// drainChannel drains all messages from a channel
func drainChannel(ch chan []byte) {
	for {
		select {
		case <-ch:
			// Drained one message
		default:
			// Channel is empty
			return
		}
	}
}

// setupWebSocketTest initializes database and WebSocket hub for testing
func setupWebSocketTest(t *testing.T) {
	// Ensure database is initialized (always call to reset if needed)
	_ = db.SetupDatabase()

	// Initialize WebSocket hub
	InitWebSocket()

	// Give hub time to start
	time.Sleep(10 * time.Millisecond)
}

func TestWebSocketConnection(t *testing.T) {
	setupWebSocketTest(t)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock context with query params
		r.URL.RawQuery = "user_id=" + testWSUserJohn

		// Upgrade to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		// Create client
		client := &Client{
			ID:            generateClientID(),
			UserID:        testWSUserJohn,
			Conn:          conn,
			Hub:           hub,
			Send:          make(chan []byte, 256),
			Conversations: make(map[string]bool),
		}

		hub.Register <- client

		// Keep connection open briefly
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	// Connect as WebSocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?user_id=" + testWSUserJohn
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer ws.Close()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	// Verify client is registered
	hub.mu.RLock()
	clients := hub.Clients[testWSUserJohn]
	hub.mu.RUnlock()

	assert.NotNil(t, clients)
	assert.Equal(t, 1, len(clients))
}

func TestWebSocketMessageBroadcast(t *testing.T) {
	setupWebSocketTest(t)

	// Create a test conversation
	conv, err := db.CreateDirectConversation(testWSUserJohn, testWSUserDaniel)
	require.NoError(t, err)
	defer func() {
		db.GetDB().Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
		db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
		db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
	}()

	// Create mock clients
	client1Chan := make(chan []byte, 256)
	client2Chan := make(chan []byte, 256)

	client1 := &Client{
		ID:            "client1",
		UserID:        testWSUserJohn,
		Conn:          nil,
		Hub:           hub,
		Send:          client1Chan,
		Conversations: map[string]bool{conv.ID: true},
	}

	client2 := &Client{
		ID:            "client2",
		UserID:        testWSUserDaniel,
		Conn:          nil,
		Hub:           hub,
		Send:          client2Chan,
		Conversations: map[string]bool{conv.ID: true},
	}

	// Register clients
	hub.Register <- client1
	hub.Register <- client2

	time.Sleep(50 * time.Millisecond)

	// Drain any presence broadcasts from registration
	// (clients receive presence updates when other clients connect)
	drainChannel(client1Chan)
	drainChannel(client2Chan)

	// Create message payload
	msgPayload := WSMessagePayload{
		ConversationID: conv.ID,
		SenderID:       testWSUserJohn,
		Content:        "Test broadcast message",
		MessageType:    "text",
	}

	payloadBytes, _ := json.Marshal(msgPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeText,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	// Broadcast message
	hub.Broadcast <- &BroadcastMessage{
		ConversationID: conv.ID,
		Data:           data,
	}

	// Wait for broadcast
	time.Sleep(100 * time.Millisecond)

	// Both clients should receive the message
	select {
	case msg := <-client1Chan:
		assert.Contains(t, string(msg), "Test broadcast message")
	case <-time.After(time.Second):
		t.Fatal("Client 1 did not receive message")
	}

	select {
	case msg := <-client2Chan:
		assert.Contains(t, string(msg), "Test broadcast message")
	case <-time.After(time.Second):
		t.Fatal("Client 2 did not receive message")
	}

	// Cleanup
	hub.Unregister <- client1
	hub.Unregister <- client2
}

func TestWebSocketTypingIndicator(t *testing.T) {
	setupWebSocketTest(t)

	// Create a test conversation
	conv, err := db.CreateDirectConversation(testWSUserJohn, testWSUserDaniel)
	require.NoError(t, err)
	defer func() {
		db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
		db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
	}()

	// Create mock receiver client
	receiverChan := make(chan []byte, 256)
	receiver := &Client{
		ID:            "receiver",
		UserID:        testWSUserDaniel,
		Conn:          nil,
		Hub:           hub,
		Send:          receiverChan,
		Conversations: map[string]bool{conv.ID: true},
	}

	hub.Register <- receiver
	time.Sleep(50 * time.Millisecond)

	// Drain presence broadcast from registration
	drainChannel(receiverChan)

	// Send typing indicator
	typingPayload := WSTypingPayload{
		ConversationID: conv.ID,
		UserID:         testWSUserJohn,
	}

	payloadBytes, _ := json.Marshal(typingPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeTyping,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	hub.Broadcast <- &BroadcastMessage{
		ConversationID: conv.ID,
		ExcludeUserID:  testWSUserJohn,
		Data:           data,
	}

	// Wait for broadcast
	time.Sleep(100 * time.Millisecond)

	// Receiver should get typing indicator
	select {
	case msg := <-receiverChan:
		var receivedMsg WSMessage
		err := json.Unmarshal(msg, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, WSMessageTypeTyping, receivedMsg.Type)
	case <-time.After(time.Second):
		t.Fatal("Receiver did not get typing indicator")
	}

	// Cleanup
	hub.Unregister <- receiver
}

func TestWebSocketPresenceBroadcast(t *testing.T) {
	setupWebSocketTest(t)

	// Create mock clients
	client1Chan := make(chan []byte, 256)
	client2Chan := make(chan []byte, 256)

	client1 := &Client{
		ID:     "client1",
		UserID: testWSUserJohn,
		Conn:   nil,
		Hub:    hub,
		Send:   client1Chan,
	}

	client2 := &Client{
		ID:     "client2",
		UserID: testWSUserDaniel,
		Conn:   nil,
		Hub:    hub,
		Send:   client2Chan,
	}

	// Register first client
	hub.Register <- client1
	time.Sleep(50 * time.Millisecond)

	// Clear presence broadcast from registration
	select {
	case <-client1Chan:
	default:
	}

	// Register second client - should broadcast presence
	hub.Register <- client2
	time.Sleep(100 * time.Millisecond)

	// First client should receive presence update
	select {
	case msg := <-client1Chan:
		var wsMsg WSMessage
		err := json.Unmarshal(msg, &wsMsg)
		assert.NoError(t, err)
		assert.Equal(t, WSMessageTypePresence, wsMsg.Type)

		var presencePayload WSPresencePayload
		err = json.Unmarshal(wsMsg.Payload, &presencePayload)
		assert.NoError(t, err)
		assert.Equal(t, testWSUserDaniel, presencePayload.UserID)
		assert.Equal(t, "online", presencePayload.Status)
	case <-time.After(time.Second):
		t.Fatal("Client 1 did not receive presence update")
	}

	// Cleanup
	hub.Unregister <- client1
	hub.Unregister <- client2
}

func TestWebSocketExcludeSender(t *testing.T) {
	setupWebSocketTest(t)

	// Create a test conversation
	conv, err := db.CreateDirectConversation(testWSUserJohn, testWSUserDaniel)
	require.NoError(t, err)
	defer func() {
		db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
		db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
	}()

	// Create mock clients
	senderChan := make(chan []byte, 256)
	receiverChan := make(chan []byte, 256)

	sender := &Client{
		ID:     "sender",
		UserID: testWSUserJohn,
		Conn:   nil,
		Hub:    hub,
		Send:   senderChan,
	}

	receiver := &Client{
		ID:     "receiver",
		UserID: testWSUserDaniel,
		Conn:   nil,
		Hub:    hub,
		Send:   receiverChan,
	}

	hub.Register <- sender
	hub.Register <- receiver
	time.Sleep(50 * time.Millisecond)

	// Drain presence broadcasts
	drainChannel(senderChan)
	drainChannel(receiverChan)

	// Send message excluding sender
	msgPayload := WSMessagePayload{
		ConversationID: conv.ID,
		SenderID:       testWSUserJohn,
		Content:        "Test message",
	}

	payloadBytes, _ := json.Marshal(msgPayload)
	wsMessage := WSMessage{
		Type:    WSMessageTypeText,
		Payload: payloadBytes,
	}

	data, _ := json.Marshal(wsMessage)

	hub.Broadcast <- &BroadcastMessage{
		ConversationID: conv.ID,
		ExcludeUserID:  testWSUserJohn,
		Data:           data,
	}

	time.Sleep(100 * time.Millisecond)

	// Sender should NOT receive message
	select {
	case <-senderChan:
		t.Fatal("Sender should not receive message")
	default:
		// Expected
	}

	// Receiver should receive message
	select {
	case msg := <-receiverChan:
		assert.Contains(t, string(msg), "Test message")
	case <-time.After(time.Second):
		t.Fatal("Receiver did not receive message")
	}

	// Cleanup
	hub.Unregister <- sender
	hub.Unregister <- receiver
}

func TestWebSocketMultipleConnectionsSameUser(t *testing.T) {
	setupWebSocketTest(t)

	// Create two connections for same user
	client1Chan := make(chan []byte, 256)
	client2Chan := make(chan []byte, 256)

	client1 := &Client{
		ID:     "client1",
		UserID: testWSUserJohn,
		Conn:   nil,
		Hub:    hub,
		Send:   client1Chan,
	}

	client2 := &Client{
		ID:     "client2",
		UserID: testWSUserJohn, // Same user
		Conn:   nil,
		Hub:    hub,
		Send:   client2Chan,
	}

	hub.Register <- client1
	hub.Register <- client2
	time.Sleep(50 * time.Millisecond)

	// Verify both clients registered for same user
	hub.mu.RLock()
	clients := hub.Clients[testWSUserJohn]
	hub.mu.RUnlock()

	assert.Equal(t, 2, len(clients))

	// Unregister one client
	hub.Unregister <- client1
	time.Sleep(50 * time.Millisecond)

	// User should still be in hub (other connection active)
	hub.mu.RLock()
	clients = hub.Clients[testWSUserJohn]
	hub.mu.RUnlock()

	assert.Equal(t, 1, len(clients))

	// Verify presence is still "online" (not offline)
	presence, err := db.GetUserPresence(testWSUserJohn)
	require.NoError(t, err)
	assert.Equal(t, "online", presence.Status)

	// Cleanup
	hub.Unregister <- client2
}
