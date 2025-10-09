package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	HOST     string = os.Getenv("DBHOST")
	PORT     string = os.Getenv("DBPORT")
	USER     string = os.Getenv("DBUSER")
	PASSWORD string = os.Getenv("DBPASSWORD")
	DATABASE string = os.Getenv("DATABASE")
)

var db *sql.DB

func SetupDatabase() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST,
		PORT,
		USER,
		PASSWORD,
		DATABASE,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

func Ping() error {
	return db.Ping()
}

func QueryUsers(username string) []User {
	var temp []User
	query := `SELECT * FROM "User" WHERE username=$1`

	rows, err := db.Query(query, username)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var username string
		var password string
		err := rows.Scan(&id, &username, &password)

		if err != nil {
			log.Fatal(err)
		}

		temp = append(temp, User{
			Id:       id,
			Username: username,
			Password: password,
		})
	}

	return temp
}
