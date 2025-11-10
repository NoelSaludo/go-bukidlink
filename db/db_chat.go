package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Conversation represents a chat conversation (direct or group)
type Conversation struct {
	ID        string    `json:"id"`
	Title     *string   `json:"title"` // NULL for direct conversations
	Type      string    `json:"type"`  // "direct" or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` // Last message timestamp
}

// ConversationParticipant represents a user's participation in a conversation
type ConversationParticipant struct {
	ConversationID string     `json:"conversation_id"`
	UserID         string     `json:"user_id"`
	JoinedAt       time.Time  `json:"joined_at"`
	LastReadAt     *time.Time `json:"last_read_at"`
	IsMuted        bool       `json:"is_muted"`
}

// Message represents a chat message
type Message struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Content        string     `json:"content"`
	MessageType    string     `json:"message_type"` // "text", "image", "file", "system"
	AttachmentURL  *string    `json:"attachment_url"`
	CreatedAt      time.Time  `json:"created_at"`
	EditedAt       *time.Time `json:"edited_at"`
	IsDeleted      bool       `json:"is_deleted"`
}

// UserPresence represents a user's online status
type UserPresence struct {
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"` // "online", "offline", "away"
	LastSeen  time.Time `json:"last_seen"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConversationWithParticipants includes participant details
type ConversationWithParticipants struct {
	Conversation
	Participants []string `json:"participants"` // Array of user IDs
}

// MessageWithSender includes sender information
type MessageWithSender struct {
	Message
	SenderUsername   string  `json:"sender_username"`
	SenderProfilePic *string `json:"sender_profile_pic"`
}

// ========== Conversation Functions ==========

// CreateDirectConversation creates a new direct (1-on-1) conversation
func CreateDirectConversation(userID1, userID2 string) (Conversation, error) {
	var conv Conversation

	// Check if conversation already exists between these two users
	existingConv, err := GetDirectConversation(userID1, userID2)
	if err == nil {
		return existingConv, nil // Return existing conversation
	}

	tx, err := db.Begin()
	if err != nil {
		return conv, err
	}

	conversationID := uuid.New().String()

	// Insert conversation
	err = tx.QueryRow(`
		INSERT INTO "Conversation" (id, title, type, created_at, updated_at)
		VALUES ($1, NULL, 'direct', NOW(), NOW())
		RETURNING id, title, type, created_at, updated_at
	`, conversationID).Scan(&conv.ID, &conv.Title, &conv.Type, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return conv, rollbackAndReturn(tx, err)
	}

	// Add both participants
	_, err = tx.Exec(`
		INSERT INTO "ConversationParticipant" (conversation_id, user_id, joined_at)
		VALUES ($1, $2, NOW()), ($1, $3, NOW())
	`, conversationID, userID1, userID2)
	if err != nil {
		return conv, rollbackAndReturn(tx, err)
	}

	return conv, tx.Commit()
}

// CreateGroupConversation creates a new group conversation
func CreateGroupConversation(title string, creatorID string, participantIDs []string) (Conversation, error) {
	var conv Conversation

	tx, err := db.Begin()
	if err != nil {
		return conv, err
	}

	conversationID := uuid.New().String()

	// Insert conversation
	err = tx.QueryRow(`
		INSERT INTO "Conversation" (id, title, type, created_at, updated_at)
		VALUES ($1, $2, 'group', NOW(), NOW())
		RETURNING id, title, type, created_at, updated_at
	`, conversationID, title).Scan(&conv.ID, &conv.Title, &conv.Type, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return conv, rollbackAndReturn(tx, err)
	}

	// Add creator as first participant
	allParticipants := append([]string{creatorID}, participantIDs...)

	// Add all participants
	for _, userID := range allParticipants {
		_, err = tx.Exec(`
			INSERT INTO "ConversationParticipant" (conversation_id, user_id, joined_at)
			VALUES ($1, $2, NOW())
		`, conversationID, userID)
		if err != nil {
			return conv, rollbackAndReturn(tx, err)
		}
	}

	return conv, tx.Commit()
}

// GetDirectConversation finds existing direct conversation between two users
func GetDirectConversation(userID1, userID2 string) (Conversation, error) {
	var conv Conversation

	err := db.QueryRow(`
		SELECT c.id, c.title, c.type, c.created_at, c.updated_at
		FROM "Conversation" c
		WHERE c.type = 'direct'
		AND c.id IN (
			SELECT cp1.conversation_id
			FROM "ConversationParticipant" cp1
			INNER JOIN "ConversationParticipant" cp2 
				ON cp1.conversation_id = cp2.conversation_id
			WHERE cp1.user_id = $1 AND cp2.user_id = $2
		)
		LIMIT 1
	`, userID1, userID2).Scan(&conv.ID, &conv.Title, &conv.Type, &conv.CreatedAt, &conv.UpdatedAt)

	return conv, err
}

// GetConversationByID retrieves a conversation by ID
func GetConversationByID(conversationID string) (ConversationWithParticipants, error) {
	var conv ConversationWithParticipants

	// Get conversation details
	err := db.QueryRow(`
		SELECT id, title, type, created_at, updated_at
		FROM "Conversation"
		WHERE id = $1
	`, conversationID).Scan(&conv.ID, &conv.Title, &conv.Type, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return conv, err
	}

	// Get participants
	rows, err := db.Query(`
		SELECT user_id
		FROM "ConversationParticipant"
		WHERE conversation_id = $1
		ORDER BY joined_at
	`, conversationID)
	if err != nil {
		return conv, err
	}
	defer rows.Close()

	conv.Participants = make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return conv, err
		}
		conv.Participants = append(conv.Participants, userID)
	}

	return conv, rows.Err()
}

// GetUserConversations retrieves all conversations for a user
func GetUserConversations(userID string, limit, offset int) ([]ConversationWithParticipants, error) {
	conversations := make([]ConversationWithParticipants, 0)

	rows, err := db.Query(`
		SELECT c.id, c.title, c.type, c.created_at, c.updated_at
		FROM "Conversation" c
		INNER JOIN "ConversationParticipant" cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return conversations, err
	}
	defer rows.Close()

	for rows.Next() {
		var conv ConversationWithParticipants
		if err := rows.Scan(&conv.ID, &conv.Title, &conv.Type, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			return conversations, err
		}

		// Get participants for each conversation
		participantRows, err := db.Query(`
			SELECT user_id
			FROM "ConversationParticipant"
			WHERE conversation_id = $1
			ORDER BY joined_at
		`, conv.ID)
		if err != nil {
			return conversations, err
		}

		conv.Participants = make([]string, 0)
		for participantRows.Next() {
			var participantID string
			if err := participantRows.Scan(&participantID); err != nil {
				participantRows.Close()
				return conversations, err
			}
			conv.Participants = append(conv.Participants, participantID)
		}
		participantRows.Close()

		conversations = append(conversations, conv)
	}

	return conversations, rows.Err()
}

