package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"bukidlink/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Id:       "01d85ea5-0c1f-457c-b1f5-04f4e48b54b6",
		Username: "JohnDoe",
		Password: "password123",
		Email:    "JohnDoe@example.com",
	}

	jData, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetUser(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/user/JohnDoe", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user db.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.NotEmpty(t, user)
}

func Test100Items(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/item/0", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetItemsByCategory(t *testing.T) {
	server := setupServer()

	categories := []string{"fruits", "vegetables", "grains", "livestock", "dairy"}
	for _, cat := range categories {
		req, _ := http.NewRequest(http.MethodGet, "/item/category/"+cat, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		var items []db.Item
		err := json.Unmarshal(w.Body.Bytes(), &items)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, i := range items {
			assert.Equal(t, cat, i.Category)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, "/item/category/nothing", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetComments(t *testing.T) {
	s := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/comment/a3e1b9f2-7d94-4d3a-9b4a-111111111111", nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var comments []db.Comment
	if err := json.Unmarshal(w.Body.Bytes(), &comments); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, comments)
}
