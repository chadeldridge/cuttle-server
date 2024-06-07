package profiles

import (
	"testing"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/stretchr/testify/require"
)

func TestNewProfile(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	// Test case 1: Valid profile name and groups
	group1 := NewGroup("Group1", testServers...)
	group2 := NewGroup("Group2", testServers[:1]...)
	profile := NewProfile("TestProfile", group1, group2)

	require.Equal("TestProfile", profile.Name, "Expected profile name to be TestProfile")
	require.Len(profile.Groups, 2, "Expected 2 groups in profile")
	require.Contains(profile.Groups, "Group1", "Expected Group1 in profile")
	require.Contains(profile.Groups, "Group2", "Expected Group2 in profile")
}

func TestSetName(t *testing.T) {
	require := require.New(t)
	profile := Profile{Name: "OldName"}

	// Test case 1: Valid name
	err := profile.SetName("NewName")
	require.NoError(err)
	require.Equal("NewName", profile.Name)

	// Test case 2: Empty name
	err = profile.SetName("")
	require.Error(err)
	require.Equal("NewName", profile.Name)
}

func testNewTile(t *testing.T, name string) Tile {
	require := require.New(t)
	tile := NewTile(name, "echo Hello", "Hello")

	require.Equal(name, tile.Name(), "Expected tile name to be %s", name)
	return tile
}

func TestAddTiles(t *testing.T) {
	require := require.New(t)
	tile1 := testNewTile(t, "Tile1")
	tile2 := testNewTile(t, "Tile2")
	profile := Profile{Tiles: make(map[string]Tile)}

	// Test case 1: Add new tiles
	err := profile.AddTiles(tile1, tile2)
	require.NoError(err)
	require.Len(profile.Tiles, 2)
	require.Contains(profile.Tiles, "Tile1")
	require.Contains(profile.Tiles, "Tile2")

	// Test case 2: Add existing tile
	tile3 := testNewTile(t, "Tile1")
	err = profile.AddTiles(tile3)
	require.Error(err)
	require.Len(profile.Tiles, 2)

	// Test case 2: Add no tiles
	err = profile.AddTiles()
	require.Error(err)
	require.Len(profile.Tiles, 2)
}

func TestAddGroups(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group1 := NewGroup("Group1", testServers...)
	group2 := NewGroup("Group2", testServers[:1]...)
	profile := NewProfile("TestProfile", group1)

	// Test case 1: Add new groups
	err := profile.AddGroups(group2)
	require.NoError(err)
	require.Len(profile.Groups, 2)
	require.Contains(profile.Groups, "Group1")
	require.Contains(profile.Groups, "Group2")

	// Test case 2: Add existing group
	group3 := NewGroup("Group1", testServers[:1]...)
	err = profile.AddGroups(group3)
	require.Error(err)
	require.Len(profile.Groups, 2)

	// Test case 2: Add no groups
	err = profile.AddGroups()
	require.Error(err)
	require.Len(profile.Groups, 2)
}

func TestGetTile(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	tile := NewTile("Tile1", "echo Hello", "Hello")

	// Test case 1: Valid tile and group
	group1 := NewGroup("Group1", testServers...)
	profile := NewProfile("TestProfile", group1)
	err := profile.AddTiles(tile)
	require.NoError(err, "Error adding tile to profile")

	got, err := profile.GetTile("Tile1")
	require.NoError(err)
	require.Equal(tile, got)

	// Test case 2: Non-existing tile
	got, err = profile.GetTile("InvalidTile")
	require.Error(err)
	require.Equal(Tile{}, got)

	// Test case 3: Empty tile name
	got, err = profile.GetTile("")
	require.Error(err)
	require.Equal(Tile{}, got)
}

func TestGetGroup(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group1 := NewGroup("Group1", testServers...)
	profile := NewProfile("TestProfile", group1)

	// Test case 1: Valid group
	got, err := profile.GetGroup("Group1")
	require.NoError(err)
	require.Equal("Group1", got.Name)

	// Test case 2: Non-existing group
	got, err = profile.GetGroup("InvaliGroup")
	require.Error(err)
	require.Equal(Group{}, got)

	// Test case 3: Empty group name
	got, err = profile.GetGroup("")
	require.Error(err)
	require.Equal(Group{}, got)
}

func TestExecute(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)

	tile1 := NewTile("Tile1", "echo Hello", "Hello")
	tile2 := NewTile("Tile2", "this is not a command", "Hello")
	group1 := NewGroup("Group1", testServers...)
	profile := NewProfile("TestProfile", group1)

	err := profile.AddTiles(tile1)
	require.NoError(err, "Error adding tile to profile")

	err = profile.AddTiles(tile2)
	require.NoError(err, "Error adding tile to profile")

	// Test case 1: Valid tile and group
	err = profile.Execute("Tile1", "Group1")
	require.NoError(err)
	require.NotContains(results.String(), "failed")
	results.Reset()
	logs.Reset()

	// Test case 2: Invalid tile
	err = profile.Execute("InvalidTile", "Group1")
	require.Error(err)

	// Test case 3: Invalid group
	err = profile.Execute("Tile1", "InvalidGroup")
	require.Error(err)

	// Test case 4: Faile running command
	err = profile.Execute("Tile2", "Group1")
	require.Error(err)
	require.NotContains(results.String(), "ok")
	results.Reset()
	logs.Reset()

	// Test case 5: Connection error
	connections.Pool.CloseAll()
	initGroupTest(t, true)
	group2 := NewGroup("Group2", testServers...)
	err = profile.AddGroups(group2)
	require.NoError(err)

	err = profile.Execute("Tile1", "Group2")
	require.Error(err)
}
