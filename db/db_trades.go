package db

import (
	"database/sql"
	"log"
)

// CreateTradeListing adds a new trade listing to the database.
func CreateTradeListing(listing TradeListing) error {
	query := `INSERT INTO "TradeListing" (id, offering_farmer_id, offered_item_id, offered_item_quantity, desired_items, status, expires_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(query, listing.ID, listing.OfferingFarmerID, listing.OfferedItemID, listing.OfferedItemQuantity, listing.DesiredItems, listing.Status, listing.ExpiresAt)
	if err != nil {
		log.Printf("Error creating trade listing: %v", err)
	}
	return err
}

// GetTradeListingByID retrieves a single trade listing from the database by its ID.
func GetTradeListingByID(id string) (*TradeListing, error) {
	query := `SELECT id, offering_farmer_id, offered_item_id, offered_item_quantity, desired_items, status, created_at, expires_at
              FROM "TradeListing" WHERE id = $1`
	row := db.QueryRow(query, id)
	var listing TradeListing
	err := row.Scan(&listing.ID, &listing.OfferingFarmerID, &listing.OfferedItemID, &listing.OfferedItemQuantity, &listing.DesiredItems, &listing.Status, &listing.CreatedAt, &listing.ExpiresAt)
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
