package models_test

import (
	"testing"

	"github.com/Azuos0/b2w_challenge/app/models"
	"github.com/stretchr/testify/require"
)

func TestPlanetValidation(t *testing.T) {
	validPlanet := models.Planet{}
	invalidPlanet := models.Planet{}

	validPlanet.Name = "New Planet"
	validPlanet.Climate = "New Climate"
	validPlanet.Terrain = "new terrain"

	invalidPlanet.Climate = "New Climate"
	invalidPlanet.Terrain = "new terrain"

	err1 := validPlanet.Validate()
	err2 := invalidPlanet.Validate()

	require.Nil(t, err1)
	require.Error(t, err2)
}
