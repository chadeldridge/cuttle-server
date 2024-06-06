package connections

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPoolsOpen(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	t.Run("good connector", func(t *testing.T) {
		pConn, err := Pool.Open(&server)
		require.Nil(t, err, "Pool.Open() returned an error: ", err)
		require.True(t, time.Now().Before(pConn.killAt), "killAt was before time.Now()")
		pConn.Close(false)
	})

	t.Run("already open", func(t *testing.T) {
		pConn, err := Pool.Open(&server)
		require.Nil(t, err, "Pool.Open() returned an error: ", err)

		pConn2, err := Pool.Open(&server)
		require.Nil(t, err, "Pool.Open() returned an error: ", err)
		require.Equal(t, pConn, pConn2, "Connection refs were not the same")
		pConn.Close(false)
	})

	t.Run("empty connector", func(t *testing.T) {
		server.SetConnector(&MockConnector{})
		_, err := Pool.Open(&server)
		require.NotNil(t, err, "Pool.Open() did not return an error")
	})

	t.Run("empty hostname", func(t *testing.T) {
		server = Server{hostname: ""}
		require.Nil(t, err, "NewServer() returned an error: ", err)

		pConn, err := Pool.Open(&server)
		require.NotNil(t, err, "Pool.Open() returned an error: ", err)
		pConn.Close(false)
	})
}

func TestPoolsGetConnection(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)
	require.True(t, time.Now().Before(pConn.killAt), "killAt was before time.Now()")

	t.Run("existing connection", func(t *testing.T) {
		got := Pool.GetConnection(server)
		require.NotNil(t, got, "Pool.GetConnection() returned nil *Connection")
	})

	pConn.Close(false)

	t.Run("no connection", func(t *testing.T) {
		got := Pool.GetConnection(server)
		require.Nil(t, got, "Pool.GetConnection() returned a *Connection")
	})
}

func TestPoolsExpires(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)
	require.Greater(t, pConn.Expires(), time.Now(), "Connection.killAt was not after time.Now()")
	pConn.Close(false)
}

func TestPoolsExpired(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)
	require.False(t, pConn.Expired(), "Connection was already expired")
	defer pConn.Close(false)

	pConn.killAt = time.Now().Add(time.Minute * time.Duration(-5))
	require.True(t, pConn.Expired(), "Connection was not expired")
}

func TestPoolsExtend(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)
	require.True(t, time.Now().Before(pConn.killAt), "killAt was before time.Now()")
	defer pConn.Close(false)

	old := time.Now().Add(time.Minute * time.Duration(TTL))
	pConn.Extend(5)
	require.Greater(t, pConn.killAt, old, "Connection.killAt was not after old killAt time")
}

func TestPoolsClose(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)

	t.Run("session active", func(t *testing.T) {
		conn.OpenSession(server)
		err = pConn.Close(false)
		require.NotNil(t, err, "Connection.Close() did not return an error")
		_, ok := Pool[server.hostname]
		require.True(t, ok, "Connection not found after failed Pool.Close()")
		conn.CloseSession()
	})

	t.Run("connection close error", func(t *testing.T) {
		conn.ErrOnConnectionClose(true)
		err = pConn.Close(false)
		require.NotNil(t, err, "Connection.Close() did not return an error")
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection found after Pool.Close()")
		conn.ErrOnConnectionClose(false)
	})

	pConn, err = pConn.Open(Pool)
	require.Nil(t, err, "Connection.Open() returned an error: ", err)

	t.Run("close connection", func(t *testing.T) {
		err = pConn.Close(false)
		require.Nil(t, err, "Connection.Close() returned an error: ", err)
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection found after close")
	})

	t.Run("not in pool", func(t *testing.T) {
		err = pConn.Close(false)
		require.NotNil(t, err, "Pool.Close() did not return an error")
		require.Equal(t, "connections.Connection.Close: Connection not found in Pool",
			err.Error(), "Did not receive 'not found' error from Connection.Close()")
	})

	server.SetConnector(&conn)
	pConn, err = pConn.Open(Pool)
	require.Nil(t, err, "Connection.Open() returned an error: ", err)
	t.Run("pool close connection", func(t *testing.T) {
		err = Pool.Close(server.hostname, false)
		require.Nil(t, err, "Pool.Close() returned an error: ", err)
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection found after close")
	})

	t.Run("pool close no connection", func(t *testing.T) {
		err = Pool.Close(server.hostname, false)
		require.Nil(t, err, "Pool.Close() returned an error: ", err)
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection found after close")
	})
}

func createPool(t *testing.T, count int) {
	for i := range count {
		var res bytes.Buffer
		var logs bytes.Buffer

		server, err := NewServer(fmt.Sprintf("host%d", i), 0, &res, &logs)
		require.Nil(t, err, "NewServer() returned an error: ", err)

		conn, err := NewMockConnector(testUser)
		require.Nil(t, err, "NewMockConnector() returned an error: ", err)
		server.SetConnector(&conn)

		_, err = Pool.Open(&server)
		require.Nil(t, err, "Pool.Open() returned an error: ", err)
	}
}

func TestPoolsCloseAll(t *testing.T) {
	hostCount := 5
	createPool(t, hostCount)

	t.Run("connection error", func(t *testing.T) {
		conn4 := Pool["host4"].Server.Connector.(*MockConnector)
		conn4.ErrOnConnectionClose(true)

		require.Equal(t, hostCount, Pool.Count(), "Pool connection count did not matched expected amount")
		err := Pool.CloseAll()
		require.NotNil(t, err, "Pool.Close() did not return an error")
		require.Equal(t, 0, Pool.Count(), "Pool.Count() did not match expected count")

		conn4.ErrOnConnectionClose(false)
	})

	createPool(t, hostCount)
	t.Run("no errors", func(t *testing.T) {
		err := Pool.CloseAll()
		require.Nil(t, err, "Pool.Close() returned an error: ", err)
		require.Equal(t, 0, Pool.Count(), "Pool.Count() did not match expected count")
	})
}

func TestPoolsTimeOut(t *testing.T) {
	var res bytes.Buffer
	var logs bytes.Buffer

	server, err := NewServer(testHost, 0, &res, &logs)
	require.Nil(t, err, "NewServer() returned an error: ", err)

	conn, err := NewMockConnector(testUser)
	require.Nil(t, err, "NewMockConnector() returned an error: ", err)
	server.SetConnector(&conn)

	pConn, err := Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)

	t.Run("not expired", func(t *testing.T) {
		err := pConn.TimeOut()
		require.Nil(t, err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.hostname]
		require.True(t, ok, "Connection not found in Pool")
	})

	pConn.Extend(-5)
	t.Run("active session", func(t *testing.T) {
		conn.OpenSession(server)

		err := pConn.TimeOut()
		require.NotNil(t, err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.hostname]
		require.True(t, ok, "Connection not found in Pool")

		conn.CloseSession()
	})

	t.Run("connection error", func(t *testing.T) {
		conn.ErrOnConnectionClose(true)

		err := pConn.TimeOut()
		require.NotNil(t, err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection not found in Pool")

		conn.ErrOnConnectionClose(false)
	})

	pConn, err = Pool.Open(&server)
	require.Nil(t, err, "Pool.Open() returned an error: ", err)
	pConn.Extend(-5)
	t.Run("expired", func(t *testing.T) {
		err := pConn.TimeOut()
		require.Nil(t, err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.hostname]
		require.False(t, ok, "Connection found in Pool")
	})
}
