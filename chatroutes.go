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
	}
}
