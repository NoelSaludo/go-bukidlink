package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test user IDs from actual User table (User.csv)
const (
	testUserJohnDoe       = "d30869ec-fb97-46d8-85a3-82608c01f803"
	testUserDanielG       = "c6554794-849f-4338-87c5-6db2e2f76514"
	testUserStewardLittle = "6a24dd2b-d441-4b39-ab85-8fa2bd61065e"
	testUserMatthew       = "543255dd-5325-4d3f-bcd2-ee6f8ac87e2e"
	testUserMaryGrace     = "9ae195a0-05ff-446b-99c0-e6f09a0150d1"
	// Legacy aliases for backward compatibility
	testUserMariaGarcia  = testUserStewardLittle
	testUserRobertSmith  = testUserMatthew
	testUserSarahJohnson = testUserMaryGrace
)

// ========== Conversation Tests ==========

func TestCreateDirectConversation(t *testing.T) {
	_ = SetupDatabase()

	// Create a direct conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)

	require.NoError(t, err)
	assert.NotEmpty(t, conv.ID)
	assert.Nil(t, conv.Title)
	assert.Equal(t, "direct", conv.Type)
	assert.NotZero(t, conv.CreatedAt)
	assert.NotZero(t, conv.UpdatedAt)

	// Verify participants were added
	convWithParticipants, err := GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convWithParticipants.Participants, 2)
	assert.Contains(t, convWithParticipants.Participants, testUserJohnDoe)
	assert.Contains(t, convWithParticipants.Participants, testUserDanielG)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestCreateGroupConversation(t *testing.T) {
	_ = SetupDatabase()

	title := "Test Group Chat"
	participants := []string{testUserDanielG, testUserMariaGarcia, testUserRobertSmith}

	conv, err := CreateGroupConversation(title, testUserJohnDoe, participants)

	require.NoError(t, err)
	assert.NotEmpty(t, conv.ID)
	assert.NotNil(t, conv.Title)
	assert.Equal(t, title, *conv.Title)
	assert.Equal(t, "group", conv.Type)
	assert.NotZero(t, conv.CreatedAt)
	assert.NotZero(t, conv.UpdatedAt)

	// Verify all participants were added (creator + participants)
	convWithParticipants, err := GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convWithParticipants.Participants, 4)
	assert.Contains(t, convWithParticipants.Participants, testUserJohnDoe)
	assert.Contains(t, convWithParticipants.Participants, testUserDanielG)
	assert.Contains(t, convWithParticipants.Participants, testUserMariaGarcia)
	assert.Contains(t, convWithParticipants.Participants, testUserRobertSmith)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetDirectConversation(t *testing.T) {
	_ = SetupDatabase()

	// Create a direct conversation
	originalConv, err := CreateDirectConversation(testUserJohnDoe, testUserSarahJohnson)
	require.NoError(t, err)

	// Test finding existing conversation
	foundConv, err := GetDirectConversation(testUserJohnDoe, testUserSarahJohnson)
	require.NoError(t, err)
	assert.Equal(t, originalConv.ID, foundConv.ID)

	// Test reverse order
	foundConvReverse, err := GetDirectConversation(testUserSarahJohnson, testUserJohnDoe)
	require.NoError(t, err)
	assert.Equal(t, originalConv.ID, foundConvReverse.ID)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, originalConv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, originalConv.ID)
}

func TestGetConversationByID(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation
	originalConv, err := CreateDirectConversation(testUserDanielG, testUserMariaGarcia)
	require.NoError(t, err)

	// Get conversation by ID
	conv, err := GetConversationByID(originalConv.ID)

	require.NoError(t, err)
	assert.Equal(t, originalConv.ID, conv.ID)
	assert.Equal(t, originalConv.Type, conv.Type)
	assert.Len(t, conv.Participants, 2)
	assert.Contains(t, conv.Participants, testUserDanielG)
	assert.Contains(t, conv.Participants, testUserMariaGarcia)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, originalConv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, originalConv.ID)
}

