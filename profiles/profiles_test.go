package profiles

import (
	"testing"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/stretchr/testify/require"
)

func testNewTile(name string) Tile { return Tile{Name: name, Cmd: "echo Hello", Exp: "Hello"} }

func TestProfilesNewProfile(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group1 := Group{Name: "Group1", Servers: testServers}
	group2 := Group{Name: "Group2", Servers: testServers[:1]}

	t.Run("single group", func(t *testing.T) {
		profile, err := NewProfile("TestProfile", group1)
		require.NoError(err, "NewProfile() returned an error: %s", err)
		require.Equal("TestProfile", profile.Name, "Expected profile name to be TestProfile")
		require.Len(profile.Groups, 1, "missing groups in profile")
		require.Contains(profile.Groups, "Group1", "Expected Group1 in profile")
	})

	t.Run("multiple groups", func(t *testing.T) {
		profile, err := NewProfile("TestProfile", group1, group2)
		require.NoError(err, "NewProfile() returned an error: %s", err)
		require.Len(profile.Groups, 2, "Expected 2 groups in profile")
		require.Contains(profile.Groups, "Group1", "Expected Group1 in profile")
		require.Contains(profile.Groups, "Group2", "Expected Group2 in profile")
	})

	t.Run("no groups", func(t *testing.T) {
		profile, err := NewProfile("TestProfile")
		require.NoError(err, "NewProfile() returned an error: %s", err)
		require.Empty(profile.Groups, "Groups was not empty")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := NewProfile("")
		require.Error(err, "NewProfile() did not return an error")
	})

	// INCOMPLETE: Add check after name validation is implemented.
	/*
		t.Run("invalid name", func(t *testing.T) {
			_, err := NewProfile("invalid name")
			require.Error(err, "NewProfile() did not return an error")
		})
	*/
}

func TestProfilesSetName(t *testing.T) {
	require := require.New(t)
	profile := Profile{Name: "Profile1"}
	name := "Profile2"

	t.Run("valid", func(t *testing.T) {
		err := profile.SetName(name)
		require.NoError(err, "SetName() returned an error: %s", err)
		require.Equal(name, profile.Name)
	})

	// INCOMPLETE: Add check after name validation is implemented.
	/*
		t.Run("invalid", func(t *testing.T) {
			err := profile.SetName("some invalid name")
			require.Error(err, "SetName() did not return an error")
			require.Equal(name, profile.Name)
		})
	*/

	t.Run("empty", func(t *testing.T) {
		err := profile.SetName("")
		require.Error(err, "SetName() did not return an error")
		require.Equal(name, profile.Name)
	})
}

func TestProfilesAddTiles(t *testing.T) {
	require := require.New(t)
	tile1 := testNewTile("Tile1")
	tile2 := testNewTile("Tile2")
	tile3 := testNewTile("Tile3")
	dupeTile := testNewTile("Tile1")
	profile := Profile{Tiles: make(map[string]Tile)}

	t.Run("add one", func(t *testing.T) {
		err := profile.AddTiles(tile1)
		require.NoError(err, "AddTiles() returned an error: %s", err)
		require.Len(profile.Tiles, 1)
		require.Contains(profile.Tiles, "Tile1")
	})

	t.Run("add multiple", func(t *testing.T) {
		err := profile.AddTiles(tile2, tile3)
		require.NoError(err, "AddTiles() returned an error: %s", err)
		require.Len(profile.Tiles, 3)
		require.Contains(profile.Tiles, "Tile1")
		require.Contains(profile.Tiles, "Tile2")
		require.Contains(profile.Tiles, "Tile3")
	})

	t.Run("add existing", func(t *testing.T) {
		err := profile.AddTiles(dupeTile)
		require.Error(err, "AddTiles() did not return an error")
		require.Len(profile.Tiles, 3)
	})

	t.Run("add none", func(t *testing.T) {
		err := profile.AddTiles()
		require.Error(err, "AddTiles() did not return an error")
		require.Len(profile.Tiles, 3)
	})
}

