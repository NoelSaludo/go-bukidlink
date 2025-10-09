package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
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

func InsertUser(user User) int64 {
	query := `INSERT INTO "User" (username, password) VALUES ($1, $2) RETURNING id`
	var id int64

	err := db.QueryRow(query, user.Username, user.Password).Scan(&id)

	checkErr(err)

	return id
}

func DeleteUser(id int64) {
	query := `DELETE FROM "User" WHERE id=$1`

	_, err := db.Exec(query, id)

	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
