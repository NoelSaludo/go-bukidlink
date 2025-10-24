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

	db.SetupDatabase()

	return r
}
func getCommentByItemID(c *gin.Context) {
	var itemid string = c.Param("itemId")

	var comments []db.Comment
	if comms, err := db.QueryCommentsOnItem(itemid); err != nil {
		retInternalServErr(err, c)
		return
	} else {
		comments = comms
	}

	c.JSON(http.StatusOK, comments)

}

func main() {
	r := setupServer()

	r.Run("localhost:8080")
}
