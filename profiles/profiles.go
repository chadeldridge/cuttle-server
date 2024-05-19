package profiles

type Profile struct {
	Name   string
	Groups []Group
}

func NewProfile(name string, groups ...Group) Profile {
	return Profile{Name: name, Groups: groups}
}
