package connections

import (
	"errors"
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

func init() { Pool = make(map[string]*Connection) }

func (p ConnectionPool) Open(server *Server) (*Connection, error) {
	if conn, ok := p[server.hostname]; ok {
		return conn, nil
	}

	conn := &Connection{Server: server}
	conn.killAt = time.Now().Add(time.Minute * time.Duration(TTL))
	err := conn.Connector.Open(*server)
	if err != nil {
		return conn, err
	}

	p[server.hostname] = conn
	return conn, nil
}

// GetConnection returns a connection for the server if one exists. Returns nil if no connection is found.
func (p ConnectionPool) GetConnection(server Server) *Connection {
	conn, ok := p[server.hostname]
	if !ok {
		return nil
	}

	return conn
}

func (c *Connection) Close(force bool) error {
	_, ok := Pool[c.hostname]
	if !ok {
		return errors.New("connections.Connection.Close: Connection not found in Pool")
	}

	err := c.Server.Close(force)
	if err != nil {
		return err
	}

	delete(Pool, c.hostname)
	return nil
}

func (p ConnectionPool) CloseAll() {
	for _, c := range p {
		err := c.Close(true)
		if err != nil {
			log.Printf("connections.ConnectionPool.CloseAll: %s", err)
		}

		delete(p, c.hostname)
	}
}

func (c *Connection) TimeOut() {
	now := time.Now()
	if now.Before(c.killAt) {
		return
	}

	err := c.Close(false)
	if err != nil && err != ErrSessionActive {
		log.Printf("connections.Pool.TimeOut: error closing connection: %s", err)
	}
}
