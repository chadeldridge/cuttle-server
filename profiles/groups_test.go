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
	servers []connections.Server
	results *bytes.Buffer
	logs    *bytes.Buffer

	serverInputs = []string{
		"test.home",
		"internal1",
		"internal2",
	}
)

func initGroupTest(t *testing.T) {
	for _, h := range serverInputs {
		s, err := connections.NewServer(h, 0, results, logs)
		if err != nil {
			t.Fatalf("profiles.TestGroupsNewGroup: server1 creation failed: %s", err)
		}
		servers = append(servers, s)
	}
}

func TestGroupsNewGroup(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, servers...)
	require.Equal(t, name1, got.Name, "profiles.TestGroupsNewGroup: group name does nto match")
	require.Greater(t, got.ServerCount(), 0, "profiles.TestGroupsNewGroup: missing servers in group")

	for i := range got.Count() {
		require.Equal(
			t,
			servers[i].Name(),
			got.Servers[i].Name(),
			"profiles.TestGroupsNewGroup: server name did not match",
		)
	}
}

func TestGroupsSetName(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, servers...)
	require.Equal(t, name1, got.Name, "profiles.TestGroupsSetName: group name does nto match")
	require.Greater(t, got.ServerCount(), 0, "profiles.TestGroupsSetName: missing servers in group")

	got.SetName(name2)
	require.Equal(t, name2, got.Name, "profiles.TestGroupsSetName: group name does nto match")
}

func TestGroupsAddServers(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, servers[0], servers[1])
	require.Equal(t, 2, got.ServerCount(), "profiles.TestGroupsAddServers: missing servers in group")

	got.AddServers(servers[2])
	require.Equal(t, 3, got.ServerCount(), "profiles.TestGroupsAddServers: missing servers in group")
	for i := range got.Count() {
		require.Equal(
			t,
			servers[i].Name(),
			got.Servers[i].Name(),
			"profiles.TestGroupsAddServers: server name did not match",
		)
	}
}

func TestGroupsNext(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, servers...)
	require.Greater(t, got.ServerCount(), 0, "profiles.TestGroupsNext: missing servers in group")

	for i := 0; i <= got.Count(); i++ {
		s, err := got.Next()
		if i == got.Count() {
			require.Equal(t, ErrEndOfList, err.Error(), "profiles.TestGroupsNext: did not get 'end of list' error")
			return
		}

		if err != nil {
			t.Fatalf("profiles.TestGroupsNext: error getting next server: %s", err)
		}
		require.Equal(t, servers[i].Name(), s.Name(), "profiles.TestGroupsNext: server name did not match")
	}
}

func TestGroupsReset(t *testing.T) {
	initGroupTest(t)
	got := NewGroup(name1, servers...)
	require.Greater(t, got.ServerCount(), 0, "profiles.TestGroupsReset: missing servers in group")

	for i := 0; i < 2; i++ {
		s, err := got.Next()
		if err != nil {
			t.Fatalf("profiles.TestGroupsReset: error getting next server: %s", err)
		}
		require.Equal(t, servers[i].Name(), s.Name(), "profiles.TestGroupsReset: server name did not match")
	}

	got.Reset()
	s, err := got.Next()
	if err != nil {
		t.Fatalf("profiles.TestGroupsReset: error getting next server: %s", err)
	}
	require.Equal(t, servers[0].Name(), s.Name(), "profiles.TestGroupsReset: server name did not match")
}
