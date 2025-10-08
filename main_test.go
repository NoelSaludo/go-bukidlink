package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	w := getTestServerResponse("GET", "/ping")

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestDatabaseConnection(t *testing.T) {
	db := setupDatabase()

	assert.Equal(t, nil, db.Ping())
}

func getTestServerResponse(method, route string) *httptest.ResponseRecorder {
	router := setupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, route, nil)
	router.ServeHTTP(w, req)

	return w
}
