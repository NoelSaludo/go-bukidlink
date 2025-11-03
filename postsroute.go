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

// encodeImageToBase64 reads an image file and returns base64 encoded data and content type
func encodeImageToBase64(imageURL *string) (string, string, error) {
	if imageURL == nil || *imageURL == "" {
		return "", "", nil
	}

	imageData, err := os.ReadFile(*imageURL)
	if err != nil {
		return "", "", err
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Determine content type from file extension
	ext := strings.ToLower(filepath.Ext(*imageURL))
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

	return base64Image, contentType, nil
}

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

		// Encode image if available
		if base64Image, contentType, err := encodeImageToBase64(post.ImageURL); err == nil && base64Image != "" {
			postData["image_base64"] = base64Image
			postData["image_content_type"] = contentType
		}

		enrichedPosts = append(enrichedPosts, postData)
	}

	c.JSON(http.StatusOK, enrichedPosts)
}

func getUserPostHandler(c *gin.Context) {
	postID := c.Param("post_id")

	post, err := db.GetPostByID(postID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Enhance post with base64 encoded image
	postData := map[string]interface{}{
		"id":         post.ID,
		"farmer_id":  post.FarmerID,
		"farm_id":    post.FarmID,
		"content":    post.Content,
		"image_url":  post.ImageURL,
		"created_at": post.CreatedAt,
	}

	// Encode image if available
	if base64Image, contentType, err := encodeImageToBase64(post.ImageURL); err == nil && base64Image != "" {
		postData["image_base64"] = base64Image
		postData["image_content_type"] = contentType
	}

	c.JSON(http.StatusOK, postData)
}

func postUserPostHandler(c *gin.Context) {

}
func updateUserPostHandler(c *gin.Context) {

}
func deleteUserPostHandler(c *gin.Context) {

}
