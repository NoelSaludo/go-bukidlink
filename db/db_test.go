package db

import (
	"testing"

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

	expectedUsers := []User{
		{
			Id:       `d30869ec-fb97-46d8-85a3-82608c01f803`,
			Username: "JohnDoe",
			Password: "P@ssw0rd",
			Email:    "JohnDoe@example.com",
		},
	}

	resultUsers, err := QueryUsers("JohnDoe")
	require.NoError(t, err)
	assert.ElementsMatch(t, expectedUsers, resultUsers)
}

func TestInsertAndDelete(t *testing.T) {
	_ = SetupDatabase()

	testUser := User{
		Id:       "01d85ea5-0c1f-457c-b1f5-04f4e48b54b6",
		Username: "DanielGalliego",
		Password: "VeryC00lP@ss",
		Email:    "DanielGalliego@example.com",
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

	items, err = QueryItembyCategory("vegetables")

	require.NoError(t, err)
	assert.NotEmpty(t, items)

	items, err = QueryItembyCategory("grains")

	require.NoError(t, err)
	assert.NotEmpty(t, items)

	items, err = QueryItembyCategory("livestock")

	require.NoError(t, err)
	assert.NotEmpty(t, items)

	items, err = QueryItembyCategory("dairy")

	require.NoError(t, err)
	assert.NotEmpty(t, items)

	items, err = QueryItembyCategory("others")

	require.NoError(t, err)
	assert.NotEmpty(t, items)
}

func TestQueryCommentsEdgeCases(t *testing.T) {
	_ = SetupDatabase()

	// Case: item with no comments (assuming id 9999 does not exist)
	comments, err := QueryCommentsOnItem("3c270f60-8934-4692-8922-011f25dda434")
	require.NoError(t, err)
	assert.IsType(t, []Comment{}, comments)
	assert.Empty(t, comments)

	// Case: validate comment fields for a known item (3)
	comments, err = QueryCommentsOnItem("a3e1b9f2-7d94-4d3a-9b4a-111111111111")
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
