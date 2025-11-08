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
	itemGroup.POST("", postItemHandler)

	userGroup := r.Group("/user")
	userGroup.GET("/:username", getUserHandler)
	userGroup.POST("", postUserHandler)

	reviewG := r.Group("/review")
	reviewG.GET("/:itemId", getReviewByItemID)

	orderG := r.Group("/order")
	orderG.GET("/:user_id", getUsersOrdersHandler) // Changed from getOrderHandler
	orderG.POST("", postOrderHandler)
	orderG.PATCH("/status", updateOrderStatusHandler)
	orderG.DELETE("", deleteOrderHandler)

	cartG := r.Group("/cart")
	cartG.GET("/:user_id", getCartHandler)
	cartG.POST("/item", addCartItemHandler)
	cartG.DELETE("/item/:cart_item_id", removeCartItemHandler)

	userPostG := r.Group("/userpost")
	userPostG.GET("/:user_id", getUserPostsHandler)
	userPostG.GET("/post/:post_id", getUserPostHandler)
	userPostG.POST("", postUserPostHandler)
	userPostG.PATCH("/:post_id", updateUserPostHandler)
	userPostG.DELETE("/:post_id", deleteUserPostHandler)

	tradeG := r.Group("/trade")
	tradeG.GET("/batch", getTradeListingsBatchHandler)
	tradeG.GET("", getTradeListingByIDHandler) // Query param: ?id=uuid
	tradeG.POST("", postTradeListingHandler)
	tradeG.PATCH("/:id", updateTradeListingStatusHandler) // Path param: /trade/:id

	// Trade bid routes
	bidG := r.Group("/bid")
	bidG.GET("", getTradeBidByIDHandler) // Query param: ?id=uuid
	bidG.POST("", postTradeBidHandler)
	bidG.PUT("/:id", updateTradeBidHandler)                // Path param: /bid/:id
	bidG.PATCH("/:id/status", updateTradeBidStatusHandler) // Path param: /bid/:id/status
	bidG.DELETE("/:id", deleteTradeBidHandler)             // Path param: /bid/:id

	bidG.GET("/farmer/:farmer_id", getTradeBidsByFarmerHandler) // Path param: /bid/farmer/:farmer_id

	// Payment routes
	balanceG := r.Group("/balance")
	balanceG.GET("/:user_id", getUserBalanceHandler)
	balanceG.POST("", createUserBalanceHandler)

	paymentG := r.Group("/payment")
	paymentG.POST("/deposit", processDepositHandler)
	paymentG.POST("/withdrawal", processWithdrawalHandler)
	paymentG.POST("/order", processOrderPaymentHandler)
	paymentG.POST("/refund", processRefundHandler)
	paymentG.GET("/transaction/:transaction_id", getPaymentTransactionHandler)
	paymentG.GET("/transactions/:user_id", getUserTransactionsHandler)

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
