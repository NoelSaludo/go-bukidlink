package main

import (
	"bytes"
	"encoding/json"
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

func TestDatabaseConnection(t *testing.T) {
	err := db.SetupDatabase()

	assert.Equal(t, nil, err)
	assert.Equal(t, nil, db.Ping())
}

func TestPostUser(t *testing.T) {
	server := setupServer()

	data := db.User{
		Id:       1,
		Username: "JohnDoe",
		Password: "password123",
	}

	jData, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, "/postuser", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		`{"id":1,"password":"password123","username":"JohnDoe"}`,
		w.Body.String())

}

func TestGetUser(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/user/JohnDoe", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	result := `{"id":1,"username":"JohnDoe","password":"password123"}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, result, w.Body.String())
}
