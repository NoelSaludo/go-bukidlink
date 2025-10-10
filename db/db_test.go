package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestDatabaseConnection(t *testing.T) {
	err := SetupDatabase()

	assert.Equal(t, nil, err)
	assert.Equal(t, nil, db.Ping())
}

func TestQueryUsers(t *testing.T) {
	_ = SetupDatabase()

	expectedUsers := []User {
		{
			Id: 2,
			Username: "JohnDoe",
			Password: "P@ssw0rd",
			Email: "JohnDoe@example.com",
		},
	}

	resultUsers := QueryUsers("JohnDoe")
	assert.ElementsMatch(t, expectedUsers, resultUsers)
}

func TestInsertAndDelete(t *testing.T) {
	_ = SetupDatabase()

	testUser := User {
		Username: "DanielGalliego",
		Password: "VeryC00lP@ss",
		Email: "DanielGalliego@example.com",
	}

	generatedId, err := InsertUser(testUser)
	assert.Equal(t, nil, err)
	err = DeleteUser(generatedId)
	assert.Equal(t, nil, err)
}

func TestGetItem(t *testing.T) {
	_ = SetupDatabase()

	expectedItem := Item{
		Id: 1,
		Name: "Carrot",
		Description: "A long, orange, and beautiful Carrot",
		Amount: 10,
	}

	var item Item
	item, err := QueryItemByID(1)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedItem, item)
}
