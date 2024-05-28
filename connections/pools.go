package connections

import (
	"log"
	"time"
)

var TTL int // Time To Live in number of minutes

type Pool struct {
	Connections map[string]*Connection // map[Open.hostame]Open
}

type Connection struct {
	*Server
	killAt time.Time
}

func NewPool() Pool {
	return Pool{Connections: make(map[string]*Connection)}
}

func (p Pool) Open(server *Server) (*Connection, error) {
	if conn, ok := p.Connections[server.hostname]; ok {
		return conn, nil
	}

	conn := &Connection{Server: server}
	conn.killAt = time.Now().Add(time.Minute * time.Duration(TTL))
	err := conn.Connector.Open(*server)
	if err != nil {
		return conn, err
	}

	p.Connections[server.hostname] = conn
	return conn, nil
}

// GetConnection returns a connection for the server if one exists. Returns nil if no connection is found.
func (p Pool) GetConnection(server Server) *Connection {
	conn, ok := p.Connections[server.hostname]
	if !ok {
		return nil
	}

	return conn
}

func (p Pool) Close(conn *Connection) {
	c, ok := p.Connections[conn.hostname]
	if !ok {
		return
	}

	if c != conn {
		log.Printf("connections.Pool.Close: found multiple connections for %s, closing both", conn.hostname)
		c.Close()
	}

	conn.Close()
	delete(p.Connections, conn.hostname)
}

func (p Pool) CloseAll(server Server) {
	for _, conn := range p.Connections {
		conn.Close()
		delete(p.Connections, conn.hostname)
	}
}

func (c *Connection) Close() { c.Server.Connector.Close() }

func (p *Pool) TimeOut(conn *Connection) {
	now := time.Now()
	if now.Before(conn.killAt) {
		return
	}

	p.Close(conn)
}
