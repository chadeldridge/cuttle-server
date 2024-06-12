package profiles

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTilesDefaultTile(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	require.False(tile.hideCmd, "hideCmd returned true")
	require.False(tile.hideExp, "hideExp returned true")
	require.Empty(tile.name, "name was not empty")
	require.Empty(tile.cmd, "cmd was not empty")
	require.Empty(tile.exp, "exp was not empty")
	require.Equal(40, tile.displaySize, "displaySize did not match")
}

func TestTilesNewTile(t *testing.T) {
	require := require.New(t)
	name := "Ping"
	cmd := "ping google.com"
	exp := "64 bytes from"

	t.Run("valid", func(t *testing.T) {
		tile := NewTile(name, cmd, exp)
		require.False(tile.hideCmd, "hideCmd returned true")
		require.False(tile.hideExp, "hideExp returned true")
		require.Equal(name, tile.name, "name did not match")
		require.Equal(cmd, tile.cmd, "cmd did not match")
		require.Equal(exp, tile.exp, "exp did not match")
		require.Equal(40, tile.displaySize, "displaySize did not match")
	})
	// TODO: Add test cases for invalid name, cmd, and exp after validation is added to NewTile.
}

func TestTilesHideCmd(t *testing.T) {
	require := require.New(t)
	tile := Tile{}

	t.Run("default", func(t *testing.T) {
		require.False(tile.HideCmd(), "Tile.HideCmd() returned true")
	})

	t.Run("true", func(t *testing.T) {
		tile.hideCmd = true
		require.True(tile.HideCmd(), "Tile.HideCmd() returned false")
	})

	t.Run("false", func(t *testing.T) {
		tile.hideCmd = false
		require.False(tile.HideCmd(), "Tile.HideCmd() returned true")
	})
}

func TestTilesHideExp(t *testing.T) {
	require := require.New(t)
	tile := Tile{}

	t.Run("default", func(t *testing.T) {
		require.False(tile.HideExp(), "Tile.HideExp() returned true")
	})

	t.Run("true", func(t *testing.T) {
		tile.hideCmd = true
		require.True(tile.HideExp(), "Tile.HideExp() returned false")
	})

	t.Run("false", func(t *testing.T) {
		tile.hideCmd = false
		require.False(tile.HideExp(), "Tile.HideExp() returned true")
	})
}

func TestTilesName(t *testing.T) {
	require := require.New(t)
	name := "Ping"
	tile := Tile{name: name}
	require.Equal(name, tile.Name(), "name did not match")
}

func TestTilesCmd(t *testing.T) {
	require := require.New(t)
	cmd := "ping google.com"
	tile := Tile{cmd: cmd}
	require.Equal(cmd, tile.Cmd(), "cmd did not match")
}

func TestTilesExp(t *testing.T) {
	require := require.New(t)
	exp := "64 bytes from"
	tile := Tile{exp: exp}
	require.Equal(exp, tile.Exp(), "exp did not match")
}

func TestTilesDisplaySize(t *testing.T) {
	require := require.New(t)
	tile := Tile{displaySize: 40}
	require.Equal(40, tile.DisplaySize(), "DisplaySize did not match")
}

func TestTilesSetHideCmd(t *testing.T) {
	require := require.New(t)
	tile := Tile{}

	t.Run("default", func(t *testing.T) {
		require.False(tile.hideCmd, "hideCmd returned true")
	})

	t.Run("true", func(t *testing.T) {
		tile.SetHideCmd(true)
		require.True(tile.hideCmd, "hideCmd returned false")
	})

	t.Run("false", func(t *testing.T) {
		tile.SetHideCmd(false)
		require.False(tile.hideCmd, "hideCmd returned true")
	})
}

func TestTilesSetHideExp(t *testing.T) {
	require := require.New(t)
	tile := Tile{}

	t.Run("default", func(t *testing.T) {
		require.False(tile.hideExp, "hideExp returned true")
	})

	t.Run("true", func(t *testing.T) {
		tile.SetHideExp(true)
		require.True(tile.hideExp, "hideExp returned false")
	})

	t.Run("false", func(t *testing.T) {
		tile.SetHideExp(false)
		require.False(tile.hideExp, "hideExp returned true")
	})
}

func TestTilesSetName(t *testing.T) {
	require := require.New(t)
	tile := NewTile("OldName", "", "")
	name := "NewName"

	t.Run("valid name", func(t *testing.T) {
		err := tile.SetName(name)
		require.NoError(err, "SetName() returned an error: %s", err)
		require.Equal(name, tile.name, "name did not match")
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
		require.Equal(name, tile.name, "name did not match")
	})
}

func TestTilesSetSize(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	t.Run("default", func(t *testing.T) {
		tile.SetSize(0)
		require.Equal(40, tile.displaySize, "displaySize did not match")
	})

	t.Run("positive", func(t *testing.T) {
		tile.SetSize(4)
		require.Equal(80, tile.displaySize, "displaySize did not match")
	})

	t.Run("negative", func(t *testing.T) {
		tile.SetSize(-4)
		require.Equal(40, tile.displaySize, "displaySize did not match")
	})
}

func TestTilesSetCmd(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	t.Run("valid", func(t *testing.T) {
		err := tile.SetCmd("ping google.com")
		require.NoError(err, "SetCmd() returned an error: %s", err)
	})

	// INCOMPLETE: Update after SetCmd checks for valid command.
	/*
		t.Run("invalid", func(t *testing.T) {
			err := tile.SetCmd("invalid command")
			require.Error(err, "SetCmd() did not return an error")
		})
	*/

	t.Run("empty", func(t *testing.T) {
		err := tile.SetCmd("")
		require.Error(err, "SetCmd() did not return an error")
	})
}

func TestTilesSetExp(t *testing.T) {
	require := require.New(t)
	tile := DefaultTile()

	t.Run("valid", func(t *testing.T) {
		err := tile.SetExp("64 bytes from")
		require.NoError(err, "SetExp() returned an error: %s", err)
	})

	/*
		t.Run("invalid", func(t *testing.T) {
			err := tile.SetExp("invalid expect string")
			require.Error(err, "SetExp() did not return an error")
		})
	*/

	t.Run("empty", func(t *testing.T) {
		err := tile.SetExp("")
		require.Error(err, "SetExp() did not return an error")
	})
}
