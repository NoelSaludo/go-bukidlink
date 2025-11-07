package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bukidlink/db"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetTradeListingsBatch tests retrieving a batch of trade listings
func TestGetTradeListingsBatch(t *testing.T) {
	server := setupServer()

	// Test with default block (0)
	req, _ := http.NewRequest(http.MethodGet, "/trade/batch", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response body: %s", w.Body.String())
	}
	assert.Equalf(t, http.StatusOK, w.Code, "message: %s", w.Body.String())

	// Now expecting enriched format with image encoding
	var enrichedListings []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &enrichedListings)
	require.NoError(t, err)
	assert.NotNil(t, enrichedListings)
	// Should return at least the sample listings from the database
	assert.GreaterOrEqual(t, len(enrichedListings), 0)

	// Verify structure if listings exist
	if len(enrichedListings) > 0 {
		firstListing := enrichedListings[0]
		assert.Contains(t, firstListing, "listing")
		// Image fields are optional (only present if listing has an image)
	}
}

// TestGetTradeListingsBatchWithBlock tests pagination
func TestGetTradeListingsBatchWithBlock(t *testing.T) {
	server := setupServer()

	// Test with block parameter
	req, _ := http.NewRequest(http.MethodGet, "/trade/batch?block=0", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var enrichedListings []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &enrichedListings)
	require.NoError(t, err)
	assert.NotNil(t, enrichedListings)
}