func TestProfilesAddGroups(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group1 := Group{Name: "Group1"}
	group2 := Group{Name: "Group2"}
	group3 := Group{Name: "Group3"}
	dupeGroup := Group{Name: "Group1"}
	profile := Profile{Name: "TestProfile", Groups: make(map[string]Group)}

	t.Run("add one", func(t *testing.T) {
		err := profile.AddGroups(group1)
		require.NoError(err, "AddGroups() returned an error: %s", err)
		require.Len(profile.Groups, 1)
		require.Contains(profile.Groups, "Group1")
	})

	t.Run("add multiple", func(t *testing.T) {
		err := profile.AddGroups(group2, group3)
		require.NoError(err, "AddGroups() returned an error: %s", err)
		require.Len(profile.Groups, 3)
		require.Contains(profile.Groups, "Group1")
		require.Contains(profile.Groups, "Group2")
		require.Contains(profile.Groups, "Group3")
	})

	t.Run("add existing", func(t *testing.T) {
		err := profile.AddGroups(dupeGroup)
		require.Error(err, "AddGroups() did not return an error")
		require.Len(profile.Groups, 3)
	})

	t.Run("add none", func(t *testing.T) {
		err := profile.AddGroups()
		require.Error(err, "AddGroups() did not return an error")
		require.Len(profile.Groups, 3)
	})
}

func TestProfilesGetTile(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	tile := testNewTile("Tile1")
	profile := Profile{Name: "TestProfile", Tiles: map[string]Tile{"Tile1": tile}}

	t.Run("exists", func(t *testing.T) {
		got, err := profile.GetTile("Tile1")
		require.NoError(err, "GetTile() returned an error: %s", err)
		require.Equal(tile, got)
	})

	t.Run("does not exist", func(t *testing.T) {
		got, err := profile.GetTile("InvalidTile")
		require.Error(err, "GetTile() did not return an error")
		require.Equal(Tile{}, got)
	})

	t.Run("empty name", func(t *testing.T) {
		got, err := profile.GetTile("")
		require.Error(err, "GetTile() did not return an error")
		require.Equal(Tile{}, got)
	})
}

func TestProfilesGetGroup(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group1 := Group{Name: "Group1"}
	profile := Profile{Name: "TestProfile", Groups: map[string]Group{"Group1": group1}}

	t.Run("exists", func(t *testing.T) {
		got, err := profile.GetGroup("Group1")
		require.NoError(err, "GetGroup() returned an error: %s", err)
		require.Equal("Group1", got.Name)
	})

	t.Run("does not exist", func(t *testing.T) {
		got, err := profile.GetGroup("missing group")
		require.Error(err, "GetGroup() did not return an error")
		require.Equal(Group{}, got)
	})

	t.Run("empty name", func(t *testing.T) {
		got, err := profile.GetGroup("")
		require.Error(err, "GetGroup() did not return an error")
		require.Equal(Group{}, got)
	})
}

func TestProfilesExecute(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	tile1 := Tile{Name: "Tile1", Cmd: "echo Hello", Exp: "Hello"}
	tile2 := Tile{Name: "Tile2", Cmd: "this is not a command", Exp: "Hello"}
	group1 := Group{Name: "Group1", Servers: testServers}
	profile := Profile{
		Name:   "TestProfile",
		Tiles:  map[string]Tile{"Tile1": tile1, "Tile2": tile2},
		Groups: map[string]Group{"Group1": group1},
	}

	t.Run("valid", func(t *testing.T) {
		err := profile.Execute("Tile1", "Group1")
		require.NoError(err, "Execute() returned an error: %s", err)
		require.NotContains(results.String(), "failed")
		results.Reset()
		logs.Reset()
	})

	t.Run("invalid tile", func(t *testing.T) {
		err := profile.Execute("InvalidTile", "Group1")
		require.Error(err, "Execute() did not return an error")
	})

	t.Run("invalid group", func(t *testing.T) {
		err := profile.Execute("Tile1", "InvalidGroup")
		require.Error(err, "Execute() did not return an error")
	})

	t.Run("invalid command", func(t *testing.T) {
		err := profile.Execute("Tile2", "Group1")
		require.Error(err, "Execute() did not return an error")
		require.NotContains(results.String(), "ok")
		results.Reset()
		logs.Reset()
	})

	t.Run("connection error", func(t *testing.T) {
		connections.Pool.CloseAll()
		initGroupTest(t, true)
		group2 := NewGroup("Group2", testServers...)
		profile.Groups["Group2"] = group2

		err := profile.Execute("Tile1", "Group2")
		require.Error(err, "Execute() did not return an error")
	})
}
