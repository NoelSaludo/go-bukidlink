package main

import (
	"bukidlink/db"
	"errors"
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

// UpdateCartItemRequest represents the PATCH payload to update a cart item's quantity
type UpdateCartItemRequest struct {
	CartItemID string `json:"cart_item_id"`
	Quantity   int    `json:"quantity"`
}

func updateCartItemHandler(c *gin.Context) {
	var req UpdateCartItemRequest
	if err := c.BindJSON(&req); err != nil {
		retBadReqErr(err, c)
		return
	}

	if req.CartItemID == "" {
		retBadReqErr(errors.New("cart_item_id is required"), c)
		return
	}

	if err := db.UpdateCartItemQuantity(req.CartItemID, req.Quantity); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "quantity updated"})
}
