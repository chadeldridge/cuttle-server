package connections

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	Pool ConnectionPool // Our shared connection pool
	TTL  int            // Time To Live in number of minutes
)

// ConnectionPool holds an array of connections used to setup our shared pool.
type ConnectionPool map[string]*Connection // map[Open.hostame]Open

// Connection holds a Server ref and our time to kill for connection cleanup.
type Connection struct {
	*Server
	killAt time.Time
}

// Always make sure we have an allocated Pool we can actually work with and set a default TTL.
func init() {
	Pool = make(map[string]*Connection)
	TTL = 2 // Two minute default TTL
}

func (p ConnectionPool) Count() int { return len(p) }

// Open creates a new Connection and adds it to Pool.
func (p ConnectionPool) Open(server *Server) (*Connection, error) {
	conn := &Connection{Server: server}
	if server.hostname == "" {
		return conn, errors.New("connections.Pool.Open: hostname was empty")
	}

	err := conn.Open(p)
	return conn, err
}

func (c *Connection) Open(pool ConnectionPool) error {
	if conn, ok := pool[c.Server.hostname]; ok {
		c = conn
	}

	c.killAt = time.Now().Add(time.Minute * time.Duration(TTL))
	err := c.Connector.Open(*c.Server)
	if err != nil {
		return err
	}

	pool[c.Server.hostname] = c
	return nil
}

// GetConnection returns a connection for the server if one exists. Returns nil if no connection is found.
func (p ConnectionPool) GetConnection(server Server) *Connection {
	conn, ok := p[server.hostname]
	if !ok {
		return nil
	}

	return conn
}

// Expires returns the Connection.killAt time.
func (c Connection) Expires() time.Time { return c.killAt }

// Expired returns true if it is currently past the Connection.killAt time.
func (c Connection) Expired() bool { return c.killAt.Before(time.Now()) }

// Extend add the specified number of minutes to the killAt time.
func (c *Connection) Extend(minutes int) {
	c.killAt = c.killAt.Add(time.Minute * time.Duration(minutes))
}

// Close closes the connection and removes it from the Pool. If the connection is not in the
// Pool, Close will return an error and will NOT try to close the connection.
func (c *Connection) Close(force bool) error {
	_, ok := Pool[c.hostname]
	if !ok {
		return errors.New("connections.Connection.Close: Connection not found in Pool")
	}

	err := c.Server.Close(force)
	log.Printf("c.Server.Close(false): %s", err)
	if err != nil && err == ErrSessionActive {
		return err
	}

	delete(Pool, c.hostname)
	return err
}

// CloseAll will force close all connections in the ConnectionPool. This means it will try to close
// the connection if it it has an active session.
func (p ConnectionPool) CloseAll() error {
	var errs error
	for _, c := range p {
		err := c.Close(true)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("connections.ConnectionPool.CloseAll: %s", err))
		}

		delete(p, c.hostname)
	}

	return errs
}

// TimeOut checks the connection to see if it is passed its killAt time. If so it will attempt to
// close the connection. If a connection is active TimeOut will extend the killAt time by the TTL.
func (c *Connection) TimeOut() error {
	if !c.Expired() {
		return nil
	}

	err := c.Close(false)
	if err != nil {
		if err != ErrSessionActive {
			return fmt.Errorf(
				"connections.Pool.TimeOut: error closing connection %s: %s",
				c.hostname, err,
			)
		}

		// If the connection was active, extend the time by the TTL.
		c.Extend(TTL)
		return err
	}

	return nil
}
