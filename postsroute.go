package main

import (
	"bukidlink/db"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// decodeAndSaveImage decodes base64 image data and saves it to a file
func decodeAndSaveImage(base64Data string, contentType string, filename string) (string, error) {
	// Decode base64 string
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", err
	}

	// Determine file extension from content type
	ext := ".jpg" // default
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}

	// Create file path
	filePath := filepath.Join("resources", "images", filename+ext)

	// Write file to disk
	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
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
			"comments":   post.Comments,
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
		"comments":   post.Comments,
	}

	// Encode image if available
	if base64Image, contentType, err := encodeImageToBase64(post.ImageURL); err == nil && base64Image != "" {
		postData["image_base64"] = base64Image
		postData["image_content_type"] = contentType
	} else {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, postData)
}

func postUserPostHandler(c *gin.Context) {
	var payload struct {
		Post struct {
			FarmerID string  `json:"farmer_id"`
			FarmID   *string `json:"farm_id"`
			Content  string  `json:"content"`
		} `json:"post"`
		PostImage struct {
			Base64      string `json:"base64"`
			ContentType string `json:"content_type"`
		} `json:"post_image"`
	}

	if err := c.BindJSON(&payload); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Generate new post ID
	postID := uuid.New().String()

	// Handle image if provided
	var imageURL *string
	if payload.PostImage.Base64 != "" && payload.PostImage.ContentType != "" {
		// Save image with post ID as filename
		filePath, err := decodeAndSaveImage(payload.PostImage.Base64, payload.PostImage.ContentType, postID+"_post")
		if err != nil {
			retConflictErr(err, c)
			return
		}
		imageURL = &filePath
	}

	// Create post object
	newPost := db.Post{
		ID:        postID,
		FarmerID:  payload.Post.FarmerID,
		FarmID:    payload.Post.FarmID,
		Content:   payload.Post.Content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}

	// Insert post into database
	if err := db.InsertPost(newPost); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post_id": postID,
		"message": "Post created successfully",
	})
}

func updateUserPostHandler(c *gin.Context) {
	postID := c.Param("post_id")

	var payload struct {
		Content   *string `json:"content"`
		PostImage *struct {
			Base64      string `json:"base64"`
			ContentType string `json:"content_type"`
		} `json:"post_image"`
	}

	if err := c.BindJSON(&payload); err != nil {
		retBadReqErr(err, c)
		return
	}

	// Get the existing post
	existingPost, err := db.GetPostByID(postID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Update content if provided
	content := existingPost.Content
	if payload.Content != nil {
		content = *payload.Content
	}

	// Handle image update if provided
	imageURL := ""
	if existingPost.ImageURL != nil {
		imageURL = *existingPost.ImageURL
	}

	if payload.PostImage != nil && payload.PostImage.Base64 != "" && payload.PostImage.ContentType != "" {
		// Delete old image if exists
		if existingPost.ImageURL != nil && *existingPost.ImageURL != "" {
			os.Remove(*existingPost.ImageURL)
		}

		// Save new image with post ID as filename
		filePath, err := decodeAndSaveImage(payload.PostImage.Base64, payload.PostImage.ContentType, postID+"_post")
		if err != nil {
			retInternalServErr(err, c)
			return
		}
		imageURL = filePath
	}

	// Update post in database
	if err := db.UpdatePost(postID, content, imageURL); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
	})
}

func deleteUserPostHandler(c *gin.Context) {
	postID := c.Param("post_id")

	// Get the existing post to retrieve image path
	existingPost, err := db.GetPostByID(postID)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	// Delete the post from database
	if err := db.DeletePost(postID); err != nil {
		retInternalServErr(err, c)
		return
	}

	// Delete the associated image file if exists
	if existingPost.ImageURL != nil && *existingPost.ImageURL != "" {
		os.Remove(*existingPost.ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}
