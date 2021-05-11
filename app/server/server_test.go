package server_test

import (
	"os"
	"testing"

	"github.com/Azuos0/b2w_challenge/app/server"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestInitializeApp(t *testing.T) {
	app := server.App{}
	var err error

	if os.Getenv("MONGODB_TEST_DATABASE") == "" {
		err = godotenv.Load("../../.env")
	}

	require.Nil(t, err)

	uri := os.Getenv("MONGODB_TEST_DATABASE")
	app.InitializeApp(uri)

	require.NotNil(t, app.DB)
	require.NotNil(t, app.Router)
}
