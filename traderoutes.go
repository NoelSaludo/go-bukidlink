package main

import (
	"bukidlink/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getTradeListingsBatchHandler retrieves a batch of available trade listings
// GET /trades/batch?block=0
func getTradeListingsBatchHandler(c *gin.Context) {
	blockStr := c.DefaultQuery("block", "0")
	block, err := strconv.Atoi(blockStr)
	if err != nil {
		retBadReqErr(err, c)
		return
	}

	listings, err := db.GetTradeListingsByBlock(block)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, listings)
}

// getTradeListingByIDHandler retrieves a single trade listing by ID
// GET /trade?id=uuid-here
func getTradeListingByIDHandler(c *gin.Context) {
	tradeID := c.Query("id")
	if tradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	listing, err := db.GetTradeListingByID(tradeID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	if listing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trade listing not found"})
		return
	}

	// Also fetch bids for this listing
	bids, err := db.GetTradeBidsByListingID(tradeID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listing": listing,
		"bids":    bids,
	})
}

// postTradeListingHandler creates a new trade listing
// POST /trade
// Body: { "offering_farmer_id": "uuid", "offered_item_id": "uuid", "offered_item_quantity": 20, "desired_items": "...", "expires_at": "2025-12-31T23:59:59Z" }
func postTradeListingHandler(c *gin.Context) {
	var listing db.TradeListing
	if err := c.ShouldBindJSON(&listing); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Generate UUID for the listing
	listing.ID = uuid.New().String()

	// Set initial status to open
	if listing.Status == "" {
		listing.Status = "open"
	}

	err := db.CreateTradeListing(listing)
	if err != nil {
		retConflictErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":     "trade listing created",
		"listing_id": listing.ID,
	})
}

// postTradeBidHandler creates a new trade bid
// POST /bid
// Body: { "trade_listing_id": "uuid", "bidding_farmer_id": "uuid", "bid_item_id": "uuid", "bid_item_quantity": 15 }
func postTradeBidHandler(c *gin.Context) {
	var bid db.TradeBid
	if err := c.ShouldBindJSON(&bid); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Generate UUID for the bid
	bid.ID = uuid.New().String()

	// Set initial status to pending
	if bid.Status == "" {
		bid.Status = "pending"
	}

	err := db.CreateTradeBid(bid)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "trade bid created",
		"bid_id": bid.ID,
	})
}

// updateTradeListingStatusHandler updates the status of a trade listing
// PATCH /trade/:id?updated_status=completed
func updateTradeListingStatusHandler(c *gin.Context) {
	tradeID := c.Param("id")
	newStatus := c.Query("updated_status")

	if newStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "updated_status query parameter is required"})
		return
	}

	// Validate status values (open, completed, cancelled)
	if newStatus != "open" && newStatus != "completed" && newStatus != "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Must be: open, completed, or cancelled"})
		return
	}

	err := db.UpdateTradeListingStatus(tradeID, newStatus)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "trade listing status updated",
		"trade_id":   tradeID,
		"new_status": newStatus,
	})
}

// getTradeBidByIDHandler retrieves a single trade bid by ID
// GET /bid?id=uuid-here
func getTradeBidByIDHandler(c *gin.Context) {
	bidID := c.Query("id")
	if bidID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	bid, err := db.GetTradeBidByID(bidID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	if bid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "trade bid not found"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// getTradeBidsByFarmerHandler retrieves all bids made by a specific farmer
// GET /bids/farmer/:farmer_id
func getTradeBidsByFarmerHandler(c *gin.Context) {
	farmerID := c.Param("farmer_id")

	bids, err := db.GetTradeBidsByFarmer(farmerID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, bids)
}

// updateTradeBidHandler updates a trade bid's item and quantity
// PUT /bid/:id
// Body: { "bid_item_id": "uuid", "bid_item_quantity": 20 }
func updateTradeBidHandler(c *gin.Context) {
	bidID := c.Param("id")

	var req struct {
		BidItemID       string  `json:"bid_item_id" binding:"required"`
		BidItemQuantity float64 `json:"bid_item_quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	err := db.UpdateTradeBid(bidID, req.BidItemID, req.BidItemQuantity)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "bid not found or not in pending status"})
			return
		}
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "trade bid updated",
		"bid_id": bidID,
	})
}

// updateTradeBidStatusHandler updates the status of a trade bid
// PATCH /bid/:id/status?updated_status=accepted
func updateTradeBidStatusHandler(c *gin.Context) {
	bidID := c.Param("id")
	newStatus := c.Query("updated_status")

	if newStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "updated_status query parameter is required"})
		return
	}

	// Validate status values (pending, accepted, rejected)
	if newStatus != "pending" && newStatus != "accepted" && newStatus != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status. Must be: pending, accepted, or rejected"})
		return
	}

	err := db.UpdateTradeBidStatus(bidID, newStatus)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "trade bid status updated",
		"bid_id":     bidID,
		"new_status": newStatus,
	})
}

// deleteTradeBidHandler deletes a trade bid
// DELETE /bid/:id
func deleteTradeBidHandler(c *gin.Context) {
	bidID := c.Param("id")

	err := db.DeleteTradeBid(bidID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "bid not found or not in pending status"})
			return
		}
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "trade bid deleted",
		"bid_id": bidID,
	})
}
