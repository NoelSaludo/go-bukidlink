package main

import (
	"bukidlink/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCartHandler(c *gin.Context) {
	userId := c.Param("user_id")
	cart, err := db.GetCartByUserID(userId)
	if err != nil {
		retInternalServErr(err, c)
		return
	}
	c.JSON(http.StatusOK, cart)
}

func addCartItemHandler(c *gin.Context) {
	var req AddToCartRequest
	if err := c.BindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	if err := db.AddItemToCart(req.CartID, req.ItemID, req.Quantity); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "item added to cart"})
}

func removeCartItemHandler(c *gin.Context) {
	cartItemId := c.Param("cart_item_id")
	if err := db.RemoveItemFromCart(cartItemId); err != nil {
		retInternalServErr(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "item removed from cart"})
}
