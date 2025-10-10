package main

import (
	"bukidlink/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	r.GET("/item/:block", get100ItemsHandler)

	r.POST("/postuser", postUserHandler)

	db.SetupDatabase()

	return r
}

func postUserHandler(c *gin.Context) {
	var data db.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int64
	id, err = db.InsertUser(data)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "created_id": id})
}

func getUserHandler(c *gin.Context) {
	temp := []db.User{}

	usernameP := c.Param("username")

	temp, err := db.QueryUsers(usernameP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, user := range temp {
		if user.Username == usernameP {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
}

func get100ItemsHandler(c *gin.Context) {

	blockP, err := strconv.Atoi(c.Param("block"))	
	if err != nil {	
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block parameter"})
		return
	}

	items, err := db.QueryAllItem100(blockP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
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
