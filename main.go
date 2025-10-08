package main

import (
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

var (

)

func setupServer() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}

func main() {
	r := setupServer()

	r.Run("localhost:8080")
}
