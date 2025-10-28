package db

import "time"

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

type Order struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userid"`
	ItemId      string    `json:"itemid"`
	Amount      int       `json:"amount"`
	Status      string    `json:"status"`
	CreatedDate time.Time `json:"created_date"`
}
