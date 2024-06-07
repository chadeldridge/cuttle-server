package profiles

import (
	"errors"
	"fmt"

	"github.com/chadeldridge/cuttle/connections"
)

// Profile holds the Groups and command Tiles uses to run tests.
type Profile struct {
	Name   string
	Tiles  map[string]Tile  // List of command Tiles that can be run against these server groups.
	Groups map[string]Group // List of groups to test against.
}

// NewProfile creates a new Profile object with a display Name and at least one Group.
func NewProfile(name string, groups ...Group) Profile {
	// Make our maps so we don't get errors later.
	p := Profile{
		Tiles:  make(map[string]Tile),
		Groups: make(map[string]Group),
	}

	p.SetName(name)
	// If ther are no groups then there's no reason to do anything.
	if len(groups) > 0 {
		p.AddGroups(groups...)
	}

	return p
}

// SetName sets the Profile.Name after validating it is safe.
func (p *Profile) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Profile.SetName: name was empty")
	}

	// INCOMPLETE: Add html safe name validation here.
	p.Name = name
	return nil
}

// AddTiles adds the list of Tiles to Group.Tiles. Only new Tiles will be added.
func (p *Profile) AddTiles(tiles ...Tile) error {
	// Providing no tiles is probably unintended behaviour so we error.
	if len(tiles) < 1 {
		return errors.New("profiles.Profile.AddTiles: no tiles provided")
	}

	var errs error
	for _, tile := range tiles {
		// If the Tile already exists, skip to the next one.
		if _, ok := p.Tiles[tile.Name()]; ok {
			errs = errors.Join(errs,
				fmt.Errorf("profiles.Profile.AddTiles: tile already exists for '%s'", tile.Name()),
			)
			continue
		}

		p.Tiles[tile.Name()] = tile
	}

	return errs
}

// AddGroups adds a list of Group to Profile.Groups.
func (p *Profile) AddGroups(groups ...Group) error {
	// Providing no tiles is probably unintended behaviour so we error.
	if len(groups) < 1 {
		return errors.New("profiles.Profile.AddGroups: no groups provided")
	}

	var errs error
	for _, group := range groups {
		// If the Group already exists, skip to the next one.
		if _, ok := p.Groups[group.Name]; ok {
			errs = errors.Join(
				errs,
				fmt.Errorf("profiles.Profile.AddGroups: group already exists for '%s'", group.Name),
			)
			continue
		}

		p.Groups[group.Name] = group
	}

	return errs
}

// GetTile retrieves the Tile by name from Profile.Tiles.
func (p Profile) GetTile(name string) (Tile, error) {
	var t Tile
	if name == "" {
		return t, errors.New("profiles.Profile.GetTile: tileName was empty")
	}

	t, ok := p.Tiles[name]
	if !ok {
		return t, errors.New("profiles.Profile.GetTile: tile not found")
	}

	return t, nil
}

// GetGroup retrieves the Group by name from Profile.Groups.
func (p Profile) GetGroup(name string) (Group, error) {
	var g Group
	if name == "" {
		return g, errors.New("profiles.Profile.GetGroup: name was empty")
	}

	g, ok := p.Groups[name]
	if !ok {
		return g, errors.New("profiles.Profile.GetGroup: group not found")
	}

	return g, nil
}

// Execute runs the Tile command against each server in the selected group. Execute also replaces
// special variables in the command and expect with the appropriate values.
func (p Profile) Execute(tileName, groupName string) error {
	tile, err := p.GetTile(tileName)
	if err != nil {
		return fmt.Errorf("profiles.Profile.Execute: %s", err)
	}

	group, err := p.GetGroup(groupName)
	if err != nil {
		return fmt.Errorf("profiles.Profile.Execute: %s", err)
	}

	var errs error
	for {
		server, err := group.Next()
		if err != nil {
			if err == ErrEndOfList {
				break
			}

			errs = errors.Join(errs, err)
			return errs
		}

		// INCOMPLETE: Add special variable replacement in the command and expect strings.

		// Make sure we have an open connection to the server.
		_, err = connections.Pool.Open(server)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		err = server.Run(tile.cmd, tile.exp)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
	}

	return errs
}
