package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

func setupServer() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:username", getUserHandler)
	r.POST("/postuser", postUserHandler)

	return r
}

func setupDatabase() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST,
		PORT,
		USER,
		PASSWORD,
		DATABASE,
	)

	db, err := sql.Open("postgres", connStr)
	checkErr(err)

	return db
}

func postUserHandler(c *gin.Context) {
	var data User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       data.Id,
		"password": data.Password,
		"username": data.Username,
	})
}

func getUserHandler(c *gin.Context) {
	// TODO: delete later
	mockdb := []User{
		{Id: 1, Username: "JohnDoe", Password: "password123"},
		{Id: 2, Username: "DanielGaliego", Password: "passing"},
	}

	usernameP := c.Param("username")

	for _, user := range mockdb {
		if user.Username == usernameP {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
}

func main() {
	r := setupServer()

	r.Run("localhost:8080")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
