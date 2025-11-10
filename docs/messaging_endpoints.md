# Messaging API Endpoints

Complete REST API and WebSocket documentation for the real-time messaging feature in BukidLink.

## Base URL
- REST API: All messaging endpoints are prefixed with `/chat`
- WebSocket: `ws://localhost:8080/ws?user_id=<uuid>`

---

## WebSocket Connection

### Connect to WebSocket
Establish a WebSocket connection for real-time messaging.

**Endpoint:** `GET /ws?user_id=<uuid>`

**Query Parameters:**
- `user_id` (required) - The UUID of the user connecting

**Connection Example (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?user_id=YOUR_USER_ID');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

**Notes:**
- User presence automatically set to "online" on connect
- User presence automatically set to "offline" on disconnect (if no other connections)
- Supports multiple simultaneous connections per user
- Connection kept alive with ping/pong (54-second intervals)

---

## WebSocket Message Format

All WebSocket messages follow this envelope structure:

```json
{
  "type": "message_type",
  "payload": { /* type-specific data */ }
}
```

### Message Types

#### Client → Server:
- `message` - Send a new message
- `typing` - Start typing indicator
- `stop_typing` - Stop typing indicator
- `presence` - Update presence status
- `read` - Mark messages as read
- `message_edit` - Edit a message
- `message_delete` - Delete a message
- `pong` - Respond to server ping

#### Server → Client:
- `message` - New message received
- `typing` - User started typing
- `stop_typing` - User stopped typing
- `presence` - User presence updated
- `read` - Messages marked as read
- `message_edit` - Message was edited
- `message_delete` - Message was deleted
- `error` - Error occurred
- `ping` - Keep-alive ping

---

## WebSocket Messages

### 1. Send Message

**Client sends:**
```json
{
  "type": "message",
  "payload": {
    "conversation_id": "uuid",
    "sender_id": "uuid",
    "content": "Hello, world!",
    "message_type": "text",
    "attachment_url": null
  }
}
```

**Server broadcasts to participants:**
```json
{
  "type": "message",
  "payload": {
    "conversation_id": "uuid",
    "message_id": "uuid",
    "sender_id": "uuid",
    "content": "Hello, world!",
    "message_type": "text",
    "attachment_url": null,
    "created_at": "2025-11-10T10:15:00Z"
  }
}
```

**Notes:**
- Message saved to database before broadcast
- All conversation participants receive the message
- Sender receives their own message (for multi-device sync)

---

### 2. Typing Indicator

**Client sends:**
```json
{
  "type": "typing",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid"
  }
}
```

**Server broadcasts to other participants:**
```json
{
  "type": "typing",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid",
    "username": null
  }
}
```

**Notes:**
- Sender excluded from broadcast
- Not persisted to database
- Client should send `stop_typing` after 3-5 seconds of inactivity

---

### 3. Stop Typing

**Client sends:**
```json
{
  "type": "stop_typing",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid"
  }
}
```

**Server broadcasts to other participants:**
```json
{
  "type": "stop_typing",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid"
  }
}
```

---

### 4. Presence Update

**Client sends:**
```json
{
  "type": "presence",
  "payload": {
    "user_id": "uuid",
    "status": "away"
  }
}
```

**Server broadcasts to all connected users:**
```json
{
  "type": "presence",
  "payload": {
    "user_id": "uuid",
    "status": "away"
  }
}
```

**Valid status values:**
- `online`
- `offline`
- `away`

**Notes:**
- Broadcast to ALL connected users, not just conversation participants
- Persisted to database
- Automatically set on connect/disconnect

---

### 5. Read Receipt

**Client sends:**
```json
{
  "type": "read",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid"
  }
}
```

**Server broadcasts to other participants:**
```json
{
  "type": "read",
  "payload": {
    "conversation_id": "uuid",
    "user_id": "uuid",
    "timestamp": "2025-11-10T10:20:00Z"
  }
}
```

**Notes:**
- Updates `last_read_at` in database
- Sender excluded from broadcast

---

### 6. Edit Message

**Client sends:**
```json
{
  "type": "message_edit",
  "payload": {
    "conversation_id": "uuid",
    "message_id": "uuid",
    "content": "Updated message content"
  }
}
```

**Server broadcasts to all participants:**
```json
{
  "type": "message_edit",
  "payload": {
    "conversation_id": "uuid",
    "message_id": "uuid",
    "content": "Updated message content"
  }
}
```

**Notes:**
- Updates message in database
- Sets `edited_at` timestamp

---

### 7. Delete Message

**Client sends:**
```json
{
  "type": "message_delete",
  "payload": {
    "conversation_id": "uuid",
    "message_id": "uuid"
  }
}
```

**Server broadcasts to all participants:**
```json
{
  "type": "message_delete",
  "payload": {
    "conversation_id": "uuid",
    "message_id": "uuid"
  }
}
```

**Notes:**
- Soft delete in database
- Content replaced with "[Message deleted]"

---

### 8. Error Message

**Server sends:**
```json
{
  "type": "error",
  "payload": {
    "message": "Error description",
    "code": "ERROR_CODE"
  }
}
```

**Common errors:**
- "Invalid message format"
- "Sender ID does not match authenticated user"
- "Cannot update another user's presence"
- "Failed to save message"

---

## REST API Endpoints

## Base URL
All messaging endpoints are prefixed with `/chat`

---

## Conversation Endpoints

### 1. Create Direct Conversation
Create a 1-on-1 conversation between two users.

**Endpoint:** `POST /chat/conversation/direct`

**Request Body:**
```json
{
  "user_id_1": "uuid",
  "user_id_2": "uuid"
}
```

**Response:** `201 Created`
```json
{
  "id": "conversation-uuid",
  "title": null,
  "type": "direct",
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

**Notes:**
- Returns existing conversation if one already exists between these users
- Users cannot create conversation with themselves (400 error)

---

### 2. Create Group Conversation
Create a group conversation with multiple participants.

**Endpoint:** `POST /chat/conversation/group`

**Request Body:**
```json
{
  "title": "Group Name",
  "creator_id": "uuid",
  "participant_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response:** `201 Created`
```json
{
  "id": "conversation-uuid",
  "title": "Group Name",
  "type": "group",
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

**Notes:**
- Creator is automatically added as a participant
- At least one additional participant required

---

### 3. Get Conversation by ID
Retrieve conversation details with participant list.

**Endpoint:** `GET /chat/conversation/:conversation_id`

**Response:** `200 OK`
```json
{
  "id": "conversation-uuid",
  "title": "Group Name",
  "type": "group",
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:30:00Z",
  "participants": ["uuid1", "uuid2", "uuid3"]
}
```

---

### 4. Get User's Conversations
List all conversations for a user with pagination.

**Endpoint:** `GET /chat/user/:user_id/conversations?limit=20&offset=0`

**Query Parameters:**
- `limit` (optional, default: 20) - Number of conversations to return
- `offset` (optional, default: 0) - Pagination offset

**Response:** `200 OK`
```json
[
  {
    "id": "conversation-uuid",
    "title": "Group Name",
    "type": "group",
    "created_at": "2025-11-10T10:00:00Z",
    "updated_at": "2025-11-10T10:30:00Z",
    "participants": ["uuid1", "uuid2", "uuid3"]
  },
  ...
]
```

**Notes:**
- Ordered by most recent activity (updated_at DESC)

---

### 5. Add Participant to Conversation
Add a user to a group conversation.

**Endpoint:** `POST /chat/conversation/:conversation_id/participant`

**Request Body:**
```json
{
  "user_id": "uuid"
}
```

**Response:** `200 OK`
```json
{
  "message": "Participant added successfully"
}
```

**Notes:**
- Only works for group conversations (direct conversations return 400)

---

### 6. Remove Participant from Conversation
Remove a user from a group conversation.

**Endpoint:** `DELETE /chat/conversation/:conversation_id/participant/:user_id`

**Response:** `200 OK`
```json
{
  "message": "Participant removed successfully"
}
```

**Notes:**
- Only works for group conversations (direct conversations return 400)

---

## Message Endpoints

### 7. Send Message
Send a new message to a conversation.

**Endpoint:** `POST /chat/conversation/:conversation_id/message`

**Request Body:**
```json
{
  "sender_id": "uuid",
  "content": "Hello, world!",
  "message_type": "text",
  "attachment_url": null
}
```

**Response:** `201 Created`
```json
{
  "id": "message-uuid",
  "conversation_id": "conversation-uuid",
  "sender_id": "uuid",
  "content": "Hello, world!",
  "message_type": "text",
  "attachment_url": null,
  "created_at": "2025-11-10T10:15:00Z",
  "edited_at": null,
  "is_deleted": false
}
```

**Message Types:**
- `text` (default)
- `image`
- `file`
- `system`

**Notes:**
- Sender must be a participant in the conversation (403 if not)
- Automatically updates conversation's `updated_at` timestamp via trigger

---

### 8. Get Conversation Messages
Retrieve messages from a conversation with pagination.

**Endpoint:** `GET /chat/conversation/:conversation_id/messages?limit=50&offset=0`

**Query Parameters:**
- `limit` (optional, default: 50) - Number of messages to return
- `offset` (optional, default: 0) - Pagination offset

**Response:** `200 OK`
```json
[
  {
    "id": "message-uuid",
    "conversation_id": "conversation-uuid",
    "sender_id": "uuid",
    "content": "Hello, world!",
    "message_type": "text",
    "attachment_url": null,
    "created_at": "2025-11-10T10:15:00Z",
    "edited_at": null,
    "is_deleted": false,
    "sender_username": "JohnDoe",
    "sender_profile_pic": "https://example.com/profile.jpg"
  },
  ...
]
```

**Notes:**
- Messages ordered by created_at DESC (newest first)
- Includes sender information (username and profile picture)

---

### 9. Edit Message
Edit an existing message's content.

**Endpoint:** `PATCH /chat/message/:message_id`

**Request Body:**
```json
{
  "content": "Updated message content"
}
```

**Response:** `200 OK`
```json
{
  "message": "Message updated successfully"
}
```

**Notes:**
- Sets `edited_at` timestamp
- Cannot edit deleted messages

---

### 10. Delete Message
Soft delete a message (marks as deleted, preserves record).

**Endpoint:** `DELETE /chat/message/:message_id`

**Response:** `200 OK`
```json
{
  "message": "Message deleted successfully"
}
```

**Notes:**
- Soft delete: sets `is_deleted = true` and `content = '[Message deleted]'`
- Message remains in database for audit purposes

---

## Read Receipt Endpoints

### 11. Get Unread Message Count
Get the number of unread messages for a user in a conversation.

**Endpoint:** `GET /chat/conversation/:conversation_id/unread?user_id=uuid`

**Query Parameters:**
- `user_id` (required) - The user to check unread count for

**Response:** `200 OK`
```json
{
  "unread_count": 5
}
```

**Notes:**
- Counts messages sent after user's `last_read_at` timestamp
- Excludes messages sent by the user themselves

---

### 12. Mark Messages as Read
Mark all messages in a conversation as read for a user.

**Endpoint:** `POST /chat/conversation/:conversation_id/read`

**Request Body:**
```json
{
  "user_id": "uuid"
}
```

**Response:** `200 OK`
```json
{
  "message": "Messages marked as read"
}
```

**Notes:**
- Updates `last_read_at` to current timestamp
- Future unread counts will only include newer messages

---

## User Presence Endpoints

### 13. Update User Presence
Set or update a user's online status.

**Endpoint:** `POST /chat/user/:user_id/presence`

**Request Body:**
```json
{
  "status": "online"
}
```

**Valid Status Values:**
- `online`
- `offline`
- `away`

**Response:** `200 OK`
```json
{
  "message": "Presence updated successfully"
}
```

**Notes:**
- Automatically updates `last_seen` and `updated_at` timestamps
- Uses UPSERT pattern (creates or updates existing record)

---

### 14. Get User Presence
Retrieve a single user's presence status.

**Endpoint:** `GET /chat/user/:user_id/presence`

**Response:** `200 OK`
```json
{
  "user_id": "uuid",
  "status": "online",
  "last_seen": "2025-11-10T10:30:00Z",
  "updated_at": "2025-11-10T10:30:00Z"
}
```

---

### 15. Get Multiple User Presence
Retrieve presence status for multiple users at once.

**Endpoint:** `POST /chat/presence/batch`

**Request Body:**
```json
{
  "user_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**Response:** `200 OK`
```json
[
  {
    "user_id": "uuid1",
    "status": "online",
    "last_seen": "2025-11-10T10:30:00Z",
    "updated_at": "2025-11-10T10:30:00Z"
  },
  {
    "user_id": "uuid2",
    "status": "away",
    "last_seen": "2025-11-10T10:25:00Z",
    "updated_at": "2025-11-10T10:25:00Z"
  },
  ...
]
```

**Notes:**
- Useful for displaying online status in conversation lists
- Only returns presence for users that have records (may return fewer results than requested)

---

## Error Responses

All endpoints follow standard HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request body or parameters
- `403 Forbidden` - User not authorized for action
- `500 Internal Server Error` - Database or server error

Error response format:
```json
{
  "error": "Error message description"
}
```

---

## Database Triggers

The following PostgreSQL triggers automatically maintain data consistency:

1. **update_conversation_timestamp** - Updates conversation's `updated_at` when a new message is sent
2. **auto_create_user_balance** - Creates wallet balance when new user is created

---

## Testing

Run all messaging endpoint tests:
```bash
go test -v -run "Test.*Message|Test.*Conversation|Test.*Presence"
```

Run specific test category:
```bash
# Message tests only
go test -v -run "Test.*Message"

# Read receipt tests
go test -v -run "TestGetUnreadCount|TestMarkAsRead"

# Presence tests
go test -v -run "Test.*Presence"
```

---

## WebSocket Implementation Details

### Architecture

**Hub Pattern:**
- Central `Hub` manages all WebSocket connections
- Clients registered/unregistered on connect/disconnect
- Broadcast channel for distributing messages to relevant clients

**Client Structure:**
- Each WebSocket connection represented by a `Client` struct
- Buffered send channel (256 messages) prevents blocking
- Separate goroutines for reading and writing (ReadPump/WritePump)

**Concurrency:**
- Thread-safe operations using `sync.RWMutex`
- Multiple connections per user supported
- Non-blocking sends with fallback disconnect

### Connection Lifecycle

1. **Connect:**
   - HTTP upgraded to WebSocket
   - Client registered in hub
   - Presence set to "online"
   - Presence broadcast to all users

2. **Active:**
   - Ping/pong keep-alive (54s interval)
   - Read timeout: 60s
   - Write timeout: 10s
   - Handles reconnection for same user

3. **Disconnect:**
   - Client unregistered from hub
   - If last connection for user, set presence to "offline"
   - Broadcast offline status

### Message Broadcasting

**Conversation-scoped:**
```go
hub.Broadcast <- &BroadcastMessage{
    ConversationID: "uuid",
    ExcludeUserID:  "sender-uuid",  // Optional
    Data:           messageBytes,
}
```

**User-scoped:**
```go
hub.Broadcast <- &BroadcastMessage{
    UserIDs: []string{"uuid1", "uuid2"},
    Data:    messageBytes,
}
```

**Global (all users):**
```go
// Used for presence updates
hub.broadcastPresence(userID, "online")
```

### Testing

**Run WebSocket tests:**
```bash
go test -v -run "TestWebSocket"
```

**Test HTML client:**
1. Start server: `go run .`
2. Open `websocket_client.html` in browser
3. Enter user ID and connect
4. Test messaging, typing indicators, presence

**Manual testing with wscat:**
```bash
npm install -g wscat
wscat -c "ws://localhost:8080/ws?user_id=YOUR_USER_ID"

# Send a message
{"type":"message","payload":{"conversation_id":"CONV_ID","sender_id":"USER_ID","content":"Hello"}}
```

---

## Security Considerations

### Current Implementation (Development)
- WebSocket accepts connections from any origin
- User ID passed as query parameter (not secure)
- No authentication/authorization checks

### Production Requirements

1. **Authentication:**
   ```go
   // Use JWT token instead of user_id query param
   token := c.Query("token")
   userID, err := validateJWT(token)
   ```

2. **Origin Validation:**
   ```go
   var upgrader = websocket.Upgrader{
       CheckOrigin: func(r *http.Request) bool {
           origin := r.Header.Get("Origin")
           return origin == "https://yourdomain.com"
       },
   }
   ```

3. **Message Validation:**
   - Verify user is participant before sending to conversation
   - Rate limiting to prevent spam
   - Content validation and sanitization

4. **TLS/WSS:**
   - Use `wss://` instead of `ws://` in production
   - Configure HTTPS on server

---

## Performance Optimization

### Current Configuration
- Send buffer: 256 messages per client
- Broadcast channel: 256 messages
- Ping interval: 54 seconds
- Read timeout: 60 seconds
- Write timeout: 10 seconds

### Scaling Considerations

**Horizontal Scaling:**
- Use Redis Pub/Sub for cross-server message distribution
- Shared session store for user presence
- Load balancer with sticky sessions

**Message Queuing:**
- Add message queue (RabbitMQ/Redis) for offline delivery
- Store undelivered messages for offline users
- Batch delivery on reconnect

**Database Optimization:**
- Index on `conversation_id` and `sender_id`
- Partition large message tables by date
- Cache active conversation participants

---

## Next Steps for Real-Time Messaging

To complete the real-time messaging feature, implement:

1. **~~WebSocket Layer~~** ✅ **COMPLETED**
   - ~~Socket connection handler~~
   - ~~Real-time message delivery~~
   - ~~Presence broadcasting~~
   - ~~Typing indicators~~

2. **File Upload**
   - Image upload endpoint
   - File storage (S3 or local)
   - Update `attachment_url` field

3. **Authentication/Authorization**
   - JWT middleware
   - Verify user permissions
   - Secure conversation access

4. **Push Notifications**
   - FCM integration for mobile
   - Email notifications for offline users
   - Notification preferences

---

## Curl Examples

### Send a message:
```bash
curl -X POST http://localhost:8080/chat/conversation/CONV_ID/message \
  -H "Content-Type: application/json" \
  -d '{
    "sender_id": "USER_ID",
    "content": "Hello!",
    "message_type": "text"
  }'
```

### Get messages:
```bash
curl http://localhost:8080/chat/conversation/CONV_ID/messages?limit=20&offset=0
```

### Mark as read:
```bash
curl -X POST http://localhost:8080/chat/conversation/CONV_ID/read \
  -H "Content-Type: application/json" \
  -d '{"user_id": "USER_ID"}'
```

### Update presence:
```bash
curl -X POST http://localhost:8080/chat/user/USER_ID/presence \
  -H "Content-Type: application/json" \
  -d '{"status": "online"}'
```