func TestGetUserConversations(t *testing.T) {
	_ = SetupDatabase()

	// Create multiple conversations for a user
	conv1, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	conv2, err := CreateDirectConversation(testUserJohnDoe, testUserMariaGarcia)
	require.NoError(t, err)

	conv3, err := CreateGroupConversation("Test Group", testUserJohnDoe, []string{testUserDanielG, testUserRobertSmith})
	require.NoError(t, err)

	// Get all conversations for the user
	conversations, err := GetUserConversations(testUserJohnDoe, 10, 0)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(conversations), 3)

	// Find our created conversations
	foundIDs := make(map[string]bool)
	for _, conv := range conversations {
		if conv.ID == conv1.ID || conv.ID == conv2.ID || conv.ID == conv3.ID {
			foundIDs[conv.ID] = true
			assert.NotEmpty(t, conv.Participants)
			assert.Contains(t, conv.Participants, testUserJohnDoe)
		}
	}

	assert.True(t, foundIDs[conv1.ID], "Should find first direct conversation")
	assert.True(t, foundIDs[conv2.ID], "Should find second direct conversation")
	assert.True(t, foundIDs[conv3.ID], "Should find group conversation")

	// Cleanup
	for _, convID := range []string{conv1.ID, conv2.ID, conv3.ID} {
		_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, convID)
		_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, convID)
	}
}

func TestAddParticipantToConversation(t *testing.T) {
	_ = SetupDatabase()

	// Create a group conversation
	conv, err := CreateGroupConversation("Test Group", testUserJohnDoe, []string{testUserDanielG})
	require.NoError(t, err)

	// Add a new participant
	err = AddParticipantToConversation(conv.ID, testUserMariaGarcia)
	require.NoError(t, err)

	// Verify participant was added
	convWithParticipants, err := GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convWithParticipants.Participants, 3)
	assert.Contains(t, convWithParticipants.Participants, testUserMariaGarcia)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestRemoveParticipantFromConversation(t *testing.T) {
	_ = SetupDatabase()

	// Create a group conversation
	conv, err := CreateGroupConversation("Test Group", testUserJohnDoe, []string{testUserDanielG, testUserMariaGarcia})
	require.NoError(t, err)

	// Remove a participant
	err = RemoveParticipantFromConversation(conv.ID, testUserDanielG)
	require.NoError(t, err)

	// Verify participant was removed
	convWithParticipants, err := GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convWithParticipants.Participants, 2)
	assert.NotContains(t, convWithParticipants.Participants, testUserDanielG)
	assert.Contains(t, convWithParticipants.Participants, testUserJohnDoe)
	assert.Contains(t, convWithParticipants.Participants, testUserMariaGarcia)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

// ========== Message Tests ==========

