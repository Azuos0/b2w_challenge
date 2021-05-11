package database_test

import (
	"os"
	"testing"

	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func getDbName() string {
	dbName := os.Getenv("MONGODB_TEST_DATABASE")

	//load env variables if they are not loaded yet,
	if dbName == "" {
		godotenv.Load("../../.env")
		_ = os.Getenv("MONGODB_TEST_DATABASE")
		dbName = os.Getenv("MONGODB_TEST_DATABASE")
	}

	return dbName
}

func TestDbConnect(t *testing.T) {
	dbName := getDbName()
	db, err := database.Connect(dbName)

	require.NotNil(t, db)
	require.Nil(t, err)
}

func TestDbGetCollection(t *testing.T) {
	dbName := getDbName()
	collectionName := "planets"

	db, _ := database.Connect(dbName)

	collection := database.GetCollection(db, collectionName)

	require.NotNil(t, collection)
}
