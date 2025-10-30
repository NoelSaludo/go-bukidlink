package main

import (
	"bukidlink/db"
	"net/http"

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

	itemGroup := r.Group("/item")
	itemGroup.GET("/:block", get100ItemsHandler)
	itemGroup.GET("/category/:category", getItemByCategory)
	itemGroup.GET("", getItembyId)

	userGroup := r.Group("/user")
	userGroup.GET("/:username", getUserHandler)
	userGroup.POST("", postUserHandler)

	reviewG := r.Group("/review")
	reviewG.GET("/:itemId", getReviewByItemID)

	orderG := r.Group("/order")
	orderG.GET("/:user_id", getUsersOrdersHandler) // Changed from getOrderHandler
	orderG.POST("", postOrderHandler)
	orderG.POST("/status", updateOrderStatusHandler)

	cartG := r.Group("/cart")
	cartG.GET("/:user_id", getCartHandler)
	cartG.POST("/item", addCartItemHandler)
	cartG.DELETE("/item/:cart_item_id", removeCartItemHandler)

	db.SetupDatabase()

	return r
}

func getReviewByItemID(c *gin.Context) {
	var itemid string = c.Param("itemId")

	var reviews []db.Review
	if comms, err := db.QueryReviewsOnItem(itemid); err != nil {
		retInternalServErr(err, c)
		return
	} else {
		reviews = comms
	}

	c.JSON(http.StatusOK, reviews)

}

func main() {
	r := setupServer()

	r.Run("localhost:8080")
}
