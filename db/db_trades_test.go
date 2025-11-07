package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetTradeListingByID tests retrieving an existing trade listing
func TestGetTradeListingByID(t *testing.T) {
	_ = SetupDatabase()

	// Test getting an existing trade listing (from sample data)
	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"
	listing, err := GetTradeListingByID(listingID)

	require.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, listingID, listing.ID)
	assert.Equal(t, "8c8c73e8-0a16-4d3a-826d-75d50d7a758f", listing.OfferingFarmerID) // JohnDoe farmer
	assert.Equal(t, "b7f2c6d4-1aeb-4f5b-9c2b-222222222222", listing.OfferedItemID)    // Tomato
	assert.Equal(t, float64(20), listing.OfferedItemQuantity)
	assert.Equal(t, "open", listing.Status)
	assert.Contains(t, listing.DesiredItems, "rice")
	assert.NotNil(t, listing.ExpiresAt)
}

// TestGetTradeListingByID_NotFound tests retrieving a non-existent listing
func TestGetTradeListingByID_NotFound(t *testing.T) {
	_ = SetupDatabase()

	// Test getting a non-existent trade listing
	listing, err := GetTradeListingByID("00000000-0000-0000-0000-000000000000")

	require.NoError(t, err)
	assert.Nil(t, listing)
}

// TestGetTradeBidsByListingID tests retrieving bids for a listing
func TestGetTradeBidsByListingID(t *testing.T) {
	_ = SetupDatabase()

	// Test getting bids for listing 1 (should have 2 bids from sample data)
	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"
	bids, err := GetTradeBidsByListingID(listingID)

	require.NoError(t, err)
	assert.NotNil(t, bids)
	assert.Len(t, bids, 2, "Listing should have 2 bids")

	// Verify bid details
	for _, bid := range bids {
		assert.NotEmpty(t, bid.ID)
		assert.Equal(t, listingID, bid.TradeListingID)
		assert.NotEmpty(t, bid.BiddingFarmerID)
		assert.Equal(t, "c9d3e8a1-55b2-4f66-a123-333333333333", bid.BidItemID) // Rice
		assert.Equal(t, "pending", bid.Status)
		assert.Greater(t, bid.BidItemQuantity, float64(0))
	}
}

// TestGetTradeBidsByListingID_NoBids tests retrieving bids for a listing with no bids
func TestGetTradeBidsByListingID_NoBids(t *testing.T) {
	_ = SetupDatabase()

	// Test getting bids for listing 3 (cancelled listing with no bids)
	listingID := "3c4d5e6f-7a8b-9c0d-1e2f-3a4b5c6d7e8f"
	bids, err := GetTradeBidsByListingID(listingID)

	require.NoError(t, err)
	assert.NotNil(t, bids)
	assert.Empty(t, bids, "Cancelled listing should have no bids")
}

// TestCreateAndDeleteTradeListing tests the full workflow of creating and managing a listing
func TestCreateAndDeleteTradeListing(t *testing.T) {
	_ = SetupDatabase()

	// Create a new trade listing
	newListingID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // Expires in 7 days
	newListing := TradeListing{
		ID:                  newListingID,
		OfferingFarmerID:    "8c8c73e8-0a16-4d3a-826d-75d50d7a758f", // JohnDoe farmer
		OfferedItemID:       "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
		OfferedItemQuantity: 50,
		DesiredItems:        "Looking for fresh vegetables, preferably tomatoes or lettuce.",
		Status:              "open",
		ExpiresAt:           &expiresAt,
	}

	err := CreateTradeListing(newListing)
	require.NoError(t, err)

	// Verify the listing was created
	createdListing, err := GetTradeListingByID(newListingID)
	require.NoError(t, err)
	assert.NotNil(t, createdListing)
	assert.Equal(t, newListingID, createdListing.ID)
	assert.Equal(t, "8c8c73e8-0a16-4d3a-826d-75d50d7a758f", createdListing.OfferingFarmerID)
	assert.Equal(t, "c9d3e8a1-55b2-4f66-a123-333333333333", createdListing.OfferedItemID)
	assert.Equal(t, float64(50), createdListing.OfferedItemQuantity)
	assert.Equal(t, "open", createdListing.Status)
	assert.Contains(t, createdListing.DesiredItems, "vegetables")

	// Test updating the listing status to completed
	err = UpdateTradeListingStatus(newListingID, "completed")
	require.NoError(t, err)

	// Verify the status was updated
	updatedListing, err := GetTradeListingByID(newListingID)
	require.NoError(t, err)
	assert.Equal(t, "completed", updatedListing.Status)

	// Clean up: Delete the listing by updating status to a marker or using direct SQL
	// Since there's no DeleteTradeListing function, we'll update to cancelled
	err = UpdateTradeListingStatus(newListingID, "cancelled")
	require.NoError(t, err)
}

