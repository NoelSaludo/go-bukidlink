package main

import (
	"bukidlink/db"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// encodeTradeImageToBase64 reads an image file and returns base64 encoded data and content type
func encodeTradeImageToBase64(imageURL *string) (string, string, error) {
	if imageURL == nil || *imageURL == "" {
		return "", "", nil
	}

	imageData, err := os.ReadFile(*imageURL)
	if err != nil {
		return "", "", err
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Determine content type from file extension
	ext := strings.ToLower(filepath.Ext(*imageURL))
	contentType := "image/jpeg" // default
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	return base64Image, contentType, nil
}

// decodeAndSaveTradeImage decodes base64 image data and saves it to a file
func decodeAndSaveTradeImage(base64Data string, contentType string, filename string) (string, error) {
	// Decode base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Determine file extension from content type
	ext := ".jpg" // default
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}

	// Create file path
	filePath := filepath.Join("resources", "images", filename+ext)

	// Write file to disk
	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

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

	// Enrich listings with base64 encoded images
	enrichedListings := make([]map[string]interface{}, 0)
	for _, listing := range listings {
		enrichedListing := map[string]interface{}{
			"listing": listing,
		}

		// Encode image if available
		if listing.ImageURL != nil && *listing.ImageURL != "" {
			if base64Image, contentType, err := encodeTradeImageToBase64(listing.ImageURL); err == nil && base64Image != "" {
				enrichedListing["image_base64"] = base64Image
				enrichedListing["image_content_type"] = contentType
			}
		}

		enrichedListings = append(enrichedListings, enrichedListing)
	}

	c.JSON(http.StatusOK, enrichedListings)
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

	// Prepare response with image encoding
	response := gin.H{
		"listing": listing,
		"bids":    bids,
	}

	// Encode image if available
	if listing.ImageURL != nil && *listing.ImageURL != "" {
		if base64Image, contentType, err := encodeTradeImageToBase64(listing.ImageURL); err == nil && base64Image != "" {
			response["image_base64"] = base64Image
			response["image_content_type"] = contentType
		}
	}

	c.JSON(http.StatusOK, response)
}

// postTradeListingHandler creates a new trade listing
// POST /trade
//
//	Body: {
//	  "listing": { "offering_farmer_id": "uuid", "offered_item_id": "uuid", "offered_item_quantity": 20, "desired_items": "...", "expires_at": "2025-12-31T23:59:59Z" },
//	  "listing_image": { "base64": "...", "content_type": "image/png" }
//	}
func postTradeListingHandler(c *gin.Context) {
	var payload struct {
		Listing struct {
			OfferingFarmerID    string  `json:"offering_farmer_id" binding:"required"`
			OfferedItemID       string  `json:"offered_item_id" binding:"required"`
			OfferedItemQuantity float64 `json:"offered_item_quantity" binding:"required"`
			DesiredItems        string  `json:"desired_items" binding:"required"`
			ExpiresAt           *string `json:"expires_at"`
		} `json:"listing"`
		ListingImage *struct {
			Base64      string `json:"base64"`
			ContentType string `json:"content_type"`
		} `json:"listing_image"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Generate UUID for the listing
	listingID := uuid.New().String()

	// Handle image if provided
	var imageURL *string
	if payload.ListingImage != nil && payload.ListingImage.Base64 != "" && payload.ListingImage.ContentType != "" {
		// Save image with listing ID as filename
		filePath, err := decodeAndSaveTradeImage(payload.ListingImage.Base64, payload.ListingImage.ContentType, listingID+"_trade")
		if err != nil {
			retConflictErr(err, c)
			return
		}
		imageURL = &filePath
	} else {
		// Set default image if no image provided
		defaultImage := "resources/images/no-image.jpg"
		imageURL = &defaultImage
	}

	// Create listing object
	listing := db.TradeListing{
		ID:                  listingID,
		OfferingFarmerID:    payload.Listing.OfferingFarmerID,
		OfferedItemID:       payload.Listing.OfferedItemID,
		OfferedItemQuantity: payload.Listing.OfferedItemQuantity,
		DesiredItems:        payload.Listing.DesiredItems,
		Status:              "open",
		ImageURL:            imageURL,
	}

	// Parse ExpiresAt if provided
	if payload.Listing.ExpiresAt != nil {
		// Assuming the database layer or consumer handles time parsing
		// For now, we'll skip parsing and let the DB layer handle it
		// In production, you'd want to parse the time string here
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