// AddParticipantToConversation adds a user to a group conversation
func AddParticipantToConversation(conversationID, userID string) error {
	_, err := db.Exec(`
		INSERT INTO "ConversationParticipant" (conversation_id, user_id, joined_at)
		VALUES ($1, $2, NOW())
	`, conversationID, userID)
	return err
}

// RemoveParticipantFromConversation removes a user from a conversation
func RemoveParticipantFromConversation(conversationID, userID string) error {
	_, err := db.Exec(`
		DELETE FROM "ConversationParticipant"
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID)
	return err
}

// ========== Message Functions ==========

// InsertMessage creates a new message in a conversation
func InsertMessage(conversationID, senderID, content, messageType string, attachmentURL *string) (Message, error) {
	var msg Message

	messageID := uuid.New().String()

	err := db.QueryRow(`
		INSERT INTO "Message" (id, conversation_id, sender_id, content, message_type, attachment_url, created_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), FALSE)
		RETURNING id, conversation_id, sender_id, content, message_type, attachment_url, created_at, edited_at, is_deleted
	`, messageID, conversationID, senderID, content, messageType, attachmentURL).Scan(
		&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content,
		&msg.MessageType, &msg.AttachmentURL, &msg.CreatedAt, &msg.EditedAt, &msg.IsDeleted,
	)

	return msg, err
}

// GetMessagesByConversation retrieves messages for a conversation with pagination
func GetMessagesByConversation(conversationID string, limit, offset int) ([]MessageWithSender, error) {
	messages := make([]MessageWithSender, 0)

	rows, err := db.Query(`
		SELECT 
			m.id, m.conversation_id, m.sender_id, m.content, m.message_type, 
			m.attachment_url, m.created_at, m.edited_at, m.is_deleted,
			u.username, u.profile_pic_url
		FROM "Message" m
		INNER JOIN "User" u ON m.sender_id = u.id
		LEFT JOIN "UserUserDetail" uud ON u.id = uud.user_id
		LEFT JOIN "UserDetail" ud ON uud.detail_id = ud.id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg MessageWithSender
		if err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &msg.MessageType,
			&msg.AttachmentURL, &msg.CreatedAt, &msg.EditedAt, &msg.IsDeleted,
			&msg.SenderUsername, &msg.SenderProfilePic,
		); err != nil {
			return messages, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// UpdateMessage edits an existing message
func UpdateMessage(messageID, newContent string) error {
	_, err := db.Exec(`
		UPDATE "Message"
		SET content = $1, edited_at = NOW()
		WHERE id = $2 AND is_deleted = FALSE
	`, newContent, messageID)
	return err
}

// DeleteMessage soft deletes a message
func DeleteMessage(messageID string) error {
	_, err := db.Exec(`
		UPDATE "Message"
		SET is_deleted = TRUE, content = '[Message deleted]'
		WHERE id = $1
	`, messageID)
	return err
}

// GetUnreadMessageCount counts unread messages for a user in a conversation
func GetUnreadMessageCount(conversationID, userID string) (int, error) {
	var count int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM "Message" m
		WHERE m.conversation_id = $1
		AND m.sender_id != $2
		AND m.created_at > COALESCE(
			(SELECT last_read_at FROM "ConversationParticipant" 
			 WHERE conversation_id = $1 AND user_id = $2),
			'1970-01-01'::TIMESTAMP
		)
	`, conversationID, userID).Scan(&count)

	return count, err
}

// UpdateLastReadAt updates the last read timestamp for a user in a conversation
func UpdateLastReadAt(conversationID, userID string) error {
	_, err := db.Exec(`
		UPDATE "ConversationParticipant"
		SET last_read_at = NOW()
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID)
	return err
}