// TestCreateAndManageTradeBid tests the full workflow of creating and managing bids
func TestCreateAndManageTradeBid(t *testing.T) {
	_ = SetupDatabase()

	// Use an existing open listing
	listingID := "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"

	// Create a new bid
	newBidID := uuid.New().String()
	newBid := TradeBid{
		ID:              newBidID,
		TradeListingID:  listingID,
		BiddingFarmerID: "ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f", // StewardLittle farmer
		BidItemID:       "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
		BidItemQuantity: 18,
		Status:          "pending",
	}

	err := CreateTradeBid(newBid)
	require.NoError(t, err)

	// Verify the bid was created
	bids, err := GetTradeBidsByListingID(listingID)
	require.NoError(t, err)

	var foundBid *TradeBid
	for _, bid := range bids {
		if bid.ID == newBidID {
			foundBid = &bid
			break
		}
	}

	require.NotNil(t, foundBid, "Created bid should be found")
	assert.Equal(t, newBidID, foundBid.ID)
	assert.Equal(t, listingID, foundBid.TradeListingID)
	assert.Equal(t, "ac7b9f28-5b1e-4b4a-9d7e-85d50d7a758f", foundBid.BiddingFarmerID)
	assert.Equal(t, float64(18), foundBid.BidItemQuantity)
	assert.Equal(t, "pending", foundBid.Status)

	// Test updating the bid status to accepted
	err = UpdateTradeBidStatus(newBidID, "accepted")
	require.NoError(t, err)

	// Verify the status was updated
	updatedBids, err := GetTradeBidsByListingID(listingID)
	require.NoError(t, err)

	var updatedBid *TradeBid
	for _, bid := range updatedBids {
		if bid.ID == newBidID {
			updatedBid = &bid
			break
		}
	}

	require.NotNil(t, updatedBid)
	assert.Equal(t, "accepted", updatedBid.Status)

	// Test updating to rejected
	err = UpdateTradeBidStatus(newBidID, "rejected")
	require.NoError(t, err)

	// Verify rejection
	rejectedBids, err := GetTradeBidsByListingID(listingID)
	require.NoError(t, err)

	var rejectedBid *TradeBid
	for _, bid := range rejectedBids {
		if bid.ID == newBidID {
			rejectedBid = &bid
			break
		}
	}

	require.NotNil(t, rejectedBid)
	assert.Equal(t, "rejected", rejectedBid.Status)
}

// TestTradeListingStatuses tests different listing statuses
func TestTradeListingStatuses(t *testing.T) {
	_ = SetupDatabase()

	testCases := []struct {
		listingID      string
		expectedStatus string
		description    string
	}{
		{
			listingID:      "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
			expectedStatus: "open",
			description:    "Open listing for tomatoes",
		},
		{
			listingID:      "2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
			expectedStatus: "completed",
			description:    "Completed listing for handicraft baskets",
		},
		{
			listingID:      "3c4d5e6f-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
			expectedStatus: "cancelled",
			description:    "Cancelled listing for bananas",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			listing, err := GetTradeListingByID(tc.listingID)
			require.NoError(t, err)
			assert.NotNil(t, listing)
			assert.Equal(t, tc.expectedStatus, listing.Status)
		})
	}
}

// TestTradeBidStatuses tests different bid statuses
func TestTradeBidStatuses(t *testing.T) {
	_ = SetupDatabase()

	testCases := []struct {
		bidID          string
		expectedStatus string
		description    string
	}{
		{
			bidID:          "4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
			expectedStatus: "pending",
			description:    "Pending bid from DanielGaliego",
		},
		{
			bidID:          "6f7a8b9c-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
			expectedStatus: "accepted",
			description:    "Accepted bid from StewardLittle",
		},
		{
			bidID:          "7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d",
			expectedStatus: "rejected",
			description:    "Rejected bid from JohnDoe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// We need to get all bids and find the one we're looking for
			// Since we know listing 2 has accepted/rejected bids, check that one
			listingID := "2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e"
			if tc.expectedStatus == "pending" {
				listingID = "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"
			}

			bids, err := GetTradeBidsByListingID(listingID)
			require.NoError(t, err)

			var foundBid *TradeBid
			for _, bid := range bids {
				if bid.ID == tc.bidID {
					foundBid = &bid
					break
				}
			}

			require.NotNil(t, foundBid, "Bid should be found")
			assert.Equal(t, tc.expectedStatus, foundBid.Status)
		})
	}
}

// TestUpdateTradeListingStatus_InvalidID tests updating a non-existent listing
func TestUpdateTradeListingStatus_InvalidID(t *testing.T) {
	_ = SetupDatabase()

	// Updating a non-existent listing should not error but will affect 0 rows
	err := UpdateTradeListingStatus("00000000-0000-0000-0000-000000000000", "open")
	assert.NoError(t, err, "Update with invalid ID should not error")
}

// TestUpdateTradeBidStatus_InvalidID tests updating a non-existent bid
func TestUpdateTradeBidStatus_InvalidID(t *testing.T) {
	_ = SetupDatabase()

	// Updating a non-existent bid should not error but will affect 0 rows
	err := UpdateTradeBidStatus("00000000-0000-0000-0000-000000000000", "accepted")
	assert.NoError(t, err, "Update with invalid ID should not error")
}
