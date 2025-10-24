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

## Usage

This section gives a concise summary of the server routes, examples of how to call them with `curl`, and what to expect in request and response bodies. For full, expanded documentation (detailed request/response shapes, examples and notes) see `docs/usage.md`.

Quick start

```bash
# run from project root
go run main.go
```

Server address: http://localhost:8080

Primary routes (summary)

- GET /ping
    - Purpose: health check
    - Response: 200 OK, JSON: {"message":"pong"}
    - Example: `curl -i http://localhost:8080/ping`

- GET /item/:block
    - Purpose: retrieve a page/block of items (tests use `/item/0`)
    - Path param: `:block` (page index)
    - Response: 200 OK, JSON array of Item objects
    - Example: `curl -i http://localhost:8080/item/0`

- GET /item/category/:category
    - Purpose: list items in a category
    - Path param: `:category` — one of `fruits, vegetables, grains, livestock, dairy, others`
    - Response: 200 OK with JSON items for valid categories; 500 for invalid category
    - Example: `curl -i http://localhost:8080/item/category/fruits`

- GET /user/:username
    - Purpose: fetch a user by username
    - Path param: `:username`
    - Response: 200 OK, JSON User object
    - Example: `curl -i http://localhost:8080/user/JohnDoe`

- POST /user
    - Purpose: insert a new user
    - Headers: `Content-Type: application/json`
    - Body: JSON `User` object (see `docs/usage.md` for an example)
    - Responses: 201 Created on success, 409 Conflict if user exists (tests expect conflict for a known user)
    - Example:
        ```bash
        curl -i -X POST -H "Content-Type: application/json" \
            -d '{"id":"01d85ea5-0c1f-457c-b1f5-04f4e48b54b6","username":"JohnDoe","password":"password123","email":"JohnDoe@example.com"}' \
            http://localhost:8080/user
        ```

- GET /comment/:itemId
    - Purpose: list comments for an item
    - Path param: `:itemId` (UUID)
    - Response: 200 OK, JSON array of Comment objects
    - Example: `curl -i http://localhost:8080/comment/a3e1b9f2-7d94-4d3a-9b4a-111111111111`

- POST /comment/:productID
    - Note: route is registered in `main.go` but no handler is implemented yet. Calling this endpoint will not create comments until a handler is provided.

Notes & pointers

- The server calls `db.SetupDatabase()` on startup; ensure your Postgres instance is available and environment variables for DB connection are set.
- Example UUIDs and expected behaviors are taken from tests (`main_test.go` and `db_test.go`) — use those values for manual testing.
- Full request/response shapes, additional examples, and implementation notes are available in `docs/usage.md`.

If you want, I can fold the entire `docs/usage.md` into README or add example JSON responses for each endpoint inline.