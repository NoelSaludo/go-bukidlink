package main

import (
	"log"
	"net/http"
	"bukidlink/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

func postUserHandler(c *gin.Context) {
	var data db.User
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
	temp := []db.User{ }

	usernameP := c.Param("username")

	for _, user := range temp {
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
