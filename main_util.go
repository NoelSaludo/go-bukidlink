package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func retInternalServErr(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
