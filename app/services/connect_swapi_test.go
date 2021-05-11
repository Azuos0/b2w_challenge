package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetExistentPlanetNumberOfApperances(t *testing.T) {
	planetName := "tatooine"

	apperances, err := getPlanetNumberOfApperances(planetName)

	require.Greater(t, apperances, 0)
	require.Nil(t, err)
}

func TestGetNumberOfApperancesFindMultiplePlanets(t *testing.T) {
	planetName := "t"

	apperances, err := getPlanetNumberOfApperances(planetName)

	require.Equal(t, apperances, 0)
	require.Nil(t, err)
}
func TestGetNonExistentPlanetNumberOfApperances(t *testing.T) {
	planetName := "DARTH VADER PLANET"

	apperances, err := getPlanetNumberOfApperances(planetName)

	require.Equal(t, apperances, 0)
	require.Nil(t, err)
}
