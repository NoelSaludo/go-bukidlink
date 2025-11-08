package db

import "time"

// User contains basic, public-facing user information
type User struct {
	Id       string `json:"id"` // uuid
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"` // This should be a hash, not plaintext
	// ProfilePicPath stores the URL/path to the user's profile picture as stored in the DB.
	// DB column name: profile_pic_url
	ProfilePicPath string     `json:"profile_pic_url"`
	Details        UserDetail `json:"details"`
}

// UserDetail contains private user information
type UserDetail struct {
	Address       string    `json:"address"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ContactNumber string    `json:"contact_number"`
	CreatedDate   time.Time `json:"created_date"`
}

type Item struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      int     `json:"amount"`
	CostPKilo   float64 `json:"costPKilo"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
	ImgPath     string  `json:"img_path"`
}

type Review struct {
	Id      string  `json:"id"`
	ItemId  string  `json:"itemid"`
	UserId  string  `json:"userid"`
	Content string  `json:"content"`
	Rating  float64 `json:"rating"`
}

// Order now contains a list of items for the order
type Order struct {
	Id         string      `json:"id"`
	UserId     string      `json:"userid"`
	Status     string      `json:"status"`
	OrderDate  time.Time   `json:"order_date"`
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `json:"items"`
}

// OrderItem represents a single item within an order
type OrderItem struct {
	Id              string  `json:"id"`
	OrderId         string  `json:"order_id"`
	ItemId          string  `json:"item_id"`
	Quantity        int     `json:"quantity"`
	PriceAtPurchase float64 `json:"price_at_purchase"`
}

// Cart now contains a list of items in the cart
type Cart struct {
	Id         string     `json:"id"`
	UserId     string     `json:"userid"`
	GrandTotal float64    `json:"grand_total"`
	CreatedAt  time.Time  `json:"created_at"`
	Items      []CartItem `json:"items"`
}

// CartItem represents a single item within a cart
type CartItem struct {
	Id       string `json:"id"`
	CartId   string `json:"cart_id"`
	ItemId   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// TradeListing represents an item offered for trade by a farmer.
type TradeListing struct {
	ID                  string     `json:"id"`
	OfferingFarmerID    string     `json:"offering_farmer_id"`
	OfferedItemID       string     `json:"offered_item_id"`
	OfferedItemQuantity float64    `json:"offered_item_quantity"`
	DesiredItems        string     `json:"desired_items"`
	Status              string     `json:"status"`
	ImageURL            *string    `json:"image_url,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
}

// TradeBid represents a bid made by a farmer on a TradeListing.
type TradeBid struct {
	ID              string    `json:"id"`
	TradeListingID  string    `json:"trade_listing_id"`
	BiddingFarmerID string    `json:"bidding_farmer_id"`
	BidItemID       string    `json:"bid_item_id"`
	BidItemQuantity float64   `json:"bid_item_quantity"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// UserBalance stores the current balance for each user.
type UserBalance struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Find the PaymentTransaction struct and update these fields to pointers:

type PaymentTransaction struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	OrderID         *string    `json:"order_id,omitempty"` // Changed to pointer (nullable)
	TransactionType string     `json:"transaction_type"`
	Amount          float64    `json:"amount"`
	BalanceBefore   float64    `json:"balance_before"`
	BalanceAfter    float64    `json:"balance_after"`
	Status          string     `json:"status"`
	PaymentMethod   string     `json:"payment_method"`
	ReferenceNumber *string    `json:"reference_number,omitempty"` // Changed to pointer (nullable)
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"` // Changed to pointer (nullable)
}
