package db

import "time"

// TODO: add first_name, last_name, ContactNumber
type User struct {
	Id       string `json:"id"` // uuid
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type Item struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      int     `json:"amount"`
	CostPKilo   float64 `json:"costPKilo"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
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
	Id          string      `json:"id"`
	UserId      string      `json:"userid"`
	Status      string      `json:"status"`
	OrderDate   time.Time   `json:"order_date"`
	TotalPrice  float64     `json:"total_price"`
	Items       []OrderItem `json:"items"`
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
