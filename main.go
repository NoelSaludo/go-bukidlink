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
	orderG.GET("", getOrderHandler)
	orderG.POST("", postOrderHandler)
	orderG.DELETE("", deleteOrderHandler)

	db.SetupDatabase()

	return r
}

func getOrderHandler(c *gin.Context) {
	var id string = c.Query("id")

	order, err := db.QueryOrder(id)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, order)
}

func postOrderHandler(c *gin.Context) {
	var order db.Order
	if err := c.BindJSON(&order); err != nil {
		retBadReqErr(err, c)
		return
	}

	if err := db.InsertOrder(order); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order created"})
}

func deleteOrderHandler(c *gin.Context) {
	var id string = c.Query("id")

	if err := db.DeleteOrder(id); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order deleted"})
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