// ========== User Presence Functions ==========

// UpsertUserPresence creates or updates user presence status
func UpsertUserPresence(userID, status string) error {
	_, err := db.Exec(`
		INSERT INTO "UserPresence" (user_id, status, last_seen, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			last_seen = NOW(),
			updated_at = NOW()
	`, userID, status)
	return err
}

// GetUserPresence retrieves presence status for a user
func GetUserPresence(userID string) (UserPresence, error) {
	var presence UserPresence

	err := db.QueryRow(`
		SELECT user_id, status, last_seen, updated_at
		FROM "UserPresence"
		WHERE user_id = $1
	`, userID).Scan(&presence.UserID, &presence.Status, &presence.LastSeen, &presence.UpdatedAt)

	return presence, err
}

// GetMultipleUserPresence retrieves presence for multiple users
func GetMultipleUserPresence(userIDs []string) ([]UserPresence, error) {
	presences := make([]UserPresence, 0)

	if len(userIDs) == 0 {
		return presences, nil
	}

	query := `
		SELECT user_id, status, last_seen, updated_at
		FROM "UserPresence"
		WHERE user_id = ANY($1)
	`

	rows, err := db.Query(query, pq.Array(userIDs))
	if err != nil {
		return presences, err
	}
	defer rows.Close()

	for rows.Next() {
		var presence UserPresence
		if err := rows.Scan(&presence.UserID, &presence.Status, &presence.LastSeen, &presence.UpdatedAt); err != nil {
			return presences, err
		}
		presences = append(presences, presence)
	}

	return presences, rows.Err()
}

// SetAllUsersOffline sets all users to offline status (useful for server restart)
func SetAllUsersOffline() error {
	_, err := db.Exec(`
		UPDATE "UserPresence"
		SET status = 'offline', last_seen = NOW(), updated_at = NOW()
		WHERE status != 'offline'
	`)
	return err
}
