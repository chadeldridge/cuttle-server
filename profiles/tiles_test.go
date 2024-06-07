package profiles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultTile(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	require.False(tile.hideCmd, "Expected hideCmd to be false")
	require.False(tile.hideExp, "Expected hideExp to be false")
	require.Empty(tile.name, "Expected name to be empty")
	require.Empty(tile.cmd, "Expected cmd to be empty")
	require.Empty(tile.exp, "Expected exp to be empty")
	require.Equal(40, tile.displaySize, "Expected displaySize to be 40")
}

func TestNewTile(t *testing.T) {
	require := require.New(t)

	name := "Ping"
	cmd := "ping google.com"
	exp := "64 bytes from"

	tile := NewTile(name, cmd, exp)

	require.False(tile.hideCmd, "Expected hideCmd to be false")
	require.False(tile.hideExp, "Expected hideExp to be false")
	require.Equal(name, tile.name, "Expected name to be %s", name)
	require.Equal(cmd, tile.cmd, "Expected cmd to be %s", cmd)
	require.Equal(exp, tile.exp, "Expected exp to be %s", exp)
	require.Equal(40, tile.displaySize, "Expected displaySize to be 40")
}

func TestTileHideCmd(t *testing.T) {
	require := require.New(t)

	tile := DefaultTile()
	require.False(tile.HideCmd(), "Expected HideCmd to return false")

	tile.SetHideCmd(true)
	require.True(tile.HideCmd(), "Expected HideCmd to return true")
}

func TestTileHideExp(t *testing.T) {
	require := require.New(t)

	tile := DefaultTile()
	require.False(tile.HideExp(), "Expected HideExp to return false")

	tile.SetHideExp(true)
	require.True(tile.HideExp(), "Expected HideExp to return true")
}

func TestTileName(t *testing.T) {
	require := require.New(t)
	name := "Ping"

	tile := NewTile(name, "", "")
	require.Equal(name, tile.Name(), "Expected Name to return %s", name)
}

func TestTileCmd(t *testing.T) {
	require := require.New(t)
	cmd := "ping google.com"

	tile := NewTile("", cmd, "")
	require.Equal(cmd, tile.Cmd(), "Expected Cmd to return %s", cmd)
}

func TestTileExp(t *testing.T) {
	require := require.New(t)
	exp := "64 bytes from"

	tile := NewTile("", "", exp)
	require.Equal(exp, tile.Exp(), "Expected Exp to return %s", exp)
}

func TestTileDisplaySize(t *testing.T) {
	require := require.New(t)

	tile := DefaultTile()
	require.Equal(40, tile.DisplaySize(), "Expected DisplaySize to return 40")

	tile.SetSize(3)
	require.Equal(60, tile.DisplaySize(), "Expected DisplaySize to return 60")
}

func TestTileSetName(t *testing.T) {
	require := require.New(t)
	tile := NewTile("OldName", "", "")

	// Test case 1: Valid name
	err := tile.SetName("NewName")
	require.NoError(err, "Expected SetName to return nil")
	require.Equal("NewName", tile.name, "Expected Exp to return NewName")

	// Test case 2: Invalid name
	// INCOMPLETE: Change to require.Error after SetName checks for valid name.
	// err = tile.SetName("invalid name")
	// require.Error(err, "Expected SetName to return an error")
	// require.Equal("NewName", tile.name, "Expected Exp to return NewName")

	// Test case 3: Empty name
	err = tile.SetName("")
	require.Error(err, "Expected SetName to return an error")
	require.Equal("NewName", tile.name, "Expected Exp to return NewName")
}

func TestTileSetSize(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	tile.SetSize(0)
	require.Equal(40, tile.displaySize, "Expected displaySize to be 40")

	tile.SetSize(4)
	require.Equal(80, tile.displaySize, "Expected displaySize to be 80")
}

func TestTileSetCmd(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	// Test case 1: Valid command
	err := tile.SetCmd("ping google.com")
	require.NoError(err, "Expected SetCmd to return nil")

	// Test case 2: Invalid command
	err = tile.SetCmd("invalid command")
	// INCOMPLETE: Update after SetCmd checks for valid command.
	// require.Error(err, "Expected SetCmd to return an error")
	require.NoError(err, "Expected SetCmd to return nil")

	// Test case 3: Empty command
	err = tile.SetCmd("")
	require.Error(err, "Expected SetCmd to return an error")
}

func TestTileSetExp(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	err := tile.SetExp("invalid expect string")
	// INCOMPLETE: Update after SetCmd checks for valid command.
	// require.Error(err, "Expected SetExp to return an error")
	require.NoError(err, "Expected SetExp to return nil")

	err = tile.SetExp("64 bytes from")
	require.NoError(err, "Expected SetExp to return nil")

	err = tile.SetExp("")
	require.Error(err, "Expected SetExp to return an error")
}
