# BukidLink API

This is a simple Go API for BukidLink.

## Prerequisites

*   [Go](https://golang.org/dl/) version 1.24.8 or higher
*   [PostgreSQL](https://www.postgresql.org/download/)

## Installation

1.  **Clone the repository:**

    ```bash
    git clone https://your-repository-url/bukidlink.git
    cd bukidlink
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Set up the database:**

    This project uses PostgreSQL. Make sure you have a running instance of PostgreSQL.

    The application requires the following environment variables to connect to the database:

    *   `DBHOST`: The host of your database (e.g., `localhost`)
    *   `DBPORT`: The port of your database (e.g., `5432`)
    *   `DBUSER`: The username for your database (e.g., "bukidlink-user")
    *   `DBPASSWORD`: The password for your database (e.g., "bukidlinkpass123")
    *   `DATABASE`: The name of your database (e.g, "BukidLink")

    You can set these environment variables in your shell or using a `.env` file.

    **Initialize the database schema and sample data:**

    Run the following command to set up all tables, types, and initial data:

    ```bash
    psql -U <your_db_user> -d <your_db_name> -f sql/backup.sql
    ```

    Replace `<your_db_user>` and `<your_db_name>` with your PostgreSQL username and database name.

## Usage

This section gives a concise summary of the server routes, examples of how to call them with `curl`, and what to expect in request and response bodies. For full, expanded documentation (detailed request/response shapes, examples and notes) see `docs/usage.md`.

### Quick start

```bash
# run from project root
go run main.go
```

Server address: http://localhost:8080

### API Routes Overview

The BukidLink API provides comprehensive endpoints for managing an agricultural marketplace platform:

#### Health Check
- **GET /ping** - Health check endpoint

#### Item Management
- **GET /item/:block** - Retrieve paginated items (100 per page)
- **GET /item/category/:category** - List items by category (fruits, vegetables, grains, livestock, dairy, others)
- **GET /item?id=<uuid>** - Fetch single item by UUID
- **POST /item** - Create new item with optional image upload

#### User Management
- **GET /user/:username** - Fetch user profile by username
- **POST /user** - Create new user with details and optional profile picture

#### Order Management
- **GET /order/:user_id** - Get all orders for a specific user
- **POST /order** - Create new order with multiple items
- **PATCH /order/status** - Update order status (Packaging, Shipping, Delivered, etc.)
- **DELETE /order** - Delete an order and all associated items

#### Shopping Cart
- **GET /cart/:user_id** - Get user's cart (auto-creates if doesn't exist)
- **POST /cart/item** - Add item to cart or update quantity
- **DELETE /cart/item/:cart_item_id** - Remove item from cart

#### Reviews
- **GET /review/:itemId** - Get all reviews for an item

#### User Posts (Farmer Updates)
- **GET /userpost/:user_id** - Get all posts by a farmer
- **GET /userpost/post/:post_id** - Get specific post by ID
- **POST /userpost** - Create new post with optional image
- **PATCH /userpost/:post_id** - Update post content and/or image
- **DELETE /userpost/:post_id** - Delete post and associated image

#### Trade Listings & Bids
- **GET /trade/batch** - Get batch of trade listings
- **GET /trade?id=<uuid>** - Get specific trade listing
- **POST /trade** - Create new trade listing
- **PATCH /trade/:id** - Update trade listing status
- **GET /bid?id=<uuid>** - Get specific trade bid
- **POST /bid** - Create new trade bid
- **PUT /bid/:id** - Update trade bid
- **PATCH /bid/:id/status** - Update bid status
- **DELETE /bid/:id** - Delete trade bid
- **GET /bid/farmer/:farmer_id** - Get all bids by a farmer

#### Payment & Balance Management
- **GET /balance/:user_id** - Get user's balance
- **POST /balance** - Create user balance account
- **POST /payment/deposit** - Process deposit transaction
- **POST /payment/withdrawal** - Process withdrawal transaction
- **POST /payment/order** - Process order payment
- **POST /payment/refund** - Process refund
- **GET /payment/transaction/:transaction_id** - Get transaction details
- **GET /payment/transactions/:user_id** - Get user's transaction history

#### Real-time Communication
- **Chat routes** - Real-time messaging between users
- **GET /ws** - WebSocket endpoint for real-time updates

### Example Usage

#### Health Check
```bash
curl http://localhost:8080/ping
```

#### Get Items
```bash
# Get first 100 items
curl http://localhost:8080/item/0

# Get items by category
curl http://localhost:8080/item/category/fruits

# Get specific item
curl "http://localhost:8080/item?id=a3e1b9f2-7d94-4d3a-9b4a-111111111111"
```

#### Create User
```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "username": "JohnDoe",
      "password": "password123",
      "email": "johndoe@example.com",
      "details": {
        "first_name": "John",
        "last_name": "Doe",
        "address": "123 Main St",
        "contact_number": "1234567890"
      }
    },
    "profile_pic": {
      "base64": "<base64-encoded-image>",
      "content_type": "image/png"
    }
  }'
```

#### Create Order
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
      }
    ]
  }'
```

#### Manage Shopping Cart
```bash
# Get cart
curl http://localhost:8080/cart/c6554794-849f-4338-87c5-6db2e2f76514

# Add item to cart
curl -X POST http://localhost:8080/cart/item \
  -H "Content-Type: application/json" \
  -d '{
    "cart_id": "31111111-1111-1111-1111-111111111111",
    "item_id": "a3e1b9f2-7d94-4d3a-9b4a-111111111111",
    "quantity": 3
  }'
```

### Features

- **Image Support**: Items, users, and posts support base64-encoded image uploads (PNG, JPEG, GIF, WEBP)
- **Transaction Safety**: Order creation, cart management, and payments use database transactions
- **Real-time Updates**: WebSocket support for live notifications and chat
- **Comprehensive Error Handling**: Appropriate HTTP status codes and error messages
- **UUID-based IDs**: All entities use UUID v4 for unique identification

### Documentation

For detailed API documentation including request/response schemas, error handling, and additional examples, see:
- **Full API Documentation**: `docs/usage.md`
- **Dart/Flutter Client Examples**: `docs/usage_dart.md`
