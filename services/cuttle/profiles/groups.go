package profiles

import (
	"errors"

	"github.com/chadeldridge/cuttle/connections"
)

var ErrEndOfList = errors.New("end of list")

type Group struct {
	Name    string
	Servers []connections.Server
	state   int
}

// NewGroup creates a new Group object with a name and servers.
func NewGroup(name string, servers ...connections.Server) Group {
	g := Group{Name: name, Servers: servers}
	g.uniq()
	return g
}

// Count returns the number of servers in the Group.Servers array. Shorthand for Group.ServerCount.
func (g Group) Count() int { return len(g.Servers) }

// SetName sets the name in Group.
func (g *Group) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Group.SetName: name was empty")
	}

	// INCOMPLETE: Add validation later (prevent html, db, or other exploits)
	g.Name = name
	return nil
}

// AddServers adds the servers to Group.Servers.
func (g *Group) AddServers(servers ...connections.Server) {
	if len(servers) < 1 {
		return
	}

	g.Servers = append(g.Servers, servers...)
	g.uniq()
}

// Reset sets the Group.state back to 0. This will cause Group.Next() to return the first element.
func (g *Group) Reset() { g.state = 0 }

// Next returns the current element and advances Group.state by 1.
func (g *Group) Next() (*connections.Server, error) {
	// If state is greater than or equal to the number of elements in Group.Servers then we're
	// reached the end of the list.
	if g.state >= len(g.Servers) {
		g.state = 0
		return nil, ErrEndOfList
	}

	// Get a ref to our current server.
	s := &(g.Servers[g.state])

	// Advance Group.state so we get the next server the next time we're ran.
	g.state++
	return s, nil
}

// uniq removes duplicates from Group.Servers.
func (g *Group) uniq() {
	var newGroup []connections.Server
	f := make(map[string]bool)

	for _, s := range g.Servers {
		if _, ok := f[s.Name]; ok {
			continue
		}

		f[s.Name] = true
		newGroup = append(newGroup, s)
	}

	g.Servers = newGroup
}