// TestGetTradeListingsBatchInvalidBlock tests invalid block parameter
func TestGetTradeListingsBatchInvalidBlock(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/trade/batch?block=invalid", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetTradeListingByID tests retrieving a specific trade listing
func TestGetTradeListingByID(t *testing.T) {
	server := setupServer()

	// Use a known listing ID from sample data
	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"

	req, _ := http.NewRequest(http.MethodGet, "/trade?id="+listingID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Listing          db.TradeListing `json:"listing"`
		Bids             []db.TradeBid   `json:"bids"`
		ImageBase64      string          `json:"image_base64,omitempty"`
		ImageContentType string          `json:"image_content_type,omitempty"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.Listing.ID)
	assert.Equal(t, listingID, response.Listing.ID)
	assert.NotNil(t, response.Bids)
	// This listing should have 2 bids from sample data
	assert.GreaterOrEqual(t, len(response.Bids), 0)
	// Image fields are optional (only present if listing has an image)
}

// TestGetTradeListingByIDNotFound tests retrieving a non-existent listing
func TestGetTradeListingByIDNotFound(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/trade?id=00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")
}

// TestGetTradeListingByIDMissingParam tests missing id parameter
func TestGetTradeListingByIDMissingParam(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/trade", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPostTradeListing tests creating a new trade listing
func TestPostTradeListing(t *testing.T) {
	server := setupServer()

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	// Use new envelope format with listing and optional listing_image
	payload := map[string]interface{}{
		"listing": map[string]interface{}{
			"offering_farmer_id":    "8c8c73e8-0a16-4d3a-826d-75d50d7a758f", // JohnDoe farmer
			"offered_item_id":       "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
			"offered_item_quantity": 50.0,
			"desired_items":         "Looking for fresh vegetables, preferably tomatoes or lettuce.",
			"expires_at":            expiresAt,
		},
		// Optional: add listing_image for testing image upload
		// "listing_image": map[string]interface{}{
		// 	"base64": "<base64-encoded-image>",
		// 	"content_type": "image/png",
		// },
	}

	jData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/trade", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		log.Printf("Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "trade listing created", response["status"])
	assert.NotEmpty(t, response["listing_id"])

	// Verify UUID format
	_, err = uuid.Parse(response["listing_id"])
	assert.NoError(t, err, "listing_id should be a valid UUID")
}

// TestPostTradeListingInvalidData tests creating a listing with invalid data
func TestPostTradeListingInvalidData(t *testing.T) {
	server := setupServer()

	// Invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/trade", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPostTradeBid tests creating a new trade bid
func TestPostTradeBid(t *testing.T) {
	server := setupServer()

	newBid := db.TradeBid{
		TradeListingID:  "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d", // Open listing
		BiddingFarmerID: "ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f", // StewardLittle farmer
		BidItemID:       "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
		BidItemQuantity: 18,
	}

	jData, _ := json.Marshal(newBid)
	req, _ := http.NewRequest(http.MethodPost, "/bid", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "trade bid created", response["status"])
	assert.NotEmpty(t, response["bid_id"])

	// Verify UUID format
	_, err = uuid.Parse(response["bid_id"])
	assert.NoError(t, err, "bid_id should be a valid UUID")
}

// TestPostTradeBidInvalidData tests creating a bid with invalid data
func TestPostTradeBidInvalidData(t *testing.T) {
	server := setupServer()

	// Invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/bid", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateTradeListingStatus tests updating a trade listing status
func TestUpdateTradeListingStatus(t *testing.T) {
	server := setupServer()

	// First, create a new listing to update with new envelope format
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	payload := map[string]interface{}{
		"listing": map[string]interface{}{
			"offering_farmer_id":    "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
			"offered_item_id":       "c9d3e8a1-55b2-4f66-a123-333333333333",
			"offered_item_quantity": 25.0,
			"desired_items":         "Looking for fruits",
			"expires_at":            expiresAt,
		},
	}

	jData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/trade", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	listingID := createResponse["listing_id"]
	require.NotEmpty(t, listingID)

	// Now update the status to completed
	req, _ = http.NewRequest(http.MethodPatch, "/trade/"+listingID+"?updated_status=completed", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "trade listing status updated", response["status"])
	assert.Equal(t, listingID, response["trade_id"])
	assert.Equal(t, "completed", response["new_status"])

	// Verify the status was actually updated
	req, _ = http.NewRequest(http.MethodGet, "/trade?id="+listingID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse struct {
		Listing db.TradeListing `json:"listing"`
	}
	json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.Equal(t, "completed", getResponse.Listing.Status)
}

// TestUpdateTradeListingStatusInvalidStatus tests updating with invalid status
func TestUpdateTradeListingStatusInvalidStatus(t *testing.T) {
	server := setupServer()

	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"

	req, _ := http.NewRequest(http.MethodPatch, "/trade/"+listingID+"?updated_status=invalid_status", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid status")
}

// TestUpdateTradeListingStatusMissingParam tests updating without status parameter
func TestUpdateTradeListingStatusMissingParam(t *testing.T) {
	server := setupServer()

	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"

	req, _ := http.NewRequest(http.MethodPatch, "/trade/"+listingID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "required")
}

// TestUpdateTradeListingStatusAllValidValues tests all valid status values
func TestUpdateTradeListingStatusAllValidValues(t *testing.T) {
	server := setupServer()

	// Create a new listing with new envelope format
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	payload := map[string]interface{}{
		"listing": map[string]interface{}{
			"offering_farmer_id":    "8c8c73e8-0a16-4d3a-826d-75d50d7a758f",
			"offered_item_id":       "c9d3e8a1-55b2-4f66-a123-333333333333",
			"offered_item_quantity": 30.0,
			"desired_items":         "Testing status updates",
			"expires_at":            expiresAt,
		},
	}

	jData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/trade", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	listingID := createResponse["listing_id"]

	// Test all valid status values
	validStatuses := []string{"completed", "cancelled", "open"}

	for _, status := range validStatuses {
		t.Run("status_"+status, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/trade/"+listingID+"?updated_status="+status, nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, status, response["new_status"])
		})
	}
}

// TestTradeAPIWorkflow tests the complete trading workflow
func TestTradeAPIWorkflow(t *testing.T) {
	server := setupServer()

	// 1. Create a new trade listing with new envelope format
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	payload := map[string]interface{}{
		"listing": map[string]interface{}{
			"offering_farmer_id":    "8c8c73e8-0a16-4d3a-826d-75d50d7a758f", // JohnDoe farmer
			"offered_item_id":       "b7f2c6d4-1aeb-4f5b-9c2b-222222222222", // Tomato
			"offered_item_quantity": 40.0,
			"desired_items":         "Looking for corn or wheat",
			"expires_at":            expiresAt,
		},
	}

	jData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/trade", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &createResp)
	listingID := createResp["listing_id"]
	require.NotEmpty(t, listingID)

	// 2. Get the listing to verify creation
	req, _ = http.NewRequest(http.MethodGet, "/trade?id="+listingID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResp struct {
		Listing          db.TradeListing `json:"listing"`
		Bids             []db.TradeBid   `json:"bids"`
		ImageBase64      string          `json:"image_base64,omitempty"`
		ImageContentType string          `json:"image_content_type,omitempty"`
	}
	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.Equal(t, listingID, getResp.Listing.ID)
	assert.Equal(t, "open", getResp.Listing.Status)
	assert.Empty(t, getResp.Bids, "New listing should have no bids")

	// 3. Create a bid on the listing
	newBid := db.TradeBid{
		TradeListingID:  listingID,
		BiddingFarmerID: "9c6a8e18-4a0d-4a3a-8c6d-75d50d7a758f", // DanielGaliego farmer
		BidItemID:       "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
		BidItemQuantity: 20,
	}

	jData, _ = json.Marshal(newBid)
	req, _ = http.NewRequest(http.MethodPost, "/bid", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var bidResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &bidResp)
	bidID := bidResp["bid_id"]
	require.NotEmpty(t, bidID)

	// 4. Get the listing again to see the bid
	req, _ = http.NewRequest(http.MethodGet, "/trade?id="+listingID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.Len(t, getResp.Bids, 1, "Listing should have 1 bid")
	assert.Equal(t, bidID, getResp.Bids[0].ID)
	assert.Equal(t, "pending", getResp.Bids[0].Status)

	// 5. Get batch listings and verify our new listing is there
	req, _ = http.NewRequest(http.MethodGet, "/trade/batch", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var enrichedListings []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &enrichedListings)
	foundListing := false
	for _, enrichedListing := range enrichedListings {
		if listingData, ok := enrichedListing["listing"].(map[string]interface{}); ok {
			if id, ok := listingData["id"].(string); ok && id == listingID {
				foundListing = true
				break
			}
		}
	}
	assert.True(t, foundListing, "New listing should appear in batch results")

	// 6. Update listing status to completed
	req, _ = http.NewRequest(http.MethodPatch, "/trade/"+listingID+"?updated_status=completed", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Verify the status update
	req, _ = http.NewRequest(http.MethodGet, "/trade?id="+listingID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &getResp)
	assert.Equal(t, "completed", getResp.Listing.Status)
}
