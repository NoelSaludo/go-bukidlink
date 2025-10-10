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

1.  **Build the application:**

    ```bash
    go build -o ./build
    ```

2.  **Run the application:**

    ```bash
    ./build/bukidlink
    ```

    The server will start on `localhost:8080`.

## API Routes

### Health Check

*   **GET /ping**

    Returns a simple "pong" message to indicate that the server is running.

    **Response:**

    ```json
    {
        "message": "pong"
    }
    ```

### User Management

*   **GET /user/:username**

    Retrieves a user by their username.

    **Parameters:**

    *   `username` (string): The username of the user to retrieve.

    **Response:**

    *   **200 OK:**

        ```json
        {
            "id": 1,
            "username": "testuser",
            "password": "password"
        }
        ```

    *   **400 Bad Request:**

        ```json
        {
            "error": "User not found"
        }
        ```

*   **POST /postuser**

    Creates a new user.

    **Request Body:**

    ```json
    {
        "username": "newuser",
        "password": "newpassword"
    }
    ```

    **Response:**

    *   **200 OK:**

        ```json
        {
            "message": "Success"
        }
        ```

    *   **400 Bad Request:**

        If the request body is invalid.
