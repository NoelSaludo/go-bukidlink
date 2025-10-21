package main

import (
	"bukidlink/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func get100ItemsHandler(c *gin.Context) {

	blockP, err := strconv.Atoi(c.Param("block"))
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	items, err := db.QueryAllItem100(blockP)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, items)
}

func getFruitsHandler(c *gin.Context) {
	fruits, err := db.QueryFruits()
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, fruits)
}

func getVegetablesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "vegetables endpoint not implemented yet"})
}

func getGrainsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "grains endpoint not implemented yet"})
}

func getLivestockHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "livestock endpoint not implemented yet"})
}

func getDiaryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "diary endpoint not implemented yet"})
}

func getOtherHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "other endpoint not implemented yet"})
}
