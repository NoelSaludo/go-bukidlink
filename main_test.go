package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestDatabaseConnection(t *testing.T) {
	db := setupDatabase()

	assert.Equal(t, nil, db.Ping())
}

func TestPostUser(t *testing.T) {
	router := setupServer()

	data := User{
		Id:       1,
		Username: "John Doe",
		Password: "password123",
	}

	jData, _:= json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, "/postuser", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		`{"id":1,"password":"password123","username":"John Doe"}`,
		w.Body.String())

}
