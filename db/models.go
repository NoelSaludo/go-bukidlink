package db

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Item struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      int     `json:"amount"`
	CostPKilo   float64 `json:"costPKilo"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
}

type Comment struct {
	Id      int     `json:"id"`
	ItemId  int     `json:"itemid"`
	UserId  int     `json:"userid"`
	Content string  `json:"content"`
	Rating  float64 `json:"rating"`
}
