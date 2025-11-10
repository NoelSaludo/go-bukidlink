package main

import (
	"bukidlink/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test user IDs from actual database
const (
	testChatUserJohnDoe   = "d30869ec-fb97-46d8-85a3-82608c01f803"
	testChatUserDanielG   = "c6554794-849f-4338-87c5-6db2e2f76514"
	testChatUserSteward   = "6a24dd2b-d441-4b39-ab85-8fa2bd61065e"
	testChatUserMatthew   = "543255dd-5325-4d3f-bcd2-ee6f8ac87e2e"
	testChatUserMaryGrace = "9ae195a0-05ff-446b-99c0-e6f09a0150d1"
)

// ========== Direct Conversation Tests ==========

func TestPostDirectConversation(t *testing.T) {
	r := setupServer()

	req := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var conv db.Conversation
	err := json.Unmarshal(w.Body.Bytes(), &conv)
	require.NoError(t, err)
	assert.NotEmpty(t, conv.ID)
	assert.Equal(t, "direct", conv.Type)
	assert.Nil(t, conv.Title)

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestPostDirectConversation_SameUser(t *testing.T) {
	r := setupServer()

	req := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserJohnDoe, // Same user
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Cannot create conversation with yourself")
}

func TestPostDirectConversation_AlreadyExists(t *testing.T) {
	r := setupServer()

	// Create first conversation
	conv1, err := db.CreateDirectConversation(testChatUserJohnDoe, testChatUserSteward)
	require.NoError(t, err)

	// Try to create duplicate
	req := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserSteward,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code) // Should return existing conversation

	var conv db.Conversation
	err = json.Unmarshal(w.Body.Bytes(), &conv)
	require.NoError(t, err)
	assert.Equal(t, conv1.ID, conv.ID) // Should return the same conversation

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv1.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv1.ID)
}

// ========== Group Conversation Tests ==========

func TestPostGroupConversation(t *testing.T) {
	r := setupServer()

	req := CreateGroupConversationRequest{
		Title:          "Test Group Chat",
		CreatorID:      testChatUserJohnDoe,
		ParticipantIDs: []string{testChatUserDanielG, testChatUserSteward},
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/group", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var conv db.Conversation
	err := json.Unmarshal(w.Body.Bytes(), &conv)
	require.NoError(t, err)
	assert.NotEmpty(t, conv.ID)
	assert.Equal(t, "group", conv.Type)
	assert.NotNil(t, conv.Title)
	assert.Equal(t, "Test Group Chat", *conv.Title)

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestPostGroupConversation_NoParticipants(t *testing.T) {
	r := setupServer()

	req := CreateGroupConversationRequest{
		Title:          "Empty Group",
		CreatorID:      testChatUserJohnDoe,
		ParticipantIDs: []string{}, // Empty participants
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/group", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "must have at least one participant")
}

// ========== Get Conversation Tests ==========

func TestGetConversationByID(t *testing.T) {
	r := setupServer()

	// Create a test conversation
	conv, err := db.CreateDirectConversation(testChatUserMatthew, testChatUserMaryGrace)
	require.NoError(t, err)

	// Get the conversation
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/chat/conversation/"+conv.ID, nil)

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var result db.ConversationWithParticipants
	err = json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, conv.ID, result.ID)
	assert.Len(t, result.Participants, 2)
	assert.Contains(t, result.Participants, testChatUserMatthew)
	assert.Contains(t, result.Participants, testChatUserMaryGrace)

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetUserConversations(t *testing.T) {
	r := setupServer()

	// Create test conversations
	conv1, err := db.CreateDirectConversation(testChatUserJohnDoe, testChatUserDanielG)
	require.NoError(t, err)

	conv2, err := db.CreateGroupConversation("Test Group", testChatUserJohnDoe, []string{testChatUserSteward})
	require.NoError(t, err)

	// Get user's conversations
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/chat/user/"+testChatUserJohnDoe+"/conversations", nil)

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var conversations []db.ConversationWithParticipants
	err = json.Unmarshal(w.Body.Bytes(), &conversations)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(conversations), 2)

	// Find our created conversations
	foundConv1 := false
	foundConv2 := false
	for _, conv := range conversations {
		if conv.ID == conv1.ID {
			foundConv1 = true
			assert.Len(t, conv.Participants, 2)
		}
		if conv.ID == conv2.ID {
			foundConv2 = true
			assert.Len(t, conv.Participants, 2)
		}
	}
	assert.True(t, foundConv1, "Should find first conversation")
	assert.True(t, foundConv2, "Should find second conversation")

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id IN ($1, $2)`, conv1.ID, conv2.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id IN ($1, $2)`, conv1.ID, conv2.ID)
}

func TestGetUserConversations_WithPagination(t *testing.T) {
	r := setupServer()

	// Get user's conversations with pagination
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/chat/user/"+testChatUserJohnDoe+"/conversations?limit=5&offset=0", nil)

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var conversations []db.ConversationWithParticipants
	err := json.Unmarshal(w.Body.Bytes(), &conversations)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(conversations), 5) // Should respect limit
}

