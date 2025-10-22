package db

import (
	"database/sql"
	"fmt"
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
	query := `SELECT * FROM "User" WHERE username=$1`

	rows, err := db.Query(query, username)
	if err != nil {
		return temp, err
	}

	for rows.Next() {
		var id int
		var username string
		var password string
		var email string
		err := rows.Scan(&id, &username, &password, &email)

		if err != nil {
			return temp, err
		}

		temp = append(temp, User{
			Id:       id,
			Username: username,
			Password: password,
			Email:    email,
		})
	}

	return temp, err
}

func InsertUser(user User) (int64, error) {
	query := `
	INSERT INTO "User" (
	username, password, email)
	VALUES ($1, $2, $3) 
	RETURNING id`
	var id int64

	err := db.QueryRow(query,
		user.Username,
		user.Password,
		user.Email).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func DeleteUser(id int64) error {
	query := `DELETE FROM "User" WHERE id=$1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func QueryItemByID(id int) (Item, error) {

	query := `SELECT * FROM "Item" WHERE id=$1`

	row := db.QueryRow(query, id)

	var item Item
	err := row.Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.CostPKilo,
		&item.Category,
		&item.Amount,
		&item.Rating)

	if err != nil {
		return Item{}, err
	}

	return item, err
}

func QueryAllItem100(block int) ([]Item, error) {
	var items []Item
	// select 100 items with offset
	query := `SELECT * FROM "Item" LIMIT 100 OFFSET $1`
	rows, err := db.Query(query, block*100)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	items, err = getItemsFromRow(rows, items)
	if err != nil {
		return items, err
	}

	return items, err
}

// category only accepts certian string
// accepted string
// fruits, vegetables, grains, livestock, dairy, others
func QueryItembyCategory(category string) ([]Item, error) {
	var items []Item
	query := `SELECT * FROM public."Item" WHERE category=$1`

	rows, err := db.Query(query, category)
	if err != nil {
		return nil, err
	}

	items, err = getItemsFromRow(rows, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func QueryCommentsOnItem(itemid int) ([]Comment, error) {
	var comments []Comment

	query := `SELECT id, userid, itemid, content, rating
			FROM public."Comment" WHERE itemid=$1`

	rows, err := db.Query(query, itemid)
	if err != nil {
		return nil, err
	}

	comments, err = getCommentFromRow(rows, comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func getItemsFromRow(rows *sql.Rows, items []Item) ([]Item, error) {
	for rows.Next() {
		var item Item

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.CostPKilo,
			&item.Category,
			&item.Amount,
			&item.Rating)

		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func getCommentFromRow(rows *sql.Rows, comments []Comment) ([]Comment, error) {
	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.Id,
			&comment.UserId,
			&comment.ItemId,
			&comment.Content,
			&comment.Rating)

		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}
