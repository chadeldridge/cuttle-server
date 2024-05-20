package profiles

import "github.com/chadeldridge/cuttle/connections"

type Group struct {
	Name    string
	Servers []connections.Server
}

/*
func NewGroup(name string, servers ...Server) Group {
	return Group{Name: name, Servers: servers}
}
*/
