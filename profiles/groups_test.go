package profiles

import (
	"bytes"
	"testing"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/stretchr/testify/require"
)

const (
	name1 = "Internal Servers"
	name2 = "Web Servers"
)

var (
	testServers []connections.Server
	results     *bytes.Buffer
	logs        *bytes.Buffer

	serverInputs = []string{
		"test.home",
		"internal1",
		"internal2",
	}
)

func initGroupTest(t *testing.T) {
	testServers = []connections.Server{}
	for _, h := range serverInputs {
		testServers = append(testServers, createNewServer(t, h))
	}
}

func createNewServer(t *testing.T, host string) connections.Server {
	s, err := connections.NewServer(host, 0, results, logs)
	if err != nil {
		t.Fatal("profiles.TestGroupsNewServer: server creation failed: ", err)
	}
	return s
}

func testNewGroup(t *testing.T, servers ...connections.Server) {
	got := NewGroup(name1, servers...)
	require.Equal(t, name1, got.Name, "profiles.TestGroupsNewGroup: group name does not match")
	require.Equal(t, len(servers), len(got.Servers), "profiles.TestGroupsNewGroup: missing servers in group")

	for i := range got.Count() {
		require.Equal(
			t,
			servers[i].Name(),
			got.Servers[i].Name(),
			"profiles.TestGroupsNewGroup: server name did not match",
		)
	}
}

func TestGroupsNewGroup(t *testing.T) {
	initGroupTest(t)
	t.Run("single server", func(t *testing.T) {
		testNewGroup(t, testServers[0])
	})

	t.Run("all servers", func(t *testing.T) {
		testNewGroup(t, testServers...)
	})

	t.Run("no servers", func(t *testing.T) {
		testNewGroup(t)
	})
}

func TestGroupsCount(t *testing.T) {
	initGroupTest(t)
	t.Run("single server", func(t *testing.T) {
		got := NewGroup(name1, testServers[0])
		require.Equal(t, 1, got.Count(),
			"profiles.TestGroupsNewGroup: missing servers in group",
		)
	})

	t.Run("all servers", func(t *testing.T) {
		got := NewGroup(name1, testServers...)
		require.Equal(t, len(testServers), got.Count(),
			"profiles.TestGroupsNewGroup: missing servers in group",
		)
	})

	t.Run("no servers", func(t *testing.T) {
		got := NewGroup(name1)
		require.Equal(t, 0, got.Count(),
			"profiles.TestGroupsNewGroup: missing servers in group",
		)
	})
}

func TestGroupsSetName(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, testServers...)

	t.Run("name2", func(t *testing.T) {
		err := got.SetName(name2)
		require.Nil(t, err, "profiles.TestGroupsSetName: Group.SetName() returned an error: ", err)
		require.Equal(t, name2, got.Name, "profiles.TestGroupsSetName: group name does nto match")
	})

	t.Run("empty name", func(t *testing.T) {
		err := got.SetName("")
		require.NotNil(t, err, "profiles.TestGroupsSetName: Group.SetName() did not return an error")
		require.Equal(t, name2, got.Name, "profiles.TestGroupsSetName: group name does nto match")
	})
}

func TestGroupsAddServers(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, testServers[0], testServers[1])
	require.Equal(t, 2, len(got.Servers), "profiles.TestGroupsAddServers: missing servers in group")

	t.Run("add server", func(t *testing.T) {
		got.AddServers(testServers[2])
		require.Equal(t, 3, len(got.Servers), "profiles.TestGroupsAddServers: missing servers in group")
		for i := range got.Count() {
			require.Equal(
				t,
				testServers[i].Name(),
				got.Servers[i].Name(),
				"profiles.TestGroupsAddServers: server name did not match",
			)
		}
	})

	got = NewGroup(name1, testServers[0], testServers[1])
	require.Equal(t, 2, len(got.Servers), "profiles.TestGroupsAddServers: missing servers in group")
	t.Run("empty list", func(t *testing.T) {
		got.AddServers()
		require.Equal(t, 2, len(got.Servers), "profiles.TestGroupsAddServers: missing servers in group")
	})
}

func TestGroupsReset(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, testServers...)
	require.Equal(t, len(testServers), len(got.Servers), "profiles.TestGroupsReset: missing servers in group")

	for i := 0; i < 2; i++ {
		_, err := got.Next()
		if err != nil {
			t.Fatalf("profiles.TestGroupsReset: error getting next server: %s", err)
		}
	}

	got.Reset()
	s, err := got.Next()
	require.Nil(t, err, "Group.Next() returned an error: ", err)
	require.Equal(t, testServers[0].Name(), s.Name(), "profiles.TestGroupsReset: server name did not match")
}

func TestGroupsNext(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, testServers...)
	require.Equal(t, len(testServers), len(got.Servers), "profiles.TestGroupsNext: missing servers in group")

	for i := 0; i <= got.Count(); i++ {
		s, err := got.Next()
		if i == got.Count() {
			require.Equal(t, ErrEndOfList, err, "profiles.TestGroupsNext: did not get 'end of list' error")
			return
		}

		require.Nil(t, err, "Group.Next() returned an error: ", err)
		require.Equal(t, testServers[i].Name(), s.Name(), "profiles.TestGroupsNext: server name did not match")
	}
}

func TestGroupsUniq(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, testServers...)
	require.Equal(t, len(testServers), len(got.Servers), "profiles.TestGroupsUniq: missing servers in group")

	got.Servers = append(got.Servers, testServers[0])
	require.Equal(t, len(testServers)+1, len(got.Servers), "profiles.TestGroupsUniq: added server missing")

	got.uniq()
	require.Equal(t, len(testServers), len(got.Servers), "profiles.TestGroupsUniq: duplicate servers not removed")

	for i := range got.Count() {
		require.Equal(
			t,
			testServers[i].Name(),
			got.Servers[i].Name(),
			"profiles.TestGroupsUniq: server name did not match",
		)
	}
}
