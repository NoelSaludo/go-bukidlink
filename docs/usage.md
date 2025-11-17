# Usage — go-bukidlink

This document describes the HTTP routes implemented by the server (from `main.go`) and shows example `curl` commands (based on `main_test.go`). It explains expected inputs and outputs for each route.

## Quick start

Start the server (project root):

```bash
# from project root
go run main.go
```

The server binds to `localhost:8080` by default (see `main.go`).

## Dart usage (Dio)

If you're building a Dart or Flutter client, a ready-to-use example using the `dio` HTTP client is available in `docs/usage_dart.md`. It contains model classes, an `ApiClient` wrapper around `Dio`, and example calls for each route implemented by the server.

Quick link:

```
docs/usage_dart.md
```

## Server behavior

- The server initializes the database with `db.SetupDatabase()` during startup.
- Routes are implemented in `main.go`, `itemroutes.go`, and `userroutes.go`. This document lists the currently implemented endpoints and how to use them.

---

## Routes

### GET /ping
**Purpose**: Health check endpoint to verify server is running.

**Response**:
- Status: `200 OK`
- Body: `{"message":"pong"}`

**Example**:
```bash
curl http://localhost:8080/ping
```

---

## Item Routes

### GET /item/:block
**Purpose**: Return a paginated list of items (100 items per block).

**Parameters**:
- `:block` (path) - Page number (0-indexed). Block 0 returns items 0-99, block 1 returns items 100-199, etc.

**Response**: Array of `Item` objects with status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/item/0
```

**Response Example**:
```json
[
    {
        "id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
        "name": "Banana",
        "description": "Ripe Cavendish bananas, sweet and soft",
        "amount": 120,
        "costPKilo": 0.8,
        "category": "fruits",
        "rating": 0
    }
]
```

### GET /item/category/:category
**Purpose**: Return all items in a specific category.

**Parameters**:
- `:category` (path) - Category name. Valid values: `fruits`, `vegetables`, `grains`, `livestock`, `dairy`, `others`

**Response**: Array of `Item` objects with status `200 OK`, or `500 Internal Server Error` for invalid categories.

**Example**:
```bash
curl http://localhost:8080/item/category/fruits
```

**Response Example**:
```json
[
    {
        "id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
        "name": "Banana",
        "description": "Ripe Cavendish bananas, sweet and soft",
        "amount": 120,
        "costPKilo": 0.8,
        "category": "fruits",
        "rating": 0
    }
]
```

### GET /item
**Purpose**: Retrieve a single item by its ID with aggregated rating from reviews.

**Query Parameters**:
- `id` - Item UUID

**Response**: Single `Item` object with status `200 OK`, includes average rating from reviews.

**Example**:
```bash
curl "http://localhost:8080/item?id=a3e1b9f2-7d94-4d3a-9b4a-111111111111"
```

**Response Example**:
```json
{
    "id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
    "name": "Banana",
    "description": "Ripe Cavendish bananas, sweet and soft",
    "amount": 120,
    "costPKilo": 0.8,
    "category": "fruits",
    "rating": 0
}
```

---

## User Routes

### GET /user/:username
**Purpose**: Retrieve a user by username, including their details.

**Parameters**:
- `:username` (path) - Username string

**Response**: Single `User` object with nested `UserDetail` with status `200 OK`, or `400 Bad Request` if user not found.

**Example**:
```bash
curl http://localhost:8080/user/JohnDoe
```

**Response Example**:
```json
{
    "id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "email": "JohnDoe@example.com",
    "username": "JohnDoe",
    "password": "P@ssw0rd",
    "details": {
        "address": "So. Pinamuntasan, Brgy. Aga, Nasugbu, Batangas",
        "first_name": "John",
        "last_name": "Doe",
        "contact_number": "+639123456789",
        "created_date": "2023-10-27T10:00:00Z"
    }
}
```

### POST /user
**Purpose**: Create a new user with their details in a single transaction.

**Request Body**: JSON `User` object with optional nested `UserDetail`.

**Response**: 
- `200 OK` with `{"message":"Success"}` on successful creation
- `409 Conflict` if user already exists
- `400 Bad Request` if JSON is malformed

**Example**:
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "NewUser",
    "password": "securepass123",
    "email": "newuser@example.com",
    "details": {
      "first_name": "John",
      "last_name": "Smith",
      "address": "123 Main St",
      "contact_number": "+1234567890"
    }
  }'
```

**Response**: `{"message":"Success"}`

---

## Order Routes

### GET /order/:user_id
**Purpose**: Get all orders for a specific user, including all order items.

