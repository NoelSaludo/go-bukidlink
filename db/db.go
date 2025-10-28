package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	HOST     string = os.Getenv("DBHOST")
	PORT     string = os.Getenv("DBPORT")
	USER     string = os.Getenv("DBUSER")
	PASSWORD string = os.Getenv("DBPASSWORD")
	DATABASE string = os.Getenv("DATABASE")
)

var db *sql.DB

func SetupDatabase() error {
	var err error

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST,
		PORT,
		USER,
		PASSWORD,
		DATABASE,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

func Ping() error {
	return db.Ping()
}

func QueryUsers(username string) ([]User, error) {
	var temp []User
	query := `SELECT id, username, password, email, address
				FROM "User" 
				WHERE username=$1`

	rows, err := db.Query(query, username)
	if err != nil {
		return temp, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Address)

		if err != nil {
			return temp, err
		}

		temp = append(temp, user)
	}

	return temp, err
}

func InsertUser(user User) error {
	query := `
	INSERT INTO "User" (
	id, username, password, email, address)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query,
		user.Id,
		user.Username,
		user.Password,
		user.Email,
		user.Address)

	if err != nil {
		return err
	}

	return err
}

func DeleteUser(id string) error {
	query := `DELETE FROM "User" WHERE id=$1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func QueryItemByID(id string) (Item, error) {

	query := `SELECT i.id as itemid,
			i.name,
			i.description,
			i.costpkilo,
			i.category,
			i.amount,
			AVG(c.rating) AS rating
			FROM "Item" AS i
			JOIN "Review" AS C ON c.itemid = i.id
			WHERE i.id=$1
			GROUP BY i.id, i.name
			`

	row := db.QueryRow(query, id)

	var item Item
	err := row.Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.CostPKilo,
		&item.Category,
		&item.Amount, &item.Rating)

	if err != nil {
		return Item{}, err
	}

	return item, err
}

func QueryAllItem100(block int) ([]Item, error) {
	var items []Item
	// select 100 items with offset
	query := `SELECT id,
			name,
			description,
			costpkilo,
			category,
			amount 
			FROM "Item" LIMIT 100 OFFSET $1`
	rows, err := db.Query(query, block*100)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	items = getItemsFromRow(rows, items)
	return items, err
}

// category only accepts certian string
// accepted string
// fruits, vegetables, grains, livestock, dairy, others
func QueryItembyCategory(category string) ([]Item, error) {
	var items []Item
	query := `SELECT id,
				name,
				description,
				costpkilo,
				category,
				amount 
				FROM public."Item" WHERE category=$1`

	rows, err := db.Query(query, category)
	if err != nil {
		return nil, err
	}

	items = getItemsFromRow(rows, items)
	return items, nil
}

func QueryReviewsOnItem(itemid string) ([]Review, error) {
	var reviews []Review

	query := `SELECT id, userid, itemid, content, rating
			FROM public."Review" WHERE itemid=$1`

	rows, err := db.Query(query, itemid)
	if err != nil {
		return nil, err
	}

	reviews, err = getReviewFromRow(rows, reviews)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func QueryUsersItem(userid string) ([]Item, error) {
	var items []Item
	query := `
	SELECT i.id   AS itemid,
       i.name,
       i.description,
       i.costpkilo,
       i.category,
       i.amount
		FROM "User" u
		JOIN "UsersItem" ui ON u.id = ui.userid
		JOIN "Item" i ON i.id = ui.itemid
		WHERE u.id = $1; `

	rows, err := db.Query(query, userid)
	if err != nil {
		return nil, err
	}

	items = getItemsFromRow(rows, items)

	return items, nil
}

// must be used only with items and with proper order
func getItemsFromRow(rows *sql.Rows, items []Item) []Item {
	for rows.Next() {
		var item Item

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.CostPKilo,
			&item.Category,
			&item.Amount)

		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return items
}

func getReviewFromRow(rows *sql.Rows, reviews []Review) ([]Review, error) {
	for rows.Next() {
		var review Review

		err := rows.Scan(
			&review.Id,
			&review.UserId,
			&review.ItemId,
			&review.Content,
			&review.Rating)

		if err != nil {
			return reviews, err
		}

		reviews = append(reviews, review)
	}
	return reviews, nil
}
