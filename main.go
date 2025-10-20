package main

import (
	"bukidlink/db"
	"log"
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
	itemGroup.GET("/fruit", getFruitsHandler)
	itemGroup.GET("/vegetables")
	itemGroup.GET("/grains")
	itemGroup.GET("/livestock")
	itemGroup.GET("/diary")
	itemGroup.GET("/other")

	userGroup := r.Group("/user")
	userGroup.GET("/:username", getUserHandler)
	userGroup.POST("", postUserHandler)

	db.SetupDatabase()

	return r
}

func main() {
	r := setupServer()

	r.Run("localhost:8080")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
