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

func NewGroup(name string, servers ...connections.Server) Group {
	g := Group{Name: name, Servers: servers}
	g.uniq()
	return g
}

func (g Group) Count() int       { return g.ServerCount() }
func (g Group) ServerCount() int { return len(g.Servers) }

func (g *Group) SetName(name string) error {
	if name == "" {
		return errors.New("could not set group name, name was empty")
	}

	// Add validation later (prevent html, db, or other exploits)
	g.Name = name
	return nil
}

func (g *Group) AddServers(servers ...connections.Server) {
	if len(servers) < 1 {
		return
	}

	g.Servers = append(g.Servers, servers...)
	g.uniq()
}

func (g *Group) Reset() { g.state = 0 }

func (g *Group) Next() (*connections.Server, error) {
	if g.state >= len(g.Servers) {
		g.state = 0
		return nil, ErrEndOfList
	}

	s := &(g.Servers[g.state])
	g.state++
	return s, nil
}

func (g *Group) uniq() {
	var newGroup []connections.Server
	f := make(map[string]bool)

	for _, s := range g.Servers {
		if _, ok := f[s.Name()]; ok {
			continue
		}

		f[s.Name()] = true
		newGroup = append(newGroup, s)
	}

	g.Servers = newGroup
}