**Parameters**:
- `:user_id` (path) - User UUID

**Response**: JSON array of `Order` objects with nested `OrderItem` arrays, status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/order/d30869ec-fb97-46d8-85a3-82608c01f803
```

**Response Example**:
```json
[
    {
        "id": "11111111-1111-1111-1111-111111111111",
        "userid": "d30869ec-fb97-46d8-85a3-82608c01f803",
        "status": "Packaging",
        "order_date": "2023-10-27T17:07:33.621Z",
        "total_price": 4,
        "items": [
            {
                "id": "21111111-1111-1111-1111-111111111111",
                "order_id": "11111111-1111-1111-1111-111111111111",
                "item_id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
                "quantity": 2,
                "price_at_purchase": 0.8
            },
            {
                "id": "21111111-1111-1111-1111-222222222222",
                "order_id": "11111111-1111-1111-1111-111111111111",
                "item_id": "b7f2c6d4-1aeb-4f5b-9c2b-222222222222",
                "quantity": 2,
                "price_at_purchase": 1.2
            }
        ]
    }
]
```

### POST /order
**Purpose**: Create a new order with multiple items in a single transaction.

**Request Body**: JSON `Order` object with nested `OrderItem` array. The server generates UUIDs for the order and items.

**Response**: `201 Created` with `{"status":"order created","order_id":"<NEW_ORDER_UUID>"}`, or `400 Bad Request`/`500 Internal Server Error` on failure.

**Example**:
```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "userid": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "status": "Packaging",
    "order_date": "2024-10-30T10:00:00Z",
    "total_price": 5.0,
    "items": [
      {
        "item_id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
        "quantity": 2,
        "price_at_purchase": 0.8
      },
      {
        "item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
        "quantity": 1,
        "price_at_purchase": 3.0
      }
    ]
  }'
```

**Response**: `{"status":"order created","order_id":"<NEW_ORDER_UUID>"}`

### PATCH /order/status
**Purpose**: Update the status of an existing order.

**Query Parameters**:
- `id` - Order UUID
- `status` - New status string (e.g., "Packaging", "Shipping", "Delivered")

**Response**: `200 OK` with `{"status":"order status updated"}`, or `500 Internal Server Error` on failure.

**Example**:
```bash
curl -X PATCH "http://localhost:8080/order/status?id=11111111-1111-1111-1111-111111111111&status=Shipping"
```

**Response**: `{"status":"order status updated"}`

### DELETE /order
**Purpose**: Delete an order and all its associated items in a transaction.

**Query Parameters**:
- `order_id` - Order UUID to delete

**Response**: `200 OK` with `{"status":"order deleted"}`, or `500 Internal Server Error` on failure.

**Example**:
```bash
curl -X DELETE "http://localhost:8080/order?order_id=11111111-1111-1111-1111-111111111111"
```

**Response**: `{"status":"order deleted"}`

---

## Cart Routes

### GET /cart/:user_id
**Purpose**: Get a user's cart with all items. Automatically creates a new empty cart if one doesn't exist.

**Parameters**:
- `:user_id` (path) - User UUID

**Response**: JSON `Cart` object with nested `CartItem` array, status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/cart/c6554794-849f-4338-87c5-6db2e2f76514
```

**Response Example**:
```json
{
    "id": "31111111-1111-1111-1111-111111111111",
    "userid": "c6554794-849f-4338-87c5-6db2e2f76514",
    "grand_total": 0,
    "created_at": "2023-10-27T17:10:00.1Z",
    "items": [
        {
            "id": "41111111-1111-1111-1111-111111111111",
            "cart_id": "31111111-1111-1111-1111-111111111111",
            "item_id": "c9d3e8a1-55b2-4f66-a123-333333333333",
            "quantity": 5
        }
    ]
}
```

### POST /cart/item
**Purpose**: Add an item to a cart or update quantity if item already exists.

**Request Body**: JSON object with `cart_id`, `item_id`, and `quantity`.

**Response**: `200 OK` with `{"status":"item added to cart"}`, or `400 Bad Request`/`500 Internal Server Error` on failure.

**Example**:
```bash
curl -X POST http://localhost:8080/cart/item \
  -H "Content-Type: application/json" \
  -d '{
    "cart_id": "31111111-1111-1111-1111-111111111111",
    "item_id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
    "quantity": 3
  }'
```

**Response**: `{"status":"item added to cart"}`

### PATCH /cart/item
**Purpose**: Update the quantity of an existing cart item. If the new `quantity` is less than or equal to 0, the cart item is removed.

