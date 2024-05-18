package cuttle

// Command Tiles

// Command Variables
// {{Group(groupName)}}		"server1, server2, server3"...
// {{Server(serverName)}}	"server1"
// {{IP(serverName)}}		"192.168.1.1"
// {{IPs(groupName)}}		"192.168.1.1, 192.168.1.2, 192.168.1.3"...

const (
	defaultTileSize = 2
)

type Tile struct {
	Name    string
	Command string
	Size    int // Size is a multiple of the smallest button size. Default = 2
}

func DefaultTile() Tile {
	return Tile{Size: defaultTileSize}
}

func NewTile(name string, cmd string) Tile {
	t := DefaultTile()
	t.Name = name
	t.Command = cmd

	return t
}
