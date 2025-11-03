package main

import (
	"bukidlink/db"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func getUserPostsHandler(c *gin.Context) {
	userID := c.Param("user_id")

	posts, err := db.GetPostsByUser(userID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Enhance posts with base64 encoded images
	var enrichedPosts []map[string]interface{}
	for _, post := range posts {
		postData := map[string]interface{}{
			"id":         post.ID,
			"farmer_id":  post.FarmerID,
			"farm_id":    post.FarmID,
			"content":    post.Content,
			"image_url":  post.ImageURL,
			"created_at": post.CreatedAt,
		}

		// If image_url exists, read and encode the image
		if post.ImageURL != nil && *post.ImageURL != "" {
			imageData, err := os.ReadFile(*post.ImageURL)
			if err == nil {
				// Encode to base64
				base64Image := base64.StdEncoding.EncodeToString(imageData)
				postData["image_base64"] = base64Image

				// Determine content type from file extension
				ext := strings.ToLower(filepath.Ext(*post.ImageURL))
				contentType := "image/jpeg" // default
				switch ext {
				case ".png":
					contentType = "image/png"
				case ".jpg", ".jpeg":
					contentType = "image/jpeg"
				case ".gif":
					contentType = "image/gif"
				case ".webp":
					contentType = "image/webp"
				}
				postData["image_content_type"] = contentType
			}
		}

		enrichedPosts = append(enrichedPosts, postData)
	}

	c.JSON(http.StatusOK, enrichedPosts)
}
func getUserPostHandler(c *gin.Context) {

}
func postUserPostHandler(c *gin.Context) {

}
func updateUserPostHandler(c *gin.Context) {

}
func deleteUserPostHandler(c *gin.Context) {

}
