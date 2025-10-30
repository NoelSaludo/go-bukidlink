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
- Purpose: Health check.
- Response:
  - Status: 200 OK
  - Body: `{"message":"pong"}`

---

## Item Routes

### GET /item/:block
- Purpose: Return a page/block of items.
- Response: Array of `Item` objects.

Example Response:
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
- Purpose: Return items matching a category.
- Response: Array of `Item` objects.

Example Response (for `/item/category/fruits`):
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
- Purpose: Retrieve a single item by its ID.
- Response: A single `Item` object.

Example Response:
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
- Purpose: Retrieve a user by username.
- Response: A single `User` object.

Example Response:
```json
{
    "id": "d30869ec-fb97-46d8-85a3-82608c01f803",
    "username": "JohnDoe",
    "password": "P@ssw0rd",
    "email": "JohnDoe@example.com",
    "address": "So. Pinamuntasan, Brgy. Aga, Nasugbu, Batangas"
}
```

### POST /user
- Purpose: Insert a new user.
- Response: `{"message":"Success"}`

---

## Order Routes

### GET /order/:user_id
- Purpose: Get all orders for a specific user.
- Response: JSON array of `Order` objects, with nested `OrderItem`s.

Example Response:
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
- Purpose: Create a new order.
- Response: `{"status":"order created","order_id":"<NEW_ORDER_UUID>"}`

### POST /order/status
- Purpose: Update the status of an order.
- Response: `{"status":"order status updated"}`

---

## Cart Routes

### GET /cart/:user_id
- Purpose: Get a user's cart. Creates a new cart if one doesn't exist.
- Response: JSON `Cart` object with nested `CartItem`s.

Example Response:
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
- Purpose: Add an item to a cart.
- Response: `{"status":"item added to cart"}`

### DELETE /cart/item/:cart_item_id
- Purpose: Remove an item from a cart.
- Response: `{"status":"item removed from cart"}`

---

## Comment Routes

### GET /comment/:itemId
- Purpose: Return comments for a given item id.
- Response: Array of `Review` objects.

Example Response:
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

### POST /comment/:productID
- Note: This route is declared but not implemented.