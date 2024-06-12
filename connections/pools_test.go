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
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("good connector", func(t *testing.T) {
		pConn, err := Pool.Open(&server)
		require.NoError(err, "Pool.Open() returned an error: ", err)
		require.True(time.Now().Before(pConn.killAt), "killAt was before time.Now()")

		conn.isConnected = false
		delete(Pool, server.Hostname)
	})

	t.Run("already open", func(t *testing.T) {
		pConn, err := Pool.Open(&server)
		require.NoError(err, "Pool.Open() returned an error: ", err)

		pConn2, err := Pool.Open(&server)
		require.NoError(err, "Pool.Open() returned an error: ", err)
		require.Equal(pConn, pConn2, "Connection refs were not the same")

		conn.isConnected = false
		delete(Pool, server.Hostname)
	})

	t.Run("empty connector", func(t *testing.T) {
		server.SetConnector(&MockConnector{})
		_, err := Pool.Open(&server)
		require.Error(err, "Pool.Open() did not return an error")
	})

	t.Run("empty hostname", func(t *testing.T) {
		server = Server{Hostname: ""}
		_, err := Pool.Open(&server)
		require.Error(err, "Pool.Open() returned an error: ", err)

		conn.isConnected = false
		delete(Pool, server.Hostname)
	})
}

func TestPoolsGetConnection(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("existing connection", func(t *testing.T) {
		Pool = ConnectionPool{server.Hostname: &Connection{Server: &server}}

		got := Pool.GetConnection(server)
		require.NotNil(got, "Pool.GetConnection() returned nil Connection")

		conn.isConnected = false
		delete(Pool, server.Hostname)
	})

	t.Run("no connection", func(t *testing.T) {
		got := Pool.GetConnection(server)
		require.Nil(got, "Pool.GetConnection() returned a Connection")
	})
}

func TestPoolsExpires(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("empty connection", func(t *testing.T) {
		pConn := Connection{}
		require.True(pConn.Expires().IsZero(), "Connection.killAt was not zero")
	})

	t.Run("valid connection", func(t *testing.T) {
		pConn := Connection{Server: &server, killAt: time.Now().Add(time.Minute * time.Duration(TTL))}
		require.False(pConn.Expires().IsZero(), "Connection.killAt was zero")
	})
}

func TestPoolsExpired(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("empty connection", func(t *testing.T) {
		pConn := Connection{}
		require.True(pConn.Expired(), "Connection was not expired")
	})

	t.Run("not expired", func(t *testing.T) {
		pConn := Connection{Server: &server, killAt: time.Now().Add(time.Minute * time.Duration(TTL))}
		require.False(pConn.Expired(), "Connection was expired")
	})

	t.Run("expired", func(t *testing.T) {
		pConn := Connection{Server: &server, killAt: time.Now().Add(time.Minute * time.Duration(-5))}
		require.True(pConn.Expired(), "Connection was not expired")
	})
}

func TestPoolsExtend(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("empty connection", func(t *testing.T) {
		pConn := Connection{}
		was := pConn.killAt
		pConn.Extend(TTL)
		require.Greater(pConn.killAt, was, "Connection.killAt was not greater than the old time")
	})

	t.Run("positive", func(t *testing.T) {
		was := time.Now().Add(time.Minute * time.Duration(TTL))
		pConn := Connection{Server: &server, killAt: was}
		pConn.Extend(TTL)
		require.Greater(pConn.killAt, was, "Connection.killAt was not greater than the old time")
	})

	t.Run("negative", func(t *testing.T) {
		was := time.Now().Add(time.Minute * time.Duration(TTL))
		pConn := Connection{Server: &server, killAt: was}
		pConn.Extend(-5)
		require.Less(pConn.killAt, was, "Connection.killAt was not greater than the old time")
	})
}