**Request Body**: JSON object matching `UpdateCartItemRequest`:
- `cart_item_id` (string) — CartItem UUID (not the Item ID)
- `quantity` (int) — New quantity. If <= 0, the item is deleted.

**Response**: `200 OK` with `{"status":"quantity updated"}`, `400 Bad Request` for malformed JSON or missing `cart_item_id`, or `500 Internal Server Error` on DB errors.

**Example**:
```bash
curl -X PATCH http://localhost:8080/cart/item \
  -H "Content-Type: application/json" \
  -d '{
    "cart_item_id": "41111111-1111-1111-1111-111111111111",
    "quantity": 2
  }'
```

**Note**: The server validates `cart_item_id` is present; see `cartroutes.go` -> `UpdateCartItemRequest`. The DB function `UpdateCartItemQuantity` will remove the cart item when `quantity <= 0`.

### DELETE /cart/item/:cart_item_id
**Purpose**: Remove a specific item from a cart.

**Parameters**:
- `:cart_item_id` (path) - CartItem UUID (not the Item ID)

**Response**: `200 OK` with `{"status":"item removed from cart"}`, or `500 Internal Server Error` on failure.

**Example**:
```bash
curl -X DELETE http://localhost:8080/cart/item/41111111-1111-1111-1111-111111111111
```

**Response**: `{"status":"item removed from cart"}`

---

## Review Routes

### GET /review/:itemId
**Purpose**: Get all reviews/comments for a specific item.

**Parameters**:
- `:itemId` (path) - Item UUID

**Response**: JSON array of `Review` objects with status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/review/a3e1b9f2-7d94-4d3a-9b4a-111111111111
```

**Response Example**:
```json
[
    {
        "id": "894169e9-c907-4d89-84c4-3f1542488c9a",
        "itemid": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
        "userid": "6a24dd2b-d441-4b39-ab85-8fa2bd61065e",
        "content": "Banana very fluffy and has a premium after taste",
        "rating": 4.9
    }
]
```

---

## Notes

### Error Responses
All endpoints may return error responses with appropriate HTTP status codes:
- `400 Bad Request` - Malformed JSON or missing required parameters
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate user)
- `500 Internal Server Error` - Database errors or server issues

Error responses follow this format:
```json
{"error": "error description"}
```

### UUIDs
All entity IDs use UUID v4 format. When creating resources, the server generates UUIDs automatically - you don't need to provide them in POST requests (except for references to existing resources like `user_id` or `item_id`).

---

## Post Routes

### GET /userpost/:user_id
**Purpose**: Get all posts by a specific user (farmer), with images encoded as base64.

**Parameters**:
- `:user_id` (path) - User UUID

**Response**: JSON array of enriched post objects with base64 encoded images, status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/userpost/d30869ec-fb97-46d8-85a3-82608c01f803
```

**Response Example**:
```json
[
    {
        "id": "11111111-1111-1111-1111-111111111111",
        "farmer_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
        "farm_id": "11111111-aaaa-aaaa-aaaa-111111111111",
        "content": "Exploring new farming techniques at Sunny Fields.",
        "image_url": "resources/images/11111111-1111-1111-1111-111111111111_post.png",
        "created_at": "2025-10-30T14:51:36Z",
        "comments": [
            {
                "id": "c1111111-1111-1111-1111-111111111111",
                "post_id": "11111111-1111-1111-1111-111111111111",
                "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
                "content": "Great farming technique!",
                "created_at": "2025-10-30T15:00:00Z"
            }
        ],
        "image_base64": "iVBORw0KGgoAAAANSUhEUgAA...",
        "image_content_type": "image/png"
    }
]
```

**Note**: The `comments` array contains all comments for each post, ordered by creation time (oldest first). If a post has no comments, the array will be empty `[]`.

### GET /userpost/post/:post_id
**Purpose**: Get a specific post by ID, with image encoded as base64.

**Parameters**:
- `:post_id` (path) - Post UUID

**Response**: JSON object with post data and base64 encoded image, status `200 OK`.

**Example**:
```bash
curl http://localhost:8080/userpost/post/11111111-1111-1111-1111-111111111111
```