func TestInsertMessage(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation first
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	// Insert a text message
	msg, err := InsertMessage(conv.ID, testUserJohnDoe, "Hello, Daniel!", "text", nil)

	require.NoError(t, err)
	assert.NotEmpty(t, msg.ID)
	assert.Equal(t, conv.ID, msg.ConversationID)
	assert.Equal(t, testUserJohnDoe, msg.SenderID)
	assert.Equal(t, "Hello, Daniel!", msg.Content)
	assert.Equal(t, "text", msg.MessageType)
	assert.Nil(t, msg.AttachmentURL)
	assert.NotZero(t, msg.CreatedAt)
	assert.Nil(t, msg.EditedAt)
	assert.False(t, msg.IsDeleted)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestInsertMessage_WithAttachment(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserMariaGarcia)
	require.NoError(t, err)

	// Insert an image message with attachment
	attachmentURL := "https://example.com/images/photo.jpg"
	msg, err := InsertMessage(conv.ID, testUserJohnDoe, "Check out this photo!", "image", &attachmentURL)

	require.NoError(t, err)
	assert.NotEmpty(t, msg.ID)
	assert.Equal(t, "image", msg.MessageType)
	assert.NotNil(t, msg.AttachmentURL)
	assert.Equal(t, attachmentURL, *msg.AttachmentURL)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetMessagesByConversation(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	// Insert multiple messages
	msg1, err := InsertMessage(conv.ID, testUserJohnDoe, "First message", "text", nil)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond) // Ensure different timestamps

	msg2, err := InsertMessage(conv.ID, testUserDanielG, "Second message", "text", nil)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	msg3, err := InsertMessage(conv.ID, testUserJohnDoe, "Third message", "text", nil)
	require.NoError(t, err)

	// Get messages
	messages, err := GetMessagesByConversation(conv.ID, 10, 0)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(messages), 3)

	// Messages should be ordered by created_at DESC (newest first)
	foundIDs := make(map[string]bool)
	for _, msg := range messages {
		if msg.ID == msg1.ID || msg.ID == msg2.ID || msg.ID == msg3.ID {
			foundIDs[msg.ID] = true
			assert.NotEmpty(t, msg.SenderUsername)
			assert.Contains(t, []string{testUserJohnDoe, testUserDanielG}, msg.SenderID)
		}
	}

	assert.True(t, foundIDs[msg1.ID])
	assert.True(t, foundIDs[msg2.ID])
	assert.True(t, foundIDs[msg3.ID])

	// Verify ordering (most recent first)
	if len(messages) >= 2 {
		assert.True(t, messages[0].CreatedAt.After(messages[1].CreatedAt) ||
			messages[0].CreatedAt.Equal(messages[1].CreatedAt))
	}

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestUpdateMessage(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation and message
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	msg, err := InsertMessage(conv.ID, testUserJohnDoe, "Original message", "text", nil)
	require.NoError(t, err)

	// Update the message
	newContent := "Edited message"
	err = UpdateMessage(msg.ID, newContent)
	require.NoError(t, err)

	// Verify the update
	messages, err := GetMessagesByConversation(conv.ID, 10, 0)
	require.NoError(t, err)

	var updatedMsg MessageWithSender
	for _, m := range messages {
		if m.ID == msg.ID {
			updatedMsg = m
			break
		}
	}

	assert.Equal(t, newContent, updatedMsg.Content)
	assert.NotNil(t, updatedMsg.EditedAt)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestDeleteMessage(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation and message
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserMariaGarcia)
	require.NoError(t, err)

	msg, err := InsertMessage(conv.ID, testUserJohnDoe, "Message to be deleted", "text", nil)
	require.NoError(t, err)

	// Delete the message
	err = DeleteMessage(msg.ID)
	require.NoError(t, err)

	// Verify the deletion (soft delete)
	messages, err := GetMessagesByConversation(conv.ID, 10, 0)
	require.NoError(t, err)

	var deletedMsg MessageWithSender
	for _, m := range messages {
		if m.ID == msg.ID {
			deletedMsg = m
			break
		}
	}

	assert.True(t, deletedMsg.IsDeleted)
	assert.Equal(t, "[Message deleted]", deletedMsg.Content)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE id = $1`, msg.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestGetUnreadMessageCount(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	// Insert messages from other user
	_, err = InsertMessage(conv.ID, testUserDanielG, "Message 1", "text", nil)
	require.NoError(t, err)

	_, err = InsertMessage(conv.ID, testUserDanielG, "Message 2", "text", nil)
	require.NoError(t, err)

	// Get unread count for JohnDoe
	count, err := GetUnreadMessageCount(conv.ID, testUserJohnDoe)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

func TestUpdateLastReadAt(t *testing.T) {
	_ = SetupDatabase()

	// Create a conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)

	// Insert messages
	_, err = InsertMessage(conv.ID, testUserDanielG, "Unread message 1", "text", nil)
	require.NoError(t, err)

	// Update last read timestamp
	err = UpdateLastReadAt(conv.ID, testUserJohnDoe)
	require.NoError(t, err)

	// Verify last_read_at was updated
	var lastReadAt *time.Time
	err = db.QueryRow(`
		SELECT last_read_at 
		FROM "ConversationParticipant" 
		WHERE conversation_id = $1 AND user_id = $2
	`, conv.ID, testUserJohnDoe).Scan(&lastReadAt)

	require.NoError(t, err)
	assert.NotNil(t, lastReadAt)
	assert.True(t, time.Since(*lastReadAt) < 5*time.Second)

	// Insert a new message after marking as read
	time.Sleep(100 * time.Millisecond)
	_, err = InsertMessage(conv.ID, testUserDanielG, "New message after read", "text", nil)
	require.NoError(t, err)

	// Should have 1 unread message
	count, err := GetUnreadMessageCount(conv.ID, testUserJohnDoe)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}

// ========== User Presence Tests ==========

func TestUpsertUserPresence(t *testing.T) {
	_ = SetupDatabase()

	// Insert new presence
	err := UpsertUserPresence(testUserJohnDoe, "online")
	require.NoError(t, err)

	// Verify insertion
	presence, err := GetUserPresence(testUserJohnDoe)
	require.NoError(t, err)
	assert.Equal(t, testUserJohnDoe, presence.UserID)
	assert.Equal(t, "online", presence.Status)
	assert.True(t, time.Since(presence.LastSeen) < 5*time.Second)
	assert.True(t, time.Since(presence.UpdatedAt) < 5*time.Second)

	// Update existing presence
	time.Sleep(100 * time.Millisecond)
	err = UpsertUserPresence(testUserJohnDoe, "away")
	require.NoError(t, err)

	// Verify update
	updatedPresence, err := GetUserPresence(testUserJohnDoe)
	require.NoError(t, err)
	assert.Equal(t, "away", updatedPresence.Status)
	assert.True(t, updatedPresence.UpdatedAt.After(presence.UpdatedAt))

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "UserPresence" WHERE user_id = $1`, testUserJohnDoe)
}

