package profiles

type Profile struct {
	Name   string
	Groups []Group
	Tiles  []Tile
}

func NewProfile(name string, groups ...Group) Profile {
	return Profile{Name: name, Groups: groups, Tiles: make([]Tile, 0)}
}
