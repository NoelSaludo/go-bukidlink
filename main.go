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

	commentG := r.Group("/comment")
	commentG.GET("/:itemId", getCommentByItemID)
	commentG.POST("/:productID")

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

func updateOrderStatusHandler(c *gin.Context) {
	var id string = c.Query("id")
	var status string = c.Query("status")

	if err := db.UpdateOrderStatus(id, status); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order status updated"})
}

func getUsersOrdersHandler(c *gin.Context) {
	userId := c.Param("user_id")

	orders, err := db.QueryUsersOrders(userId)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func postOrderHandler(c *gin.Context) {
	var order db.Order
	if err := c.BindJSON(&order); err != nil {
		retBadReqErr(err, c)
		return
	}

	orderId, err := db.InsertOrder(order)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "order created", "order_id": orderId})
}

func getCommentByItemID(c *gin.Context) {
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
