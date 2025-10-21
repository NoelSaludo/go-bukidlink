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
			Id:       2,
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
		Username: "DanielGalliego",
		Password: "VeryC00lP@ss",
		Email:    "DanielGalliego@example.com",
	}

	generatedId, err := InsertUser(testUser)
	require.NoError(t, err)
	assert.NotEmpty(t, generatedId)
	err = DeleteUser(generatedId)
	require.NoError(t, err)
}

func TestGetItem(t *testing.T) {
	_ = SetupDatabase()

	var item Item
	item, err := QueryItemByID(3)

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
