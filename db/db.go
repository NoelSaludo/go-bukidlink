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
	checkErr(err)

	return nil
}

func Ping() error {
	return db.Ping()
}

func QueryUsers(username string) []User {
	var temp []User
	query := `SELECT * FROM "User" WHERE username=$1`

	rows, err := db.Query(query, username)
	checkErr(err)

	for rows.Next() {
		var id int
		var username string
		var password string
		var email string
		err := rows.Scan(&id, &username, &password, &email)

		if err != nil {
			log.Fatal(err)
		}

		temp = append(temp, User{
			Id:       id,
			Username: username,
			Password: password,
			Email:    email,
		})
	}

	return temp
}

func InsertUser(user User) (int64, error) {
	query := `INSERT INTO "User" (username, password, email) VALUES ($1, $2, $3) RETURNING id`
	var id int64

	err := db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&id)

	return id, err
}

func DeleteUser(id int64) error {
	query := `DELETE FROM "User" WHERE id=$1`

	_, err := db.Exec(query, id)

	return err
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func QueryItemByID(id int) (Item, error) {

	query := `SELECT * FROM "Item" WHERE id=$1`

	row := db.QueryRow(query, id)

	var itemId int
	var name string
	var description string
	var amount int

	err := row.Scan(&itemId, &name, &description, &amount)

	if err != nil {
		return Item{}, err
	}

	resItem := Item{
		Id:          itemId,
		Name:        name,
		Description: description,
		Amount:      amount,
	}

	return resItem, err
}

func QueryAllItem100(block int) ([]Item, error) {
	var item []Item
	// select 100 items with offset
	query := `SELECT * FROM "Item" LIMIT 100 OFFSET $1`
	rows, err := db.Query(query, block*100)
	if err != nil {
		return item, err
	}

	defer rows.Close()

	for rows.Next() {
		var itemId int
		var name string
		var description string
		var amount int
		err := rows.Scan(&itemId, &name, &description, &amount)
		if err != nil {
			return item, err
		}
		item = append(item, Item{
			Id:          itemId,
			Name:        name,
			Description: description,
			Amount:      amount,
		})
	}

	return item, err
}
