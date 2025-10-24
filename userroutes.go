package main

import (
	"bukidlink/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postUserHandler(c *gin.Context) {
	var data db.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.InsertUser(data)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
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
