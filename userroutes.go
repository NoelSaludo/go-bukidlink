package main

import (
	"bukidlink/db"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func postUserHandler(c *gin.Context) {
	var data db.User
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = db.InsertUser(data)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func getUserHandler(c *gin.Context) {
	var temp db.User

	usernameP := c.Param("username")

	temp, err := db.QueryUser(usernameP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if temp.Username != "" {
		// Attempt to fetch and embed profile picture (non-fatal).
		if temp.ProfilePicPath != "" {
			const maxImageSize = 1 << 20 // 1 MiB limit
			var data []byte
			var contentType string

			// Treat ProfilePicPath as a local filesystem path. Resolve relative
			// paths from the process working directory and read the file.
			fp := temp.ProfilePicPath
			if !filepath.IsAbs(fp) {
				fp = filepath.Join(".", fp)
			}
			f, err := os.Open(fp)
			if err == nil {
				defer f.Close()
				data, err = io.ReadAll(io.LimitReader(f, maxImageSize+1))
				if err == nil && int64(len(data)) <= maxImageSize {
					contentType = http.DetectContentType(data)
				} else {
					data = nil
				}
			}

			if len(data) > 0 {
				b64 := base64.StdEncoding.EncodeToString(data)
				c.JSON(http.StatusOK, gin.H{
					"user": temp,
					"profile_pic": gin.H{
						"content_type": contentType,
						"base64":       b64,
					},
				})
				return
			}
			// if fetch failed or image too large, fallthrough and return user only
		}

		// Default: return user only
		c.JSON(http.StatusOK, temp)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
}
