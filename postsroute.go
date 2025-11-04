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

// EnrichedComment represents a comment with base64-encoded profile picture data
type EnrichedComment struct {
	ID                    string    `json:"id"`
	PostID                string    `json:"post_id"`
	UserID                string    `json:"user_id"`
	Content               string    `json:"content"`
	CreatedAt             time.Time `json:"created_at"`
	Username              string    `json:"username"`
	ProfilePicUrl         string    `json:"profile_pic_url"`
	ProfilePicBase64      string    `json:"profile_pic_base64,omitempty"`
	ProfilePicContentType string    `json:"profile_pic_content_type,omitempty"`
}

// EnrichedPost represents a post with base64-encoded image data and enriched comments
type EnrichedPost struct {
	ID               string            `json:"id"`
	FarmerID         string            `json:"farmer_id"`
	FarmID           *string           `json:"farm_id"`
	Content          string            `json:"content"`
	ImageURL         *string           `json:"image_url"`
	CreatedAt        time.Time         `json:"created_at"`
	Comments         []EnrichedComment `json:"comments"`
	ImageBase64      string            `json:"image_base64,omitempty"`
	ImageContentType string            `json:"image_content_type,omitempty"`
}

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
	var enrichedPosts []EnrichedPost
	for _, post := range posts {
		// Initialize empty slice for enriched comments (ensures JSON returns [] not null)
		enrichedComments := make([]EnrichedComment, 0)

		// Enrich all comments with base64-encoded profile pictures
		for _, comment := range post.Comments {
			enrichedComment := EnrichedComment{
				ID:            comment.ID,
				PostID:        comment.PostID,
				UserID:        comment.UserID,
				Content:       comment.Content,
				CreatedAt:     comment.CreatedAt,
				Username:      comment.Username,
				ProfilePicUrl: comment.ProfilePicUrl,
			}

			// Encode profile picture if available
			if comment.ProfilePicUrl != "" {
				if base64Image, contentType, err := encodeImageToBase64(&comment.ProfilePicUrl); err == nil && base64Image != "" {
					enrichedComment.ProfilePicBase64 = base64Image
					enrichedComment.ProfilePicContentType = contentType
				}
			}

			enrichedComments = append(enrichedComments, enrichedComment)
		}

		// Build post with enriched comments
		postData := EnrichedPost{
			ID:        post.ID,
			FarmerID:  post.FarmerID,
			FarmID:    post.FarmID,
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			CreatedAt: post.CreatedAt,
			Comments:  enrichedComments,
		}

		// Encode post image if available
		if post.ImageURL != nil {
			if base64Image, contentType, err := encodeImageToBase64(post.ImageURL); err == nil && base64Image != "" {
				postData.ImageBase64 = base64Image
				postData.ImageContentType = contentType
			}
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

	// Initialize empty slice for enriched comments (ensures JSON returns [] not null)
	enrichedComments := make([]EnrichedComment, 0)

	// Enrich all comments with base64-encoded profile pictures
	for _, comment := range post.Comments {
		enrichedComment := EnrichedComment{
			ID:            comment.ID,
			PostID:        comment.PostID,
			UserID:        comment.UserID,
			Content:       comment.Content,
			CreatedAt:     comment.CreatedAt,
			Username:      comment.Username,
			ProfilePicUrl: comment.ProfilePicUrl,
		}

		// Encode profile picture if available
		if comment.ProfilePicUrl != "" {
			if base64Image, contentType, err := encodeImageToBase64(&comment.ProfilePicUrl); err == nil && base64Image != "" {
				enrichedComment.ProfilePicBase64 = base64Image
				enrichedComment.ProfilePicContentType = contentType
			}
		}

		enrichedComments = append(enrichedComments, enrichedComment)
	}

	// Build post with enriched comments
	postData := EnrichedPost{
		ID:        post.ID,
		FarmerID:  post.FarmerID,
		FarmID:    post.FarmID,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		Comments:  enrichedComments,
	}

	// Encode post image if available
	if post.ImageURL != nil {
		if base64Image, contentType, err := encodeImageToBase64(post.ImageURL); err == nil && base64Image != "" {
			postData.ImageBase64 = base64Image
			postData.ImageContentType = contentType
		}
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
