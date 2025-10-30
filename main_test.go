package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bukidlink/db"

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
		Id:       "01d85ea5-0c1f-457c-b1f5-04f4e48b54b6",
		Username: "JohnDoe",
		Password: "password123",
		Email:    "JohnDoe@example.com",
	}

	jData, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetUser(t *testing.T) {
	server := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/user/JohnDoe", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
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
}

func TestGetItemsByCategory(t *testing.T) {
	server := setupServer()

	categories := []string{"fruits", "vegetables", "grains", "livestock", "dairy"}
	for _, cat := range categories {
		req, _ := http.NewRequest(http.MethodGet, "/item/category/"+cat, nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		var items []db.Item
		err := json.Unmarshal(w.Body.Bytes(), &items)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		for _, i := range items {
			assert.Equal(t, cat, i.Category)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, "/item/category/nothing", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetComments(t *testing.T) {
	s := setupServer()

	req, _ := http.NewRequest(http.MethodGet, "/comment/a3e1b9f2-7d94-4d3a-9b4a-111111111111", nil)
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

	var item db.Item
	if err := json.Unmarshal(w.Body.Bytes(), &item); err != nil {
		log.Fatal(err)
	}

	assert.NotEmpty(t, item)
	assert.NotEmpty(t, item.Rating)
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

	// 3. POST /order/status
	req, _ = http.NewRequest(http.MethodPost, "/order/status?id="+orderID+"&status=Shipping", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
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