func TestPoolsConnectionClose(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	t.Run("session active", func(t *testing.T) {
		pConn := Connection{Server: &server}
		Pool = ConnectionPool{server.Hostname: &pConn}
		conn.isConnected = true
		conn.hasSession = true

		err := pConn.Close(false)
		require.Error(err, "Connection.Close() did not return an error")
		_, ok := Pool[server.Hostname]
		require.True(ok, "Connection not found after failed Pool.Close()")

		conn.hasSession = false
	})

	t.Run("connection close error", func(t *testing.T) {
		pConn := Connection{Server: &server}
		Pool = ConnectionPool{server.Hostname: &pConn}
		conn.isConnected = true
		conn.connCloseErr = true

		err := pConn.Close(false)
		require.Error(err, "Connection.Close() did not return an error")
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection found after Connection.Close()")

		conn.connCloseErr = false
	})

	t.Run("not in pool", func(t *testing.T) {
		pConn := Connection{Server: &server}
		conn.isConnected = true

		err := pConn.Close(false)
		require.Error(err, "Pool.Close() did not return an error")
		require.Equal(
			"connections.Connection.Close: Connection not found in Pool",
			err.Error(),
			"Did not receive 'not found' error from Connection.Close(): %s", err,
		)
	})

	t.Run("close connection", func(t *testing.T) {
		pConn := Connection{Server: &server}
		Pool = ConnectionPool{server.Hostname: &pConn}
		conn.isConnected = true

		err := pConn.Close(false)
		require.NoError(err, "Connection.Close() returned an error: %s", err)
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection found after close")
	})
}

func TestPoolsConnectionPoolClose(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	pConn := Connection{Server: &server}
	Pool = ConnectionPool{server.Hostname: &pConn}
	conn.isConnected = true

	t.Run("close connection", func(t *testing.T) {
		err := Pool.Close(server.Hostname, false)
		require.NoError(err, "Pool.Close() returned an error: %s", err)
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection found after close")
	})

	t.Run("no connection", func(t *testing.T) {
		err := Pool.Close(server.Hostname, false)
		require.NoError(err, "Pool.Close() returned an error: %s", err)
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection found after close")
	})
}

func createPool(count int) {
	for i := range count {
		var res bytes.Buffer
		var log bytes.Buffer

		conn := MockConnector{user: testUser}
		server := Server{
			Hostname:  fmt.Sprintf("host%d", i),
			Results:   &res,
			Logs:      &log,
			Connector: &conn,
		}

		pConn := Connection{Server: &server}
		Pool[server.Hostname] = &pConn
		conn.isConnected = true
	}
}

func TestPoolsCloseAll(t *testing.T) {
	require := require.New(t)
	hostCount := 5
	createPool(hostCount)

	t.Run("connection error", func(t *testing.T) {
		Pool["host4"].Server.Connector.(*MockConnector).connCloseErr = true

		require.Equal(hostCount, Pool.Count(), "Pool connection count did not matched expected amount")
		err := Pool.CloseAll()
		require.Error(err, "Pool.Close() did not return an error")
		require.Equal(0, Pool.Count(), "Pool.Count() did not match expected count")
	})

	createPool(hostCount)
	t.Run("no errors", func(t *testing.T) {
		err := Pool.CloseAll()
		require.NoError(err, "Pool.Close() returned an error: ", err)
		require.Equal(0, Pool.Count(), "Pool.Count() did not match expected count")
	})
}

func TestPoolsTimeOut(t *testing.T) {
	var res bytes.Buffer
	var log bytes.Buffer

	require := require.New(t)
	conn := MockConnector{user: testUser}
	server := Server{
		Hostname:  testHost,
		Results:   &res,
		Logs:      &log,
		Connector: &conn,
	}

	pConn := Connection{Server: &server, killAt: time.Now().Add(time.Minute * time.Duration(TTL))}
	Pool = ConnectionPool{server.Hostname: &pConn}
	conn.isConnected = true

	t.Run("not expired", func(t *testing.T) {
		err := pConn.TimeOut()
		require.NoError(err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.Hostname]
		require.True(ok, "Connection not found in Pool")
	})

	pConn.killAt = time.Now().Add(time.Minute * time.Duration(-5))

	t.Run("active session", func(t *testing.T) {
		conn.hasSession = true

		err := pConn.TimeOut()
		require.Error(err, "Connection.TimeOut() did not return an error")
		_, ok := Pool[server.Hostname]
		require.True(ok, "Connection not found in Pool")

		conn.hasSession = false
	})

	t.Run("connection error", func(t *testing.T) {
		conn.connCloseErr = true

		err := pConn.TimeOut()
		require.Error(err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection not found in Pool")

		conn.connCloseErr = false
	})

	Pool = ConnectionPool{server.Hostname: &pConn}
	conn.isConnected = true
	t.Run("expired", func(t *testing.T) {
		err := pConn.TimeOut()
		require.NoError(err, "Connection.TimeOut() returned an error")
		_, ok := Pool[server.Hostname]
		require.False(ok, "Connection found in Pool")
	})
}
