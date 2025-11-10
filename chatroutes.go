package main

import (
	"bukidlink/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ========== Request/Response Types ==========

// CreateDirectConversationRequest represents a request to create a direct conversation
type CreateDirectConversationRequest struct {
	UserID1 string `json:"user_id_1" binding:"required"`
	UserID2 string `json:"user_id_2" binding:"required"`
}

// CreateGroupConversationRequest represents a request to create a group conversation
type CreateGroupConversationRequest struct {
	Title          string   `json:"title" binding:"required"`
	CreatorID      string   `json:"creator_id" binding:"required"`
	ParticipantIDs []string `json:"participant_ids" binding:"required"`
}

// AddParticipantRequest represents a request to add a participant to a conversation
type AddParticipantRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	SenderID      string  `json:"sender_id" binding:"required"`
	Content       string  `json:"content" binding:"required"`
	MessageType   string  `json:"message_type"` // defaults to "text"
	AttachmentURL *string `json:"attachment_url"`
}

// EditMessageRequest represents a request to edit a message
type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// MarkAsReadRequest represents a request to mark messages as read
type MarkAsReadRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// ========== Conversation Handlers ==========

// postDirectConversationHandler creates a new direct (1-on-1) conversation
func postDirectConversationHandler(c *gin.Context) {
	var req CreateDirectConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Validate that the two users are different
	if req.UserID1 == req.UserID2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create conversation with yourself"})
		return
	}

	// Create or get existing direct conversation
	conv, err := db.CreateDirectConversation(req.UserID1, req.UserID2)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, conv)
}

// postGroupConversationHandler creates a new group conversation
func postGroupConversationHandler(c *gin.Context) {
	var req CreateGroupConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Validate participant count
	if len(req.ParticipantIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group conversation must have at least one participant besides creator"})
		return
	}

	conv, err := db.CreateGroupConversation(req.Title, req.CreatorID, req.ParticipantIDs)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, conv)
}

// getConversationByIDHandler retrieves a conversation by ID with participants
func getConversationByIDHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	conv, err := db.GetConversationByID(conversationID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, conv)
}

// getUserConversationsHandler retrieves all conversations for a user
func getUserConversationsHandler(c *gin.Context) {
	userID := c.Param("user_id")

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	conversations, err := db.GetUserConversations(userID, limit, offset)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, conversations)
}

// postAddParticipantHandler adds a participant to a group conversation
func postAddParticipantHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Verify it's a group conversation (not direct)
	conv, err := db.GetConversationByID(conversationID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	if conv.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add participants to direct conversations"})
		return
	}

	err = db.AddParticipantToConversation(conversationID, req.UserID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant added successfully"})
}

// deleteRemoveParticipantHandler removes a participant from a conversation
func deleteRemoveParticipantHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	userID := c.Param("user_id")

	// Verify it's a group conversation (not direct)
	conv, err := db.GetConversationByID(conversationID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	if conv.Type != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove participants from direct conversations"})
		return
	}

	err = db.RemoveParticipantFromConversation(conversationID, userID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed successfully"})
}

// ========== Message Handlers ==========

// postSendMessageHandler sends a new message to a conversation
func postSendMessageHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Default message type to "text" if not specified
	messageType := req.MessageType
	if messageType == "" {
		messageType = "text"
	}

	// Verify sender is a participant
	conv, err := db.GetConversationByID(conversationID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Check if sender is in participants list
	isParticipant := false
	for _, participantID := range conv.Participants {
		if participantID == req.SenderID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sender is not a participant in this conversation"})
		return
	}

	msg, err := db.InsertMessage(conversationID, req.SenderID, req.Content, messageType, req.AttachmentURL)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// getConversationMessagesHandler retrieves messages from a conversation
func getConversationMessagesHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	messages, err := db.GetMessagesByConversation(conversationID, limit, offset)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// patchEditMessageHandler edits an existing message
func patchEditMessageHandler(c *gin.Context) {
	messageID := c.Param("message_id")

	var req EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.UpdateMessage(messageID, req.Content)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

// deleteMessageHandler soft deletes a message
func deleteMessageHandler(c *gin.Context) {
	messageID := c.Param("message_id")

	err := db.DeleteMessage(messageID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// ========== Read Receipt Handlers ==========

// getUnreadCountHandler gets the unread message count for a user in a conversation
func getUnreadCountHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}

	count, err := db.GetUnreadMessageCount(conversationID, userID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

// postMarkAsReadHandler marks all messages in a conversation as read for a user
func postMarkAsReadHandler(c *gin.Context) {
	conversationID := c.Param("conversation_id")

	var req MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.UpdateLastReadAt(conversationID, req.UserID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

// ========== User Presence Handlers ==========

// postUpdatePresenceHandler updates a user's presence status
func postUpdatePresenceHandler(c *gin.Context) {
	userID := c.Param("user_id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=online offline away"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.UpsertUserPresence(userID, req.Status)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Presence updated successfully"})
}

// getUserPresenceHandler gets a user's presence status
func getUserPresenceHandler(c *gin.Context) {
	userID := c.Param("user_id")

	presence, err := db.GetUserPresence(userID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, presence)
}

// getMultipleUserPresenceHandler gets presence for multiple users
func getMultipleUserPresenceHandler(c *gin.Context) {
	var req struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	presences, err := db.GetMultipleUserPresence(req.UserIDs)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, presences)
}

// ========== Route Registration ==========

// setupChatRoutes registers all chat-related routes
func setupChatRoutes(r *gin.Engine) {
	chat := r.Group("/chat")
	{
		// Conversation endpoints
		chat.POST("/conversation/direct", postDirectConversationHandler)
		chat.POST("/conversation/group", postGroupConversationHandler)
		chat.GET("/conversation/:conversation_id", getConversationByIDHandler)
		chat.GET("/user/:user_id/conversations", getUserConversationsHandler)
		chat.POST("/conversation/:conversation_id/participant", postAddParticipantHandler)
		chat.DELETE("/conversation/:conversation_id/participant/:user_id", deleteRemoveParticipantHandler)

		// Message endpoints
		chat.POST("/conversation/:conversation_id/message", postSendMessageHandler)
		chat.GET("/conversation/:conversation_id/messages", getConversationMessagesHandler)
		chat.PATCH("/message/:message_id", patchEditMessageHandler)
		chat.DELETE("/message/:message_id", deleteMessageHandler)

		// Read receipt endpoints
		chat.GET("/conversation/:conversation_id/unread", getUnreadCountHandler)
		chat.POST("/conversation/:conversation_id/read", postMarkAsReadHandler)

		// User presence endpoints
		chat.POST("/user/:user_id/presence", postUpdatePresenceHandler)
		chat.GET("/user/:user_id/presence", getUserPresenceHandler)
		chat.POST("/presence/batch", getMultipleUserPresenceHandler)
	}
}
