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

func getItemByCategory(c *gin.Context) {
	cat := c.Param("category")

	fruits, err := db.QueryItembyCategory(cat)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, fruits)
}

func getItembyId(c *gin.Context) {
	itemid := c.Query("id")

	item, err := db.QueryItemByID(itemid)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, item)
}