func TestGetUserPresence(t *testing.T) {
	_ = SetupDatabase()

	// Set presence
	err := UpsertUserPresence(testUserDanielG, "online")
	require.NoError(t, err)

	// Get presence
	presence, err := GetUserPresence(testUserDanielG)

	require.NoError(t, err)
	assert.Equal(t, testUserDanielG, presence.UserID)
	assert.Equal(t, "online", presence.Status)
	assert.NotZero(t, presence.LastSeen)
	assert.NotZero(t, presence.UpdatedAt)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "UserPresence" WHERE user_id = $1`, testUserDanielG)
}

func TestGetMultipleUserPresence(t *testing.T) {
	_ = SetupDatabase()

	// Set presence for multiple users
	err := UpsertUserPresence(testUserJohnDoe, "online")
	require.NoError(t, err)

	err = UpsertUserPresence(testUserDanielG, "away")
	require.NoError(t, err)

	err = UpsertUserPresence(testUserMariaGarcia, "offline")
	require.NoError(t, err)

	// Get multiple presences
	userIDs := []string{testUserJohnDoe, testUserDanielG, testUserMariaGarcia}
	presences, err := GetMultipleUserPresence(userIDs)

	require.NoError(t, err)
	assert.Len(t, presences, 3)

	// Verify all users are present
	foundUsers := make(map[string]string)
	for _, p := range presences {
		foundUsers[p.UserID] = p.Status
	}

	assert.Equal(t, "online", foundUsers[testUserJohnDoe])
	assert.Equal(t, "away", foundUsers[testUserDanielG])
	assert.Equal(t, "offline", foundUsers[testUserMariaGarcia])

	// Cleanup
	for _, userID := range userIDs {
		_, _ = db.Exec(`DELETE FROM "UserPresence" WHERE user_id = $1`, userID)
	}
}

func TestSetAllUsersOffline(t *testing.T) {
	_ = SetupDatabase()

	// Set multiple users online
	users := []string{testUserJohnDoe, testUserDanielG, testUserMariaGarcia}
	for _, userID := range users {
		err := UpsertUserPresence(userID, "online")
		require.NoError(t, err)
	}

	// Set all users offline
	err := SetAllUsersOffline()
	require.NoError(t, err)

	// Verify all users are offline
	presences, err := GetMultipleUserPresence(users)
	require.NoError(t, err)

	for _, p := range presences {
		assert.Equal(t, "offline", p.Status)
	}

	// Cleanup
	for _, userID := range users {
		_, _ = db.Exec(`DELETE FROM "UserPresence" WHERE user_id = $1`, userID)
	}
}

// ========== Complete Chat Workflow Tests ==========

func TestCompleteChatWorkflow(t *testing.T) {
	_ = SetupDatabase()

	// Step 1: Create a direct conversation
	conv, err := CreateDirectConversation(testUserJohnDoe, testUserDanielG)
	require.NoError(t, err)
	assert.Equal(t, "direct", conv.Type)

	// Step 2: Set user presence
	err = UpsertUserPresence(testUserJohnDoe, "online")
	require.NoError(t, err)
	err = UpsertUserPresence(testUserDanielG, "online")
	require.NoError(t, err)

	// Step 3: Send messages
	msg1, err := InsertMessage(conv.ID, testUserJohnDoe, "Hey Daniel!", "text", nil)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	msg2, err := InsertMessage(conv.ID, testUserDanielG, "Hi John! How are you?", "text", nil)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	attachmentURL := "https://example.com/image.jpg"
	msg3, err := InsertMessage(conv.ID, testUserJohnDoe, "Check this out!", "image", &attachmentURL)
	require.NoError(t, err)

	// Step 4: Check unread count for Daniel
	unreadCount, err := GetUnreadMessageCount(conv.ID, testUserDanielG)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, unreadCount, 2) // msg1 and msg3 from John

	// Step 5: Daniel marks as read
	err = UpdateLastReadAt(conv.ID, testUserDanielG)
	require.NoError(t, err)

	// Step 6: Verify unread count is now 0
	unreadCount, err = GetUnreadMessageCount(conv.ID, testUserDanielG)
	require.NoError(t, err)
	assert.Equal(t, 0, unreadCount)

	// Step 7: Get messages with sender info
	messages, err := GetMessagesByConversation(conv.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(messages), 3)

	// Verify message details
	for _, msg := range messages {
		assert.NotEmpty(t, msg.SenderUsername)
		if msg.ID == msg3.ID {
			assert.Equal(t, "image", msg.MessageType)
			assert.NotNil(t, msg.AttachmentURL)
		}
	}

	// Step 8: Edit a message
	err = UpdateMessage(msg1.ID, "Hey Daniel! (edited)")
	require.NoError(t, err)

	// Step 9: Delete a message
	err = DeleteMessage(msg2.ID)
	require.NoError(t, err)

	// Step 10: Get user conversations
	conversations, err := GetUserConversations(testUserJohnDoe, 10, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, conversations)

	foundConv := false
	for _, c := range conversations {
		if c.ID == conv.ID {
			foundConv = true
			assert.Len(t, c.Participants, 2)
		}
	}
	assert.True(t, foundConv)

	// Step 11: Set user offline
	err = UpsertUserPresence(testUserJohnDoe, "offline")
	require.NoError(t, err)

	presence, err := GetUserPresence(testUserJohnDoe)
	require.NoError(t, err)
	assert.Equal(t, "offline", presence.Status)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "UserPresence" WHERE user_id IN ($1, $2)`, testUserJohnDoe, testUserDanielG)
}