// ========== Participant Management Tests ==========

func TestAddParticipantToConversation(t *testing.T) {
	r := setupServer()

	// Create a group conversation
	conv, err := db.CreateGroupConversation("Test Group", testChatUserJohnDoe, []string{testChatUserDanielG})
	require.NoError(t, err)

	// Add a participant
	req := AddParticipantRequest{
		UserID: testChatUserSteward,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/participant", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify participant was added
	convDetails, err := db.GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convDetails.Participants, 3)
	assert.Contains(t, convDetails.Participants, testChatUserSteward)

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestAddParticipantToDirectConversation_ShouldFail(t *testing.T) {
	r := setupServer()

	// Create a direct conversation
	conv, err := db.CreateDirectConversation(testChatUserJohnDoe, testChatUserDanielG)
	require.NoError(t, err)

	// Try to add a participant
	req := AddParticipantRequest{
		UserID: testChatUserSteward,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/participant", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Cannot add participants to direct conversations")

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestRemoveParticipantFromConversation(t *testing.T) {
	r := setupServer()

	// Create a group conversation
	conv, err := db.CreateGroupConversation("Test Group", testChatUserJohnDoe, []string{testChatUserDanielG, testChatUserSteward})
	require.NoError(t, err)

	// Remove a participant
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/chat/conversation/"+conv.ID+"/participant/"+testChatUserDanielG, nil)

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify participant was removed
	convDetails, err := db.GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convDetails.Participants, 2)
	assert.NotContains(t, convDetails.Participants, testChatUserDanielG)

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestRemoveParticipantFromDirectConversation_ShouldFail(t *testing.T) {
	r := setupServer()

	// Create a direct conversation
	conv, err := db.CreateDirectConversation(testChatUserMatthew, testChatUserMaryGrace)
	require.NoError(t, err)

	// Try to remove a participant
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/chat/conversation/"+conv.ID+"/participant/"+testChatUserMatthew, nil)

	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Cannot remove participants from direct conversations")

	// Cleanup
	database := db.GetDB()
	_, _ = database.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = database.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

// ========== Message Endpoint Tests ==========

func TestSendMessage(t *testing.T) {
	router := setupServer()

	// First create a conversation
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Send a message
	msgReq := SendMessageRequest{
		SenderID:    testChatUserJohnDoe,
		Content:     "Hello from test!",
		MessageType: "text",
	}
	body, _ := json.Marshal(msgReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var msg db.Message
	err := json.Unmarshal(w.Body.Bytes(), &msg)
	assert.NoError(t, err)
	assert.Equal(t, "Hello from test!", msg.Content)
	assert.Equal(t, testChatUserJohnDoe, msg.SenderID)
	assert.Equal(t, "text", msg.MessageType)

	// Cleanup
	db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestSendMessage_NonParticipant(t *testing.T) {
	router := setupServer()

	// Create a conversation between JohnDoe and DanielG
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Try to send message as Steward (not a participant)
	msgReq := SendMessageRequest{
		SenderID: testChatUserSteward,
		Content:  "I shouldn't be able to send this",
	}
	body, _ := json.Marshal(msgReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	// Cleanup
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetConversationMessages(t *testing.T) {
	router := setupServer()

	// Create a conversation and send some messages
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Send 3 messages
	messageIDs := make([]string, 0)
	for i := 1; i <= 3; i++ {
		msgReq := SendMessageRequest{
			SenderID: testChatUserJohnDoe,
			Content:  "Message " + strconv.Itoa(i),
		}
		body, _ := json.Marshal(msgReq)
		msgResp := httptest.NewRecorder()
		msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
		msgReq2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(msgResp, msgReq2)

		var msg db.Message
		json.Unmarshal(msgResp.Body.Bytes(), &msg)
		messageIDs = append(messageIDs, msg.ID)
	}

	// Get messages
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/chat/conversation/"+conv.ID+"/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var messages []db.MessageWithSender
	err := json.Unmarshal(w.Body.Bytes(), &messages)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(messages), 3)

	// Verify sender information is populated
	assert.NotEmpty(t, messages[0].SenderUsername)

	// Cleanup
	for _, msgID := range messageIDs {
		db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msgID)
	}
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetConversationMessages_WithPagination(t *testing.T) {
	router := setupServer()

	// Create a conversation
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Send 5 messages
	messageIDs := make([]string, 0)
	for i := 1; i <= 5; i++ {
		msgReq := SendMessageRequest{
			SenderID: testChatUserJohnDoe,
			Content:  "Message " + strconv.Itoa(i),
		}
		body, _ := json.Marshal(msgReq)
		msgResp := httptest.NewRecorder()
		msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
		msgReq2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(msgResp, msgReq2)

		var msg db.Message
		json.Unmarshal(msgResp.Body.Bytes(), &msg)
		messageIDs = append(messageIDs, msg.ID)
	}

	// Get first 2 messages
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/chat/conversation/"+conv.ID+"/messages?limit=2&offset=0", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var messages []db.MessageWithSender
	err := json.Unmarshal(w.Body.Bytes(), &messages)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(messages), 2)

	// Cleanup
	for _, msgID := range messageIDs {
		db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msgID)
	}
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestEditMessage(t *testing.T) {
	router := setupServer()

	// Create conversation and send message
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	msgReq := SendMessageRequest{
		SenderID: testChatUserJohnDoe,
		Content:  "Original message",
	}
	msgBody, _ := json.Marshal(msgReq)
	msgResp := httptest.NewRecorder()
	msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(msgBody))
	msgReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(msgResp, msgReq2)

	var msg db.Message
	json.Unmarshal(msgResp.Body.Bytes(), &msg)

	// Edit the message
	editReq := EditMessageRequest{
		Content: "Edited message",
	}
	editBody, _ := json.Marshal(editReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/chat/message/"+msg.ID, bytes.NewBuffer(editBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the edit in database
	var editedContent string
	var editedAt *time.Time
	err := db.GetDB().QueryRow(`SELECT content, edited_at FROM "Message" WHERE id = $1`, msg.ID).Scan(&editedContent, &editedAt)
	assert.NoError(t, err)
	assert.Equal(t, "Edited message", editedContent)
	assert.NotNil(t, editedAt)

	// Cleanup
	db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestDeleteMessage(t *testing.T) {
	router := setupServer()

	// Create conversation and send message
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	msgReq := SendMessageRequest{
		SenderID: testChatUserJohnDoe,
		Content:  "Message to delete",
	}
	msgBody, _ := json.Marshal(msgReq)
	msgResp := httptest.NewRecorder()
	msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(msgBody))
	msgReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(msgResp, msgReq2)

	var msg db.Message
	json.Unmarshal(msgResp.Body.Bytes(), &msg)

	// Delete the message
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/chat/message/"+msg.ID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the deletion in database
	var isDeleted bool
	var content string
	err := db.GetDB().QueryRow(`SELECT is_deleted, content FROM "Message" WHERE id = $1`, msg.ID).Scan(&isDeleted, &content)
	assert.NoError(t, err)
	assert.True(t, isDeleted)
	assert.Equal(t, "[Message deleted]", content)

	// Cleanup
	db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

// ========== Read Receipt Tests ==========

func TestGetUnreadCount(t *testing.T) {
	router := setupServer()

	// Create conversation and send messages
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Send 3 messages from JohnDoe
	messageIDs := make([]string, 0)
	for i := 1; i <= 3; i++ {
		msgReq := SendMessageRequest{
			SenderID: testChatUserJohnDoe,
			Content:  "Unread message " + strconv.Itoa(i),
		}
		body, _ := json.Marshal(msgReq)
		msgResp := httptest.NewRecorder()
		msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
		msgReq2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(msgResp, msgReq2)

		var msg db.Message
		json.Unmarshal(msgResp.Body.Bytes(), &msg)
		messageIDs = append(messageIDs, msg.ID)
	}

	// Check unread count for DanielG
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/chat/conversation/"+conv.ID+"/unread?user_id="+testChatUserDanielG, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, response["unread_count"], 3)

	// Cleanup
	for _, msgID := range messageIDs {
		db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msgID)
	}
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestMarkAsRead(t *testing.T) {
	router := setupServer()

	// Create conversation and send messages
	convReq := CreateDirectConversationRequest{
		UserID1: testChatUserJohnDoe,
		UserID2: testChatUserDanielG,
	}
	convBody, _ := json.Marshal(convReq)
	convResp := httptest.NewRecorder()
	convReq2, _ := http.NewRequest("POST", "/chat/conversation/direct", bytes.NewBuffer(convBody))
	convReq2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(convResp, convReq2)

	var conv db.Conversation
	json.Unmarshal(convResp.Body.Bytes(), &conv)

	// Send messages
	messageIDs := make([]string, 0)
	for i := 1; i <= 2; i++ {
		msgReq := SendMessageRequest{
			SenderID: testChatUserJohnDoe,
			Content:  "Message " + strconv.Itoa(i),
		}
		body, _ := json.Marshal(msgReq)
		msgResp := httptest.NewRecorder()
		msgReq2, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/message", bytes.NewBuffer(body))
		msgReq2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(msgResp, msgReq2)

		var msg db.Message
		json.Unmarshal(msgResp.Body.Bytes(), &msg)
		messageIDs = append(messageIDs, msg.ID)
	}

	// Mark as read
	readReq := MarkAsReadRequest{
		UserID: testChatUserDanielG,
	}
	readBody, _ := json.Marshal(readReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/conversation/"+conv.ID+"/read", bytes.NewBuffer(readBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify last_read_at was updated
	var lastReadAt *time.Time
	err := db.GetDB().QueryRow(`SELECT last_read_at FROM "ConversationParticipant" WHERE conversation_id = $1 AND user_id = $2`, conv.ID, testChatUserDanielG).Scan(&lastReadAt)
	assert.NoError(t, err)
	assert.NotNil(t, lastReadAt)

	// Cleanup
	for _, msgID := range messageIDs {
		db.GetDB().Exec(`DELETE FROM "Message" WHERE id = $1`, msgID)
	}
	db.GetDB().Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	db.GetDB().Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

// ========== User Presence Tests ==========

func TestUpdateUserPresence(t *testing.T) {
	router := setupServer()

	presenceReq := struct {
		Status string `json:"status"`
	}{
		Status: "online",
	}

	body, _ := json.Marshal(presenceReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/user/"+testChatUserJohnDoe+"/presence", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify in database
	var status string
	err := db.GetDB().QueryRow(`SELECT status FROM "UserPresence" WHERE user_id = $1`, testChatUserJohnDoe).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "online", status)
}

func TestUpdateUserPresence_InvalidStatus(t *testing.T) {
	router := setupServer()

	presenceReq := struct {
		Status string `json:"status"`
	}{
		Status: "invisible", // Invalid status
	}

	body, _ := json.Marshal(presenceReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/user/"+testChatUserJohnDoe+"/presence", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserPresence(t *testing.T) {
	router := setupServer()

	// First set presence
	db.UpsertUserPresence(testChatUserJohnDoe, "online")

	// Get presence
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/chat/user/"+testChatUserJohnDoe+"/presence", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var presence db.UserPresence
	err := json.Unmarshal(w.Body.Bytes(), &presence)
	assert.NoError(t, err)
	assert.Equal(t, testChatUserJohnDoe, presence.UserID)
	assert.Equal(t, "online", presence.Status)
}

func TestGetMultipleUserPresence(t *testing.T) {
	router := setupServer()

	// Set presence for multiple users
	db.UpsertUserPresence(testChatUserJohnDoe, "online")
	db.UpsertUserPresence(testChatUserDanielG, "away")
	db.UpsertUserPresence(testChatUserSteward, "offline")

	presenceReq := struct {
		UserIDs []string `json:"user_ids"`
	}{
		UserIDs: []string{testChatUserJohnDoe, testChatUserDanielG, testChatUserSteward},
	}

	body, _ := json.Marshal(presenceReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat/presence/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var presences []db.UserPresence
	err := json.Unmarshal(w.Body.Bytes(), &presences)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(presences), 3)
}
