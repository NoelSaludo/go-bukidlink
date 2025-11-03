package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bukidlink/db"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingRoute(t *testing.T) {
	server := setupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestPostUser(t *testing.T) {
	server := setupServer()

	data := db.User{
		Id:       uuid.New().String(),
		Username: "JohnDoe",
		Password: "password123",
		Email:    "JohnDoe" + uuid.NewString() + "@example.com",
		Details: db.UserDetail{
			Address:       "123 Main St",
			FirstName:     "John",
			LastName:      "Doe",
			ContactNumber: "1234567890",
			CreatedDate:   time.Now(),
		},
	}

	// decode image to base64
	imgPath := "resources/images/nanakusa-nazuna-icons.png"
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	// POST payload now expects an envelope: { "user": {...}, "profile_pic": {"base64":"","content_type":""} }
	payload := map[string]interface{}{
		"user": data,
		"profile_pic": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}

	jData, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		log.Fatalf("Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUser(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/user/JohnDoe", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// server may respond with an envelope: { "user": {...}, "profile_pic": {...} }
	// try to decode that envelope and extract the user
	var envelope struct {
		User       db.User         `json:"user"`
		ProfilePic json.RawMessage `json:"profile_pic"`
	}

	// First try to decode into envelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err == nil && envelope.User.Username != "" {
		assert.NotEmpty(t, envelope.User)
		return
	}

	// Fallback: maybe the response is the user object directly
	var user db.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.NotEmpty(t, user)
}

func Test100Items(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/item/0", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response is now array of envelopes: [{"item": {...}, "item_pic": {...}}, ...]
	var itemEnvelopes []struct {
		Item    db.Item         `json:"item"`
		ItemPic json.RawMessage `json:"item_pic"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &itemEnvelopes)
	require.NoError(t, err)
	assert.NotEmpty(t, itemEnvelopes)
}

func TestGetItemsByCategory(t *testing.T) {
	server := setupServer()

	categories := []string{"fruits", "vegetables", "grains", "livestock", "dairy"}
	for _, cat := range categories {
		req, _ := http.NewRequest(http.MethodGet, "/item/category/"+cat, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		// Response is now array of envelopes: [{"item": {...}, "item_pic": {...}}, ...]
		var itemEnvelopes []struct {
			Item    db.Item         `json:"item"`
			ItemPic json.RawMessage `json:"item_pic"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &itemEnvelopes)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, envelope := range itemEnvelopes {
			assert.Equal(t, cat, envelope.Item.Category)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, "/item/category/nothing", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetComments(t *testing.T) {
	s := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/review/a3e1b9f2-7d94-4d3a-9b4a-111111111111", nil)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var comments []db.Review
	if err := json.Unmarshal(w.Body.Bytes(), &comments); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, comments)
}

func TestGetItemById(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/item?id=a3e1b9f2-7d94-4d3a-9b4a-111111111111", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// Response may be envelope: {"item": {...}, "item_pic": {...}}
	var envelope struct {
		Item    db.Item         `json:"item"`
		ItemPic json.RawMessage `json:"item_pic"`
	}

	// Try to decode as envelope first
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err == nil && envelope.Item.Id != "" {
		assert.NotEmpty(t, envelope.Item)
		assert.NotEmpty(t, envelope.Item.Rating)
		return
	}

	// Fallback: maybe the response is the item object directly
	var item db.Item
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		log.Fatal(err)
	}

	assert.NotEmpty(t, item)
	assert.NotEmpty(t, item.Rating)
}

