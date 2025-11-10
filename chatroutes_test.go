package main

import (
	"bukidlink/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