func TestGroupChatWorkflow(t *testing.T) {
	_ = SetupDatabase()

	// Step 1: Create a group conversation
	title := "Project Team"
	participants := []string{testUserDanielG, testUserMariaGarcia}
	conv, err := CreateGroupConversation(title, testUserJohnDoe, participants)
	require.NoError(t, err)
	assert.Equal(t, "group", conv.Type)
	assert.Equal(t, title, *conv.Title)

	// Step 2: Verify all participants
	convDetails, err := GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convDetails.Participants, 3)

	// Step 3: Add another participant
	err = AddParticipantToConversation(conv.ID, testUserRobertSmith)
	require.NoError(t, err)

	convDetails, err = GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convDetails.Participants, 4)

	// Step 4: Send group messages
	_, err = InsertMessage(conv.ID, testUserJohnDoe, "Welcome to the team!", "text", nil)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = InsertMessage(conv.ID, testUserDanielG, "Thanks for adding me!", "text", nil)
	require.NoError(t, err)

	// Step 5: Get messages
	messages, err := GetMessagesByConversation(conv.ID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(messages), 2)

	// Step 6: Remove a participant
	err = RemoveParticipantFromConversation(conv.ID, testUserMariaGarcia)
	require.NoError(t, err)

	convDetails, err = GetConversationByID(conv.ID)
	require.NoError(t, err)
	assert.Len(t, convDetails.Participants, 3)
	assert.NotContains(t, convDetails.Participants, testUserMariaGarcia)

	// Cleanup
	_, _ = db.Exec(`DELETE FROM "Message" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "ConversationParticipant" WHERE conversation_id = $1`, conv.ID)
	_, _ = db.Exec(`DELETE FROM "Conversation" WHERE id = $1`, conv.ID)
}
