package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func retInternalServErr(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func retBadReqErr(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

type AddToCartRequest struct {
	CartID   string `json:"cart_id"`
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}
