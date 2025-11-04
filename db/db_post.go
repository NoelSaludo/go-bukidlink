package db

import "time"

// Post represents a post in the Posts table
type Post struct {
	ID        string    `json:"id"`
	FarmerID  string    `json:"farmer_id"`
	FarmID    *string   `json:"farm_id"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	ID            string    `json:"id"`
	PostID        string    `json:"post_id"`
	UserID        string    `json:"user_id"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	Username      string    `json:"username"`
	ProfilePicUrl string    `json:"profile_pic_url"`
}

// Function to get all 100 posts by block
func GetPostsByBlock(block int) ([]Post, error) {
	offset := (block - 1) * 100
	rows, err := db.Query(`SELECT id, farmer_id, farm_id, content, image_url, created_at FROM "Posts" ORDER BY created_at DESC LIMIT 100 OFFSET $1`, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.FarmerID, &post.FarmID, &post.Content, &post.ImageURL, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch comments for this post
		comments, err := getCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, post)
	}
	return posts, nil
}

// Function to get a specific post
func GetPostByID(postID string) (*Post, error) {
	var post Post
	err := db.QueryRow(`SELECT id, farmer_id, farm_id, content, image_url, created_at FROM "Posts" WHERE id = $1`, postID).
		Scan(&post.ID, &post.FarmerID, &post.FarmID, &post.Content, &post.ImageURL, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch comments for this post
	comments, err := getCommentsByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return &post, nil
}

// Function to delete a post
func DeletePost(postID string) error {
	_, err := db.Exec(`DELETE FROM "Posts" WHERE id = $1`, postID)
	return err
}

// Function to update a post
func UpdatePost(postID string, content string, imageURL string) error {
	_, err := db.Exec(`UPDATE "Posts" SET content = $1, image_url = $2 WHERE id = $3`, content, imageURL, postID)
	return err
}

// Function to get all posts by a user
func GetPostsByUser(userID string) ([]Post, error) {
	rows, err := db.Query(`SELECT id, farmer_id, farm_id, content, image_url, created_at FROM "Posts" WHERE farmer_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.FarmerID, &post.FarmID, &post.Content, &post.ImageURL, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch comments for this post
		comments, err := getCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, post)
	}
	return posts, nil
}

// InsertPost inserts a new post into the Posts table
func InsertPost(post Post) error {
	query := `INSERT INTO "Posts" (id, farmer_id, farm_id, content, image_url, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, post.ID, post.FarmerID, post.FarmID, post.Content, post.ImageURL, post.CreatedAt)
	return err
}

// getCommentsByPostID retrieves all comments for a specific post with user info
func getCommentsByPostID(postID string) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.profile_pic_url 
		FROM "Comment" c
		JOIN "User" u ON c.user_id = u.id
		WHERE c.post_id = $1 
		ORDER BY c.created_at ASC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username, &comment.ProfilePicUrl); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Return empty slice instead of nil if no comments
	if comments == nil {
		comments = []Comment{}
	}

	return comments, nil
}
