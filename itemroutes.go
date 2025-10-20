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
