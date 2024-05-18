package cuttle

type Group struct {
	Name    string
	Servers []Server
}

func NewGroup(name string, servers ...Server) Group {
	return Group{Name: name, Servers: servers}
}
