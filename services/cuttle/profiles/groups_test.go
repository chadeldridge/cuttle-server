package profiles

import (
	"bytes"
	"testing"

	"github.com/chadeldridge/cuttle-server/services/cuttle/connections"
	"github.com/stretchr/testify/require"
)

const (
	name1 = "Internal Servers"
	name2 = "Web Servers"
)

var (
	testServers []connections.Server
	results     bytes.Buffer
	logs        bytes.Buffer

	testServerNames = []string{
		"host1",
		"host2",
		"host3",
	}
)

func initGroupTest(t *testing.T, errConn bool) {
	testServers = []connections.Server{}
	for _, h := range testServerNames {
		testServers = append(testServers, createNewServer(t, h, errConn))
	}
}

func createNewServer(t *testing.T, host string, errConn bool) connections.Server {
	s, err := connections.NewServer(host, 0, &results, &logs)
	if err != nil {
		t.Fatal("", err)
	}

	conn, err := connections.NewMockConnector("my connector", "test")
	if err != nil {
		t.Fatalf("failed to create mock connector: %s", err)
	}
	conn.ErrOnConnectionOpen(errConn)

	err = s.SetConnector(&conn)
	if err != nil {
		t.Fatalf("failed to add connector to server: %s", err)
	}

	return s
}

func testNewGroup(t *testing.T, servers ...connections.Server) {
	require := require.New(t)
	group := NewGroup(name1, servers...)
	require.Equal(name1, group.Name, "group name does not match expected name")
	require.Len(group.Servers, len(servers), "missing servers in group")

	for i := range group.Count() {
		require.Equal(servers[i].Name, group.Servers[i].Name, "server name did not match")
	}
}

func TestGroupsNewGroup(t *testing.T) {
	initGroupTest(t, false)

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
	initGroupTest(t, false)
	require := require.New(t)

	t.Run("single server", func(t *testing.T) {
		group := Group{Name: name1, Servers: testServers[:1]}
		require.Equal(1, group.Count(), "missing server in group")
	})

	t.Run("all servers", func(t *testing.T) {
		group := Group{Name: name1, Servers: testServers}
		require.Equal(3, group.Count(), "missing servers in group")
	})

	t.Run("no servers", func(t *testing.T) {
		group := Group{Name: name1}
		require.Equal(0, group.Count(), "servers found in group")
	})
}

func TestGroupsSetName(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group := Group{Name: name1}

	t.Run("name2", func(t *testing.T) {
		err := group.SetName(name2)
		require.NoError(err, "Group.SetName returned an error: %s", err)
		require.Equal(name2, group.Name, "group name does not match")
	})

	t.Run("empty name", func(t *testing.T) {
		err := group.SetName("")
		require.Error(err, "Group.SetName did not return an error")
		require.Equal(name2, group.Name, "group name does nto match")
	})
}

func TestGroupsAddServers(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)

	t.Run("add one", func(t *testing.T) {
		group := Group{Name: name1}
		group.AddServers(testServers[0])
		require.Len(group.Servers, 1, "missing servers in group")
	})

	t.Run("add multiple", func(t *testing.T) {
		group := Group{Name: name1}
		group.AddServers(testServers[:2]...)
		require.Len(group.Servers, 2, "missing servers in group")
	})

	t.Run("add duplicate", func(t *testing.T) {
		group := Group{Name: name1, Servers: testServers}
		group.AddServers(testServers[0])
		require.Len(group.Servers, 3, "missing servers in group")
	})

	t.Run("add none", func(t *testing.T) {
		group := Group{Name: name1, Servers: testServers}
		group.AddServers()
		require.Len(group.Servers, 3, "missing servers in group")
	})
}

func TestGroupsReset(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group := Group{Name: name1, Servers: testServers}

	for i := 0; i < 2; i++ {
		_, err := group.Next()
		require.NoError(err, "Group.Next() returned an error: %s", err)
	}

	group.Reset()
	s, err := group.Next()
	require.NoError(err, "Group.Next() returned an error: %s", err)
	require.Equal(testServers[0].Name, s.Name, "server name did not match")
}

func TestGroupsNext(t *testing.T) {
	initGroupTest(t, false)
	require := require.New(t)
	group := Group{Name: name1, Servers: testServers}

	for i := 0; i <= group.Count(); i++ {
		s, err := group.Next()
		if i == group.Count() {
			require.Equal(ErrEndOfList, err, "did not get 'end of list' error")
			return
		}

		require.NoError(err, "Group.Next() returned an error: %s", err)
		require.Equal(testServers[i].Name, s.Name, "server name did not match")
	}
}

func TestGroupsUniq(t *testing.T) {
	require := require.New(t)
	initGroupTest(t, false)
	group := Group{Name: name1, Servers: testServers}

	group.Servers = append(group.Servers, testServers[0])
	require.Len(group.Servers, 4, "added server missing")

	group.uniq()
	require.Len(group.Servers, 3, "duplicate servers not removed")

	for i := range group.Count() {
		require.Equal(testServers[i].Name, group.Servers[i].Name, "server name did not match")
	}
}
