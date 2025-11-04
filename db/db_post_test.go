package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPostsByBlock(t *testing.T) {
	_ = SetupDatabase()

	posts, err := GetPostsByBlock(1)
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.LessOrEqual(t, len(posts), 100)

	// Verify that comments are included with user info
	for _, post := range posts {
		assert.NotNil(t, post.Comments, "Comments should be initialized (not nil)")
		// If post has comments, verify they include user info
		for _, comment := range post.Comments {
			assert.NotEmpty(t, comment.ID, "Comment should have ID")
			assert.NotEmpty(t, comment.UserID, "Comment should have UserID")
			assert.NotEmpty(t, comment.Username, "Comment should have Username")
			// ProfilePicUrl can be empty but field should exist
			assert.NotNil(t, comment.ProfilePicUrl, "Comment should have ProfilePicUrl field")
		}
	}
}

func TestGetPostByID(t *testing.T) {
	_ = SetupDatabase()

	post, err := GetPostByID("11111111-1111-1111-1111-111111111111")
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", post.ID)
	assert.NotNil(t, post.Comments, "Comments should be initialized (not nil)")

	// Verify comments include user info
	for _, comment := range post.Comments {
		assert.NotEmpty(t, comment.ID, "Comment should have ID")
		assert.NotEmpty(t, comment.UserID, "Comment should have UserID")
		assert.NotEmpty(t, comment.Username, "Comment should have Username")
		// ProfilePicUrl can be empty but field should exist (not checking NotEmpty)
	}
}

func TestUpdatePost(t *testing.T) {
	_ = SetupDatabase()

	err := UpdatePost("11111111-1111-1111-1111-111111111111", "Updated content", "resources/images/updated-image-url.jpg")
	assert.NoError(t, err)
}

func TestGetPostsByUser(t *testing.T) {
	_ = SetupDatabase()

	posts, err := GetPostsByUser("d30869ec-fb97-46d8-85a3-82608c01f803")
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	// Verify that comments are included for each post with user info
	for _, post := range posts {
		assert.NotNil(t, post.Comments, "Comments should be initialized (not nil)")
		// If post has comments, verify they include user info
		for _, comment := range post.Comments {
			assert.NotEmpty(t, comment.ID, "Comment should have ID")
			assert.NotEmpty(t, comment.UserID, "Comment should have UserID")
			assert.NotEmpty(t, comment.Username, "Comment should have Username")
			// ProfilePicUrl can be empty but field should exist
		}
	}
}

func TestInsertAndDeletePost(t *testing.T) {
	_ = SetupDatabase()

	// Insert a new post
	newPost := Post{
		ID:        "77777777-7777-7777-7777-777777777777",
		FarmerID:  "d30869ec-fb97-46d8-85a3-82608c01f803",
		FarmID:    nil,
		Content:   "This is a test post for insert and delete functionality.",
		ImageURL:  nil,
		CreatedAt: time.Now(),
	}

	err := InsertPost(newPost)
	assert.NoError(t, err)

	// Verify the post was inserted
	insertedPost, err := GetPostByID(newPost.ID)
	assert.NoError(t, err)
	assert.NotNil(t, insertedPost)
	assert.Equal(t, newPost.ID, insertedPost.ID)
	assert.NotNil(t, insertedPost.Comments, "Comments should be initialized (not nil)")
	assert.Empty(t, insertedPost.Comments, "New post should have no comments")

	// Delete the post
	err = DeletePost(newPost.ID)
	assert.NoError(t, err)

	// Verify the post was deleted
	deletedPost, err := GetPostByID(newPost.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedPost)
}
