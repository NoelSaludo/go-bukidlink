package db

import (
	"database/sql"
	"log"
)

// CreateTradeListing adds a new trade listing to the database.
func CreateTradeListing(listing TradeListing) error {
	query := `INSERT INTO "TradeListing" (id, offering_farmer_id, offered_item_id, offered_item_quantity, desired_items, status, image_url, expires_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, listing.ID, listing.OfferingFarmerID, listing.OfferedItemID, listing.OfferedItemQuantity, listing.DesiredItems, listing.Status, listing.ImageURL, listing.ExpiresAt)
	if err != nil {
		log.Printf("Error creating trade listing: %v", err)
	}
	return err
}

// GetTradeListingByID retrieves a single trade listing from the database by its ID.
func GetTradeListingByID(id string) (*TradeListing, error) {
	query := `SELECT id, offering_farmer_id, offered_item_id, offered_item_quantity, desired_items, status, image_url, created_at, expires_at
			  FROM "TradeListing" WHERE id = $1`
	row := db.QueryRow(query, id)
	var listing TradeListing
	err := row.Scan(&listing.ID, &listing.OfferingFarmerID, &listing.OfferedItemID, &listing.OfferedItemQuantity, &listing.DesiredItems, &listing.Status, &listing.ImageURL, &listing.CreatedAt, &listing.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		log.Printf("Error scanning trade listing: %v", err)
		return nil, err
	}
	return &listing, nil
}

// CreateTradeBid adds a new trade bid to the database.
func CreateTradeBid(bid TradeBid) error {
	query := `INSERT INTO "TradeBid" (id, trade_listing_id, bidding_farmer_id, bid_item_id, bid_item_quantity, status)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, bid.ID, bid.TradeListingID, bid.BiddingFarmerID, bid.BidItemID, bid.BidItemQuantity, bid.Status)
	if err != nil {
		log.Printf("Error creating trade bid: %v", err)
	}
	return err
}

// GetTradeBidsByListingID retrieves all bids for a given trade listing.
func GetTradeBidsByListingID(listingID string) ([]TradeBid, error) {
	query := `SELECT id, trade_listing_id, bidding_farmer_id, bid_item_id, bid_item_quantity, status, created_at
              FROM "TradeBid" WHERE trade_listing_id = $1`
	rows, err := db.Query(query, listingID)
	if err != nil {
		log.Printf("Error querying trade bids: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bids []TradeBid = []TradeBid{}
	for rows.Next() {
		var bid TradeBid
		err := rows.Scan(&bid.ID, &bid.TradeListingID, &bid.BiddingFarmerID, &bid.BidItemID, &bid.BidItemQuantity, &bid.Status, &bid.CreatedAt)
		if err != nil {
			log.Printf("Error scanning trade bid: %v", err)
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

// UpdateTradeListingStatus updates the status of a trade listing.
func UpdateTradeListingStatus(id string, status string) error {
	query := `UPDATE "TradeListing" SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, id)
	if err != nil {
		log.Printf("Error updating trade listing status: %v", err)
	}
	return err
}

// UpdateTradeBidStatus updates the status of a trade bid.
func UpdateTradeBidStatus(id string, status string) error {
	query := `UPDATE "TradeBid" SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, status, id)
	if err != nil {
		log.Printf("Error updating trade bid status: %v", err)
	}
	return err
}

// GetTradeBidByID retrieves a single trade bid from the database by its ID.
func GetTradeBidByID(id string) (*TradeBid, error) {
	query := `SELECT id, trade_listing_id, bidding_farmer_id, bid_item_id, bid_item_quantity, status, created_at
              FROM "TradeBid" WHERE id = $1`
	row := db.QueryRow(query, id)
	var bid TradeBid
	err := row.Scan(&bid.ID, &bid.TradeListingID, &bid.BiddingFarmerID, &bid.BidItemID, &bid.BidItemQuantity, &bid.Status, &bid.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		log.Printf("Error scanning trade bid: %v", err)
		return nil, err
	}
	return &bid, nil
}

// GetTradeBidsByFarmer retrieves all bids made by a specific farmer.
func GetTradeBidsByFarmer(farmerID string) ([]TradeBid, error) {
	query := `SELECT id, trade_listing_id, bidding_farmer_id, bid_item_id, bid_item_quantity, status, created_at
              FROM "TradeBid" WHERE bidding_farmer_id = $1
              ORDER BY created_at DESC`
	rows, err := db.Query(query, farmerID)
	if err != nil {
		log.Printf("Error querying trade bids by farmer: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bids []TradeBid = []TradeBid{}
	for rows.Next() {
		var bid TradeBid
		err := rows.Scan(&bid.ID, &bid.TradeListingID, &bid.BiddingFarmerID, &bid.BidItemID, &bid.BidItemQuantity, &bid.Status, &bid.CreatedAt)
		if err != nil {
			log.Printf("Error scanning trade bid: %v", err)
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

// UpdateTradeBid updates a trade bid's item and quantity (only allowed for pending bids).
func UpdateTradeBid(id string, bidItemID string, bidItemQuantity float64) error {
	query := `UPDATE "TradeBid" 
              SET bid_item_id = $1, bid_item_quantity = $2 
              WHERE id = $3 AND status = 'pending'`
	result, err := db.Exec(query, bidItemID, bidItemQuantity, id)
	if err != nil {
		log.Printf("Error updating trade bid: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Bid not found or not in pending status
	}

	return nil
}

// DeleteTradeBid deletes a trade bid from the database (only allowed for pending bids).
func DeleteTradeBid(id string) error {
	query := `DELETE FROM "TradeBid" WHERE id = $1 AND status = 'pending'`
	result, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting trade bid: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Bid not found or not in pending status
	}

	return nil
}

// GetTradeListingsByBlock retrieves a batch of trade listings (100 per block).
func GetTradeListingsByBlock(block int) ([]TradeListing, error) {
	query := `SELECT id,
				offering_farmer_id,
				offered_item_id,
				offered_item_quantity,
				desired_items,
				status,
			image_url,
			created_at,
			expires_at
              FROM "TradeListing"
              ORDER BY created_at DESC
              LIMIT 100 OFFSET $1`
	rows, err := db.Query(query, block*100)
	if err != nil {
		log.Printf("Error querying trade listings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var listings []TradeListing = []TradeListing{}
	for rows.Next() {
		var listing TradeListing
		err := rows.Scan(&listing.ID, &listing.OfferingFarmerID, &listing.OfferedItemID, &listing.OfferedItemQuantity, &listing.DesiredItems, &listing.Status, &listing.ImageURL, &listing.CreatedAt, &listing.ExpiresAt)
		if err != nil {
			log.Printf("Error scanning trade listing: %v", err)
			return nil, err
		}
		listings = append(listings, listing)
	}
	return listings, nil
}
