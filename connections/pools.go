package connections

import (
	"log"
	"time"
)

var (
	Pool ConnectionPool
	TTL  int // Time To Live in number of minutes
)

type ConnectionPool map[string]*Connection // map[Open.hostame]Open

type Connection struct {
	*Server
	killAt time.Time
}

func init() { Pool = make(map[string]*Connection) }

// func init()                   { Pool = NewPool() }
// func NewPool() ConnectionPool { return make(map[string]*Connection) }

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

func (p ConnectionPool) Close(conn *Connection) {
	c, ok := p[conn.hostname]
	if !ok {
		return
	}

	if c != conn {
		log.Printf("connections.Pool.Close: found multiple connections for %s, closing both", conn.hostname)
		c.Close()
	}

	conn.Close()
	delete(p, conn.hostname)
}

func (p ConnectionPool) CloseAll() {
	for _, c := range p {
		c.Close()
		delete(p, c.hostname)
	}
}

func (c *Connection) Close() { c.Server.Connector.Close() }

func (p *ConnectionPool) TimeOut(conn *Connection) {
	now := time.Now()
	if now.Before(conn.killAt) {
		return
	}

	p.Close(conn)
}
