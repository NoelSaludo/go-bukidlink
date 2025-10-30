package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	err := SetupDatabase()

	require.NoError(t, err)
	assert.Equal(t, nil, db.Ping())
}

func TestQueryUsers(t *testing.T) {
	_ = SetupDatabase()

	resultUsers, err := QueryUsers("JohnDoe")
	require.NoError(t, err)
	assert.NotEmpty(t, resultUsers)
}

func TestInsertAndDelete(t *testing.T) {
	_ = SetupDatabase()

	testUser := User{
		Id:       "01d85ea5-0c1f-457c-b1f5-04f4e48b54b6",
		Username: "DanielGalliego",
		Password: "VeryC00lP@ss",
		Email:    "DanielGalliego@example.com",
		Address:  "In our hearts",
	}

	err := InsertUser(testUser)
	require.NoError(t, err)
	err = DeleteUser("01d85ea5-0c1f-457c-b1f5-04f4e48b54b6")
	require.NoError(t, err)

}

func TestGetItem(t *testing.T) {
	_ = SetupDatabase()

	var item Item
	item, err := QueryItemByID("a3e1b9f2-7d94-4d3a-9b4a-111111111111")

	require.NoError(t, err)
	assert.NotEmpty(t, item)
}

func TestGetAllItems(t *testing.T) {
	_ = SetupDatabase()

	var items []Item
	items, err := QueryAllItem100(0)

	require.NoError(t, err)
	assert.NotEmpty(t, items)
}

func TestQueryItemsbyCategory(t *testing.T) {
	_ = SetupDatabase()

	var items []Item
	items, err := QueryItembyCategory("fruits")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "fruits", i.Category)
	}

	items, err = QueryItembyCategory("vegetables")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "vegetables", i.Category)
	}

	items, err = QueryItembyCategory("grains")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "grains", i.Category)
	}

	items, err = QueryItembyCategory("livestock")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "livestock", i.Category)
	}

	items, err = QueryItembyCategory("dairy")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "dairy", i.Category)
	}

	items, err = QueryItembyCategory("others")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
	for _, i := range items {
		assert.Equal(t, "others", i.Category)
	}
}

func TestQueryCommentsEdgeCases(t *testing.T) {
	_ = SetupDatabase()

	// Case: item with no comments (assuming id 9999 does not exist)
	comments, err := QueryReviewsOnItem("3c270f60-8934-4692-8922-011f25dda434")
	require.NoError(t, err)
	assert.IsType(t, []Review{}, comments)
	assert.Empty(t, comments)

	// Case: validate comment fields for a known item (3)
	comments, err = QueryReviewsOnItem("a3e1b9f2-7d94-4d3a-9b4a-111111111111")
	require.NoError(t, err)
	require.NotEmpty(t, comments)

	for _, c := range comments {
		// IDs should be positive
		assert.IsType(t, "", c.Id)
		assert.IsType(t, "", c.UserId)
		assert.IsType(t, "", c.ItemId)
		// Content should be non-empty
		assert.NotEmpty(t, c.Content)
		// Rating should be within expected bounds (0-5)
		assert.GreaterOrEqual(t, int(c.Rating), 0)
		assert.LessOrEqual(t, int(c.Rating), 5)
	}
}

func TestUsersItemQuery(t *testing.T) {
	_ = SetupDatabase()

	items, err := QueryUsersItem("d30869ec-fb97-46d8-85a3-82608c01f803")
	require.NoError(t, err)
	assert.NotEmpty(t, items)
}

func TestOrderWorkflow(t *testing.T) {
	_ = SetupDatabase()

	// 1. Create a new order
	newOrder := Order{
		UserId:     "d30869ec-fb97-46d8-85a3-82608c01f803", // JohnDoe
		OrderDate:  time.Now(),
		Status:     "Packaging",
		TotalPrice: 5.0, // 2 * 1.0 (Banana) + 1 * 3.0 (Rice)
		Items: []OrderItem{
			{
				ItemId:          "a3e1b9f2-7d94-4d3a-9b4a-111111111111", // Banana
				Quantity:        2,
				PriceAtPurchase: 1.0,
			},
			{
				ItemId:          "c9d3e8a1-55b2-4f66-a123-333333333333", // Rice
				Quantity:        1,
				PriceAtPurchase: 3.0,
			},
		},
	}

	orderId, err := InsertOrder(newOrder)
	require.NoError(t, err)
	require.NotEmpty(t, orderId)

	// 2. Query the user's orders and verify the new order is there
	orders, err := QueryUsersOrders("d30869ec-fb97-46d8-85a3-82608c01f803")
	require.NoError(t, err)
	assert.NotEmpty(t, orders)

	var foundOrder Order
	for _, o := range orders {
		if o.Id == orderId {
			foundOrder = o
			break
		}
	}
	require.NotEmpty(t, foundOrder.Id, "Could not find inserted order")

	assert.Equal(t, "Packaging", foundOrder.Status)
	assert.Equal(t, 5.0, foundOrder.TotalPrice)
	assert.Len(t, foundOrder.Items, 2)

	// 3. Update the order status
	err = UpdateOrderStatus(orderId, "Shipping")
	require.NoError(t, err)

	// 4. Query again to verify the status update
	orders, err = QueryUsersOrders("d30869ec-fb97-46d8-85a3-82608c01f803")
	require.NoError(t, err)

	for _, o := range orders {
		if o.Id == orderId {
			foundOrder = o
			break
		}
	}
	require.NotEmpty(t, foundOrder.Id, "Could not find inserted order after update")
	assert.Equal(t, "Shipping", foundOrder.Status)
}

func TestCartWorkflow(t *testing.T) {
	_ = SetupDatabase()

	userID := "d30869ec-fb97-46d8-85a3-82608c01f803" // JohnDoe
	bananaID := "a3e1b9f2-7d94-4d3a-9b4a-111111111111"
	tomatoID := "b7f2c6d4-1aeb-4f5b-9c2b-222222222222"

	// 1. Get cart for a user (should create one if it doesn't exist)
	cart, err := GetCartByUserID(userID)
	require.NoError(t, err)
	assert.Equal(t, userID, cart.UserId)

	// 2. Add items to the cart
	err = AddItemToCart(cart.Id, bananaID, 2) // Add 2 bananas
	require.NoError(t, err)

	err = AddItemToCart(cart.Id, tomatoID, 3) // Add 3 tomatoes
	require.NoError(t, err)

	err = AddItemToCart(cart.Id, bananaID, 1) // Add 1 more banana
	require.NoError(t, err)

	// 3. Verify cart contents
	cart, err = GetCartByUserID(userID)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 2)

	var bananaItem CartItem
	var tomatoItem CartItem
	for _, item := range cart.Items {
		if item.ItemId == bananaID {
			bananaItem = item
		} else if item.ItemId == tomatoID {
			tomatoItem = item
		}
	}

	assert.Equal(t, 3, bananaItem.Quantity, "Should have 3 bananas in total")
	assert.Equal(t, 3, tomatoItem.Quantity, "Should have 3 tomatoes")

	// 4. Remove an item from the cart
	err = RemoveItemFromCart(tomatoItem.Id)
	require.NoError(t, err)

	// 5. Verify item removal
	cart, err = GetCartByUserID(userID)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, bananaID, cart.Items[0].ItemId)
}
