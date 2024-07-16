package profiles

import (
	"testing"

	"github.com/chadeldridge/cuttle/tests"
	"github.com/stretchr/testify/require"
)

func TestTilesDefaultTile(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	require.Empty(tile.Name, "name was not empty")
	require.Equal(SmallestTileSize*DefaultSizeMultiplier, tile.DisplaySize, "displaySize did not match")
	require.Empty(tile.Tests, "tests was not empty")
	require.False(tile.AllMustPass, "allMustPass returned true")
	require.False(tile.InParallel, "inParallel returned true")
}

func TestTilesNewTile(t *testing.T) {
	require := require.New(t)
	name := "Ping Test"
	test := tests.Test{Name: name, MustSucceed: true}

	t.Run("valid", func(t *testing.T) {
		tile := NewTile(name, test)
		require.Equal(name, tile.Name, "name did not match")
		require.Equal(SmallestTileSize*DefaultSizeMultiplier, tile.DisplaySize, "displaySize did not match")
		require.Len(tile.Tests, 1, "tests did not have 1 test")
		require.False(tile.AllMustPass, "allMustPass returned true")
		require.False(tile.InParallel, "inParallel returned true")
	})
	// TODO: Add test cases for invalid name after validation is added to NewTile.
}

func TestTilesSetName(t *testing.T) {
	require := require.New(t)
	tile := NewTile("OldName")
	name := "NewName"

	t.Run("valid name", func(t *testing.T) {
		err := tile.SetName(name)
		require.NoError(err, "SetName() returned an error: %s", err)
		require.Equal(name, tile.Name, "name did not match")
	})

	// INCOMPLETE: Change to require.Error after SetName checks for valid name.
	/*
		t.Run("invalid name", func(t *testing.T) {
			err = tile.SetName("invalid name")
			require.Error(err, "SetName() did not return an error")
			require.Equal(name, tile.name, "name did not match")
		})
	*/

	t.Run("empty name", func(t *testing.T) {
		err := tile.SetName("")
		require.Error(err, "SetName() did not return an error")
		require.Equal(name, tile.Name, "name did not match")
	})
}

func TestTilesSetSize(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	t.Run("default", func(t *testing.T) {
		tile.SetSize(0)
		require.Equal(40, tile.DisplaySize, "displaySize did not match")
	})

	t.Run("positive", func(t *testing.T) {
		tile.SetSize(4)
		require.Equal(80, tile.DisplaySize, "displaySize did not match")
	})

	t.Run("negative", func(t *testing.T) {
		tile.SetSize(-4)
		require.Equal(40, tile.DisplaySize, "displaySize did not match")
	})
}