**Response Example**:
```json
{
    "id": "11111111-1111-1111-1111-111111111111",
    "farmer_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "farm_id": "11111111-aaaa-aaaa-aaaa-111111111111",
    "content": "Exploring new farming techniques at Sunny Fields.",
    "image_url": "resources/images/11111111-1111-1111-1111-111111111111_post.png",
    "created_at": "2025-10-30T14:51:36Z",
    "comments": [
        {
            "id": "c1111111-1111-1111-1111-111111111111",
            "post_id": "11111111-1111-1111-1111-111111111111",
            "user_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
            "content": "Great farming technique!",
            "created_at": "2025-10-30T15:00:00Z"
        },
        {
            "id": "c2222222-2222-2222-2222-222222222222",
            "post_id": "11111111-1111-1111-1111-111111111111",
            "user_id": "c6554794-849f-4338-87c5-6db2e2f76514",
            "content": "Thanks for sharing this!",
            "created_at": "2025-10-30T16:30:00Z"
        }
    ],
    "image_base64": "iVBORw0KGgoAAAANSUhEUgAA...",
    "image_content_type": "image/png"
}
```

**Note**: The `comments` array contains all comments for the post, ordered chronologically (oldest first). Empty array `[]` if no comments exist.

### POST /userpost
**Purpose**: Create a new post with optional image upload.

**Request Body**: JSON object with `post` data and optional `post_image` with base64 encoded image.

**Response**: 
- `201 Created` with `{"post_id":"<NEW_POST_UUID>","message":"Post created successfully"}` on success
- `400 Bad Request` if JSON is malformed
- `500 Internal Server Error` on database errors

**Example**:
```bash
curl -X POST http://localhost:8080/userpost \
  -H "Content-Type: application/json" \
  -d '{
    "post": {
      "farmer_id": "d30869ec-fb97-46d8-85a3-82608c01f803",
      "farm_id": "11111111-aaaa-aaaa-aaaa-111111111111",
      "content": "Check out our new harvest!"
    },
    "post_image": {
      "base64": "iVBORw0KGgoAAAANSUhEUgAA...",
      "content_type": "image/png"
    }
  }'
```

**Response**: `{"post_id":"<NEW_POST_UUID>","message":"Post created successfully"}`

**Note**: The `post_image` field is optional. If omitted, the post will be created without an image.

### PATCH /userpost/:post_id
**Purpose**: Update an existing post's content and/or image.

**Parameters**:
- `:post_id` (path) - Post UUID

**Request Body**: JSON object with optional `content` (string) and/or `post_image` (base64 encoded).

**Response**: 
- `200 OK` with `{"message":"Post updated successfully"}` on success
- `400 Bad Request` if JSON is malformed
- `500 Internal Server Error` on database errors

**Example - Update content only**:
```bash
curl -X PATCH http://localhost:8080/userpost/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated post content"
  }'
```

**Example - Update image only**:
```bash
curl -X PATCH http://localhost:8080/userpost/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -d '{
    "post_image": {
      "base64": "iVBORw0KGgoAAAANSUhEUgAA...",
      "content_type": "image/png"
    }
  }'
```

**Example - Update both content and image**:
```bash
curl -X PATCH http://localhost:8080/userpost/11111111-1111-1111-1111-111111111111 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Final updated content",
    "post_image": {
      "base64": "iVBORw0KGgoAAAANSUhEUgAA...",
      "content_type": "image/png"
    }
  }'
```

**Response**: `{"message":"Post updated successfully"}`

**Note**: Both fields are optional. You can update just the content, just the image, or both. If updating the image, the old image file will be deleted.

### DELETE /userpost/:post_id
**Purpose**: Delete a post and its associated image file.

**Parameters**:
- `:post_id` (path) - Post UUID to delete

**Response**: `200 OK` with `{"message":"Post deleted successfully"}`, or `500 Internal Server Error` on failure.

**Example**:
```bash
curl -X DELETE http://localhost:8080/userpost/11111111-1111-1111-1111-111111111111
```

**Response**: `{"message":"Post deleted successfully"}`

---

## Notes

### Image Encoding
Posts, users, and items support image uploads via base64 encoding:
- **Upload**: Send images in request body as `{"base64":"<encoded_data>","content_type":"image/png"}`
- **Download**: GET requests return images embedded in response as `image_base64` and `image_content_type` fields
- **Supported formats**: PNG, JPEG/JPG, GIF, WEBP
- Images are stored in `resources/images/` directory

### Testing
Test data includes specific UUIDs that can be used for manual testing:
- Test user "JohnDoe": `d30869ec-fb97-46d8-85a3-82608c01f803`
- Test user "DanielGaliego": `c6554794-849f-4338-87c5-6db2e2f76514`
- Test item "Banana": `a3e1b9f2-7d94-4d3a-9b4a-111111111111`
- Test item "Tomato": `b7f2c6d4-1aeb-4f5b-9c2b-222222222222`
- Test post: `11111111-1111-1111-1111-111111111111`
- Test farm "Sunny Fields": `11111111-aaaa-aaaa-aaaa-111111111111`
