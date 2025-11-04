package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bukidlink/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUsersPosts(t *testing.T) {
	server := setupServer()
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803" // JohnDoe

	req, _ := http.NewRequest(http.MethodGet, "/userpost/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response now includes enriched posts with base64 encoded images
	var enrichedPosts []EnrichedPost
	err := json.Unmarshal(w.Body.Bytes(), &enrichedPosts)
	require.NoError(t, err)
	assert.NotEmpty(t, enrichedPosts)

	// Verify that posts have the expected fields including image encoding
	for _, post := range enrichedPosts {
		assert.NotEmpty(t, post.ID)
		assert.NotEmpty(t, post.FarmerID)
		assert.NotEmpty(t, post.Content)

		fmt.Printf("Post ID: %s, Content: %s, Comments: %d\n", post.ID, post.Content, len(post.Comments))
		// Verify comments field exists (can be empty array)
		assert.NotNil(t, post.Comments, "comments field should be present")

		// Verify each comment has user info and optional profile pic encoding (only if comments exist)
		for _, comment := range post.Comments {
			assert.NotEmpty(t, comment.ID, "comment should have id")
			assert.NotEmpty(t, comment.UserID, "comment should have user_id")
			assert.NotEmpty(t, comment.Username, "comment should have username")
			// ProfilePicUrl field always exists (it's a string, not pointer)

			// If profile_pic_url exists and is not empty, verify base64 encoding
			if comment.ProfilePicUrl != "" {
				// Profile pic base64 should be present
				if comment.ProfilePicBase64 != "" {
					assert.NotEmpty(t, comment.ProfilePicBase64, "profile_pic_base64 should be present when profile_pic_url exists")
					assert.NotEmpty(t, comment.ProfilePicContentType, "profile_pic_content_type should be present when profile pic is encoded")
				}
			}
		}

		// If image_url exists, verify base64 and content_type are included
		if post.ImageURL != nil && *post.ImageURL != "" {
			assert.NotEmpty(t, post.ImageBase64, "image_base64 should be present when image_url exists")
			assert.NotEmpty(t, post.ImageContentType, "image_content_type should be present when image_url exists")
		}
	}
}

func TestGetUserPost(t *testing.T) {
	server := setupServer()

	// Use a known post ID from your test data
	// You may need to adjust this ID based on your test database
	postID := "11111111-1111-1111-1111-111111111111"

	req, _ := http.NewRequest(http.MethodGet, "/userpost/post/"+postID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	if w.Code != http.StatusOK {
		log.Fatalf("Response Body: %s", w.Body.String())
	}

	// Response includes post data with optional base64 encoded image
	var post EnrichedPost
	err := json.Unmarshal(w.Body.Bytes(), &post)
	require.NoError(t, err)
	assert.NotEmpty(t, post)

	// Verify that post has the expected fields
	assert.NotEmpty(t, post.ID)
	assert.NotEmpty(t, post.FarmerID)
	assert.NotEmpty(t, post.Content)

	// Verify comments field exists (can be empty array)
	assert.NotNil(t, post.Comments, "comments field should be present")

	// Verify each comment has user info and optional profile pic encoding (only if comments exist)
	for _, comment := range post.Comments {
		assert.NotEmpty(t, comment.ID, "comment should have id")
		assert.NotEmpty(t, comment.UserID, "comment should have user_id")
		assert.NotEmpty(t, comment.Username, "comment should have username")
		// ProfilePicUrl field always exists (it's a string, not pointer)

		// If profile_pic_url exists and is not empty, verify base64 encoding
		if comment.ProfilePicUrl != "" {
			// Profile pic base64 should be present
			if comment.ProfilePicBase64 != "" {
				assert.NotEmpty(t, comment.ProfilePicBase64, "profile_pic_base64 should be present when profile_pic_url exists")
				assert.NotEmpty(t, comment.ProfilePicContentType, "profile_pic_content_type should be present when profile pic is encoded")
			}
		}
	}

	// If image_url exists, verify base64 and content_type are included
	if post.ImageURL != nil && *post.ImageURL != "" {
		assert.NotEmpty(t, post.ImageBase64, "image_base64 should be present when image_url exists")
		assert.NotEmpty(t, post.ImageContentType, "image_content_type should be present when image_url exists")
	}
}

func TestPostUserPost(t *testing.T) {
	server := setupServer()

	farmid := "11111111-aaaa-aaaa-aaaa-111111111111" // Sunny Fields
	newPost := db.Post{
		FarmerID: "d30869ec-fb97-46d8-85a3-82608c01f803", // JohnDoe
		FarmID:   &farmid,
		Content:  "This is a test post from JohnDoe.",
	}

	// Sample image to encode
	imgPath := "resources/images/nanakusa-nazuna-icons.png"
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	payload := map[string]interface{}{
		"post": newPost,
		"post_image": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}

	jData, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/userpost", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		log.Fatalf("Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUpdateUserPost(t *testing.T) {
	server := setupServer()

	// First, create a post to update
	farmid := "11111111-aaaa-aaaa-aaaa-111111111111" // Sunny Fields
	newPost := db.Post{
		FarmerID: "d30869ec-fb97-46d8-85a3-82608c01f803", // JohnDoe
		FarmID:   &farmid,
		Content:  "Original post content",
	}

	imgPath := "resources/images/nanakusa-nazuna-icons.png"
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		t.Skipf("Test image not found: %v", err)
		return
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	createPayload := map[string]interface{}{
		"post": newPost,
		"post_image": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}

	jData, _ := json.Marshal(createPayload)
	req, _ := http.NewRequest(http.MethodPost, "/userpost", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &createResp)
	postID := createResp["post_id"]
	require.NotEmpty(t, postID)

	// Test 1: Update content only
	updatePayload1 := map[string]interface{}{
		"content": "Updated post content",
	}
	jData1, _ := json.Marshal(updatePayload1)
	req, _ = http.NewRequest(http.MethodPatch, "/userpost/"+postID, bytes.NewBuffer(jData1))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the update
	req, _ = http.NewRequest(http.MethodGet, "/userpost/post/"+postID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedPost EnrichedPost
	json.Unmarshal(w.Body.Bytes(), &updatedPost)
	assert.Equal(t, "Updated post content", updatedPost.Content)

	// Verify comments field exists in updated post (can be empty array)
	assert.NotNil(t, updatedPost.Comments, "comments field should be present in updated post")

	// Test 2: Update image only
	updatePayload2 := map[string]interface{}{
		"post_image": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}
	jData2, _ := json.Marshal(updatePayload2)
	req, _ = http.NewRequest(http.MethodPatch, "/userpost/"+postID, bytes.NewBuffer(jData2))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test 3: Update both content and image
	updatePayload3 := map[string]interface{}{
		"content": "Final updated content",
		"post_image": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}
	jData3, _ := json.Marshal(updatePayload3)
	req, _ = http.NewRequest(http.MethodPatch, "/userpost/"+postID, bytes.NewBuffer(jData3))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the final update
	req, _ = http.NewRequest(http.MethodGet, "/userpost/post/"+postID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &updatedPost)
	assert.Equal(t, "Final updated content", updatedPost.Content)
	assert.NotEmpty(t, updatedPost.ImageBase64)

	// Cleanup: Delete the test post
	req, _ = http.NewRequest(http.MethodDelete, "/userpost/"+postID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the post was deleted
	req, _ = http.NewRequest(http.MethodGet, "/userpost/post/"+postID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Deleted post should not be found")
}