func TestPostItem(t *testing.T) {
	server := setupServer()

	newItem := db.Item{
		Name:        "Test Item",
		Description: "This is a test item",
		CostPKilo:   10.0,
		Category:    "dairy",
		Amount:      50,
	}

	// Sample image to encode
	imgPath := "resources/images/nanakusa-nazuna-icons.png"
	imgData, err := os.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	payload := map[string]interface{}{
		"item": newItem,
		"item_pic": map[string]string{
			"base64":       imgBase64,
			"content_type": "image/png",
		},
	}

	jData, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/item", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		log.Fatalf("Response Body: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrderAPIWorkflow(t *testing.T) {
	server := setupServer()
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803" // JohnDoe

	// 1. POST /order
	newOrder := db.Order{
		UserId:     userID,
		OrderDate:  time.Now(),
		Status:     "Packaging",
		TotalPrice: 5.0,
		Items: []db.OrderItem{
			{ItemId: "a3e1b9f2-7d94-4d3a-9b4a-111111111111", Quantity: 2, PriceAtPurchase: 1.0},
			{ItemId: "c9d3e8a1-55b2-4f66-a123-333333333333", Quantity: 1, PriceAtPurchase: 3.0},
		},
	}
	jData, _ := json.Marshal(newOrder)
	req, _ := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	orderID, ok := resp["order_id"]
	require.True(t, ok, "order_id not found in response")

	// 2. GET /order/:user_id
	req, _ = http.NewRequest(http.MethodGet, "/order/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var orders []db.Order
	json.Unmarshal(w.Body.Bytes(), &orders)
	assert.NotEmpty(t, orders)

	// 3. PATCH /order/status
	req, _ = http.NewRequest(http.MethodPatch, "/order/status?id="+orderID+"&status=Shipping", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GET /order/:user_id again to verify status update
	req, _ = http.NewRequest(http.MethodGet, "/order/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &orders)

	var updatedOrder db.Order
	for _, o := range orders {
		if o.Id == orderID {
			updatedOrder = o
			break
		}
	}
	assert.Equal(t, "Shipping", updatedOrder.Status)

	// 5. DELETE /order
	req, _ = http.NewRequest(http.MethodDelete, "/order?order_id="+orderID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. GET /order/:user_id to verify deletion
	req, _ = http.NewRequest(http.MethodGet, "/order/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &orders)

	var deletedOrder db.Order
	for _, o := range orders {
		if o.Id == orderID {
			deletedOrder = o
			break
		}
	}
	assert.Empty(t, deletedOrder.Id, "Deleted order should not be found")
}

func TestCartAPIWorkflow(t *testing.T) {
	server := setupServer()
	userID := "c6554794-849f-4338-87c5-6db2e2f76514" // DanielGaliego
	bananaID := "a3e1b9f2-7d94-4d3a-9b4a-111111111111"
	tomatoID := "b7f2c6d4-1aeb-4f5b-9c2b-222222222222"

	// Helper function to clear the cart
	clearCart := func() {
		req, _ := http.NewRequest(http.MethodGet, "/cart/"+userID, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		var cart db.Cart
		json.Unmarshal(w.Body.Bytes(), &cart)
		for _, item := range cart.Items {
			req, _ := http.NewRequest(http.MethodDelete, "/cart/item/"+item.Id, nil)
			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)
		}
	}

	// 1. Setup: Ensure the cart is empty before starting
	clearCart()

	// Defer cleanup for after the test runs
	defer clearCart()

	// 2. GET /cart/:user_id - This will create a cart if one doesn't exist
	req, _ := http.NewRequest(http.MethodGet, "/cart/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var cart db.Cart
	json.Unmarshal(w.Body.Bytes(), &cart)
	require.NotEmpty(t, cart.Id, "Cart ID should not be empty")
	assert.Empty(t, cart.Items, "Cart should be empty at the start")

	// 3. POST /cart/item - Add a banana
	addItemReq1 := AddToCartRequest{CartID: cart.Id, ItemID: bananaID, Quantity: 2}
	jData1, _ := json.Marshal(addItemReq1)
	req, _ = http.NewRequest(http.MethodPost, "/cart/item", bytes.NewBuffer(jData1))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GET /cart/:user_id again to verify one item
	req, _ = http.NewRequest(http.MethodGet, "/cart/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &cart)
	assert.Len(t, cart.Items, 1, "Cart should have 1 item type after adding banana")
	assert.Equal(t, 2, cart.Items[0].Quantity)
	assert.Equal(t, bananaID, cart.Items[0].ItemId)

	// 5. POST /cart/item - Add a tomato
	addItemReq2 := AddToCartRequest{CartID: cart.Id, ItemID: tomatoID, Quantity: 3}
	jData2, _ := json.Marshal(addItemReq2)
	req, _ = http.NewRequest(http.MethodPost, "/cart/item", bytes.NewBuffer(jData2))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. GET /cart/:user_id to verify two items
	req, _ = http.NewRequest(http.MethodGet, "/cart/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &cart)
	assert.Len(t, cart.Items, 2, "Cart should have 2 item types after adding tomato")

	// 7. DELETE /cart/item/:cart_item_id - Remove the banana
	var bananaCartItemID string
	for _, item := range cart.Items {
		if item.ItemId == bananaID {
			bananaCartItemID = item.Id
			break
		}
	}
	req, _ = http.NewRequest(http.MethodDelete, "/cart/item/"+bananaCartItemID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. GET /cart/:user_id to verify deletion
	req, _ = http.NewRequest(http.MethodGet, "/cart/"+userID, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &cart)
	assert.Len(t, cart.Items, 1, "Cart should have 1 item type after deleting banana")
	assert.Equal(t, tomatoID, cart.Items[0].ItemId, "The remaining item should be the tomato")
}

func TestGetUsersPosts(t *testing.T) {
	server := setupServer()
	userID := "d30869ec-fb97-46d8-85a3-82608c01f803" // JohnDoe

	req, _ := http.NewRequest(http.MethodGet, "/userpost/"+userID, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Response now includes enriched posts with base64 encoded images
	var enrichedPosts []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &enrichedPosts)
	require.NoError(t, err)
	assert.NotEmpty(t, enrichedPosts)

	// Verify that posts have the expected fields including image encoding
	for _, post := range enrichedPosts {
		assert.NotEmpty(t, post["id"])
		assert.NotEmpty(t, post["farmer_id"])
		assert.NotEmpty(t, post["content"])

		// If image_url exists, verify base64 and content_type are included
		if post["image_url"] != nil && post["image_url"] != "" {
			assert.NotEmpty(t, post["image_base64"], "image_base64 should be present when image_url exists")
			assert.NotEmpty(t, post["image_content_type"], "image_content_type should be present when image_url exists")
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

	// Response includes post data with optional base64 encoded image
	var post map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &post)
	require.NoError(t, err)
	assert.NotEmpty(t, post)

	// Verify that post has the expected fields
	assert.NotEmpty(t, post["id"])
	assert.NotEmpty(t, post["farmer_id"])
	assert.NotEmpty(t, post["content"])

	// If image_url exists, verify base64 and content_type are included
	if post["image_url"] != nil && post["image_url"] != "" {
		assert.NotEmpty(t, post["image_base64"], "image_base64 should be present when image_url exists")
		assert.NotEmpty(t, post["image_content_type"], "image_content_type should be present when image_url exists")
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

	var updatedPost map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updatedPost)
	assert.Equal(t, "Updated post content", updatedPost["content"])

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
	assert.Equal(t, "Final updated content", updatedPost["content"])
	assert.NotEmpty(t, updatedPost["image_base64"])

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
