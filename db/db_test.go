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
	}

	generatedId := InsertUser(testUser)
	DeleteUser(generatedId)
}
