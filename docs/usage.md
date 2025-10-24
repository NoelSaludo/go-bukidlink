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
- Routes are implemented in `main.go`. This document lists the currently implemented endpoints and how to use them.

---

## Routes

### GET /ping
- Purpose: Health check.
- Request: none
- Response:
  - Status: 200 OK
  - Body: JSON

Example response:
```json
{"message":"pong"}
```

Curl example:
```bash
curl -i http://localhost:8080/ping
```

---

### GET /item/:block
- Purpose: Return a page/block of items (the tests call `/item/0`).
- Path parameters:
  - `:block` (string/int) — page or block index. The tests use `0` to request the first block.
- Request body: none
- Response:
  - Status: 200 OK on success
  - Body: JSON array of Item objects

Item object shape (from `db/models.go`):
```json
{
  "id": "<uuid>",
  "name": "<string>",
  "description": "<string>",
  "amount": <int>,
  "costPKilo": <float>,
  "category": "<fruits|vegetables|grains|livestock|dairy|others>",
  "rating": <float> // if present from DB
}
```

Curl example (first block):
```bash
curl -i http://localhost:8080/item/0
```

---

### GET /item/category/:category
- Purpose: Return items matching a category.
- Path parameters:
  - `:category` — one of: `fruits`, `vegetables`, `grains`, `livestock`, `dairy`, `others`.
- Request body: none
- Response:
  - Status: 200 OK and JSON array of `Item` objects for a valid category (the tests assert non-empty arrays for each valid category).
  - Status: 500 Internal Server Error if the category is invalid (see tests: category `nothing` returns 500).

Curl example:
```bash
curl -i http://localhost:8080/item/category/fruits
```

Note: The tests expect each returned item's `category` field to equal the requested category.

---

### GET /user/:username
- Purpose: Retrieve a user by username.
- Path parameters:
  - `:username` — the user's username (e.g. `JohnDoe`).
- Request body: none
- Response:
  - Status: 200 OK on success
  - Body: JSON representation of a `User` object

User object shape (from `db/models.go`):
```json
{
  "id": "<uuid>",
  "username": "<string>",
  "password": "<string>",
  "email": "<string>"
}
```

Curl example:
```bash
curl -i http://localhost:8080/user/JohnDoe
```

---

### POST /user
- Purpose: Insert a new user.
- Request headers:
  - `Content-Type: application/json`
- Request body: JSON `User` object. Example:
```json
{
  "id": "01d85ea5-0c1f-457c-b1f5-04f4e48b54b6",
  "username": "JohnDoe",
  "password": "password123",
  "email": "JohnDoe@example.com"
}
```
- Response:
  - Status: 201 Created (expected when a new user is inserted) OR
  - Status: 409 Conflict if the user already exists (the tests submit a known user and expect `409`).

Curl example (in tests they expect conflict):
```bash
curl -i -X POST -H "Content-Type: application/json" -d '{"id":"01d85ea5-0c1f-457c-b1f5-04f4e48b54b6","username":"JohnDoe","password":"password123","email":"JohnDoe@example.com"}' http://localhost:8080/user
```

---

### GET /comment/:itemId
- Purpose: Return comments for a given item id.
- Path parameters:
  - `:itemId` — an item UUID (e.g. `a3e1b9f2-7d94-4d3a-9b4a-111111111111`).
- Request body: none
- Response:
  - Status: 200 OK on success (tests expect non-empty comments for the sample item id above).
  - Body: JSON array of `Comment` objects

Comment object shape (from `db/models.go`):
```json
{
  "id": <int>,
  "itemid": <int>,
  "userid": <int>,
  "content": "<string>",
  "rating": <float>
}
```

Curl example:
```bash
curl -i http://localhost:8080/comment/a3e1b9f2-7d94-4d3a-9b4a-111111111111
```

---

### POST /comment/:productID
- Note: `main.go` contains a route registration `commentG.POST("/:productID")` but no handler function is provided. That means the POST route is declared in the router, but the project currently does not implement the handler to accept or create comments. Attempting to call this endpoint will either produce a server error or will not behave as intended until the handler is implemented.

Action item: Implement the handler (e.g. `postCommentHandler(c *gin.Context)`) and register it:
```go
commentG.POST("/:productID", postCommentHandler)
```