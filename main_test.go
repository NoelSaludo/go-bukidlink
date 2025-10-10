package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bukidlink/db"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	server := setupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestPostUser(t *testing.T) {
	server := setupServer()

	data := db.User{
		Id:       1,
		Username: "JohnDoe",
		Password: "password123",
		Email:    "JohnDoe@example.com",
	}

	jData, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, "/postuser", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Print(w.Body.String())
}

func TestGetUser(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/user/JohnDoe", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	result := `{"id":2,"username":"JohnDoe","password":"P@ssw0rd","email":"JohnDoe@example.com"}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, result, w.Body.String())
}
