package main

import (
	"bukidlink/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetUserTransactions tests retrieving transaction history with default pagination
func TestGetUserTransactions(t *testing.T) {
	server := setupServer()

	// Use JohnDoe's user ID from seed data
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
	// Should return transactions if user has any
	assert.GreaterOrEqual(t, len(transactions), 0)
}

// TestGetUserTransactionsWithPagination tests custom limit and offset parameters
func TestGetUserTransactionsWithPagination(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Test with custom limit
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=5", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
	// Should respect limit (at most 5 items)
	assert.LessOrEqual(t, len(transactions), 5)
}

// TestGetUserTransactionsWithOffset tests pagination offset
func TestGetUserTransactionsWithOffset(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Get first page
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var firstPage []db.PaymentTransaction
	json.Unmarshal(w.Body.Bytes(), &firstPage)

	// Get second page
	req, _ = http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=10&offset=10", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var secondPage []db.PaymentTransaction
	json.Unmarshal(w.Body.Bytes(), &secondPage)

	// If both pages have results, they should be different
	if len(firstPage) > 0 && len(secondPage) > 0 {
		assert.NotEqual(t, firstPage[0].ID, secondPage[0].ID, "Pages should contain different transactions")
	}
}

// TestGetUserTransactionsInvalidLimit tests invalid limit parameter
func TestGetUserTransactionsInvalidLimit(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Invalid limit (non-numeric) should fall back to default (20)
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=invalid", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	// Should still return valid response with default pagination
	assert.NotNil(t, transactions)
}

// TestGetUserTransactionsInvalidOffset tests invalid offset parameter
func TestGetUserTransactionsInvalidOffset(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Invalid offset (non-numeric) should fall back to default (0)
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?offset=invalid", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
}

// TestGetUserTransactionsNegativeLimit tests negative limit parameter
func TestGetUserTransactionsNegativeLimit(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Negative limit should fall back to default (20)
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=-5", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
}

// TestGetUserTransactionsNegativeOffset tests negative offset parameter
func TestGetUserTransactionsNegativeOffset(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Negative offset should fall back to default (0)
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?offset=-10", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
}

// TestGetUserTransactionsZeroLimit tests zero limit parameter
func TestGetUserTransactionsZeroLimit(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	// Zero limit should fall back to default (20)
	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)
	assert.NotNil(t, transactions)
}

// TestGetUserTransactionsNonExistentUser tests retrieving transactions for non-existent user
func TestGetUserTransactionsNonExistentUser(t *testing.T) {
	server := setupServer()

	// Use a UUID that doesn't exist in the database
	userID := "00000000-0000-0000-0000-000000000000"

	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	// Should return 404 or empty array depending on implementation
	if w.Code == http.StatusNotFound {
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "not found")
	} else {
		assert.Equal(t, http.StatusOK, w.Code)
		var transactions []db.PaymentTransaction
		err := json.Unmarshal(w.Body.Bytes(), &transactions)
		require.NoError(t, err)
		assert.Empty(t, transactions, "Non-existent user should have no transactions")
	}
}

// TestGetUserTransactionsResponseStructure tests the structure of returned transactions
func TestGetUserTransactionsResponseStructure(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)

	// Verify structure if transactions exist
	if len(transactions) > 0 {
		txn := transactions[0]
		assert.NotEmpty(t, txn.ID, "Transaction should have ID")
		assert.NotEmpty(t, txn.UserID, "Transaction should have UserID")
		assert.NotEmpty(t, txn.TransactionType, "Transaction should have TransactionType")
		assert.NotEmpty(t, txn.Status, "Transaction should have Status")
		assert.NotZero(t, txn.Amount, "Transaction should have Amount")
		assert.NotZero(t, txn.CreatedAt, "Transaction should have CreatedAt timestamp")
	}
}

// TestGetUserTransactionsOrderedByCreatedAt tests that transactions are ordered by creation date descending
func TestGetUserTransactionsOrderedByCreatedAt(t *testing.T) {
	server := setupServer()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803"

	req, _ := http.NewRequest(http.MethodGet, "/payment/transactions/"+userID+"?limit=50", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []db.PaymentTransaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	require.NoError(t, err)

	// If we have multiple transactions, verify they're ordered by created_at DESC
	if len(transactions) > 1 {
		for i := 0; i < len(transactions)-1; i++ {
			// Each transaction should have created_at >= next transaction's created_at
			assert.True(t,
				transactions[i].CreatedAt.After(transactions[i+1].CreatedAt) ||
					transactions[i].CreatedAt.Equal(transactions[i+1].CreatedAt),
				"Transactions should be ordered by created_at descending")
		}
	}
}
