package profiles

import (
	"errors"
	"fmt"
	"log"

	"github.com/chadeldridge/cuttle/connections"
)

type Profile struct {
	Name   string
	Tiles  map[string]Tile
	Groups map[string]Group
}

func NewProfile(name string, groups ...Group) Profile {
	p := Profile{
		Tiles:  make(map[string]Tile, 0),
		Groups: make(map[string]Group),
	}

	p.SetName(name)
	if len(groups) > 0 {
		p.AddGroups(groups...)
	}

	return p
}

func (p *Profile) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Profile.SetName: name was empty")
	}

	// Add validation here
	p.Name = name
	return nil
}

func (p *Profile) AddTiles(tiles ...Tile) error {
	if len(tiles) < 1 {
		return errors.New("profiles.Profile.AddTiles: no tiles provided")
	}

	for _, tile := range tiles {
		if _, ok := p.Tiles[tile.Name()]; ok {
			continue
		}

		p.Tiles[tile.Name()] = tile
	}

	return nil
}

func (p *Profile) AddGroups(groups ...Group) error {
	if len(groups) < 1 {
		return errors.New("profiles.Profile.AddGroups: no groups provided")
	}

	for _, group := range groups {
		if _, ok := p.Groups[group.Name]; ok {
			continue
		}

		p.Groups[group.Name] = group
	}

	return nil
}

func (p Profile) Execute(pool connections.Pool, tileName, groupName string) error {
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

			return err
		}

		// Make sure we have an open connection to the server.
		_, err = pool.Open(server)
		if err != nil {
			log.Printf("profiles.Profile.Execute: failed to open connection to %s: %s", server.Hostname(), err)
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

func (p Profile) GetGroup(name string) (Group, error) {
	var g Group
	if name == "" {
		return g, errors.New("profiles.GetGroup: name was empty")
	}

	g, ok := p.Groups[name]
	if !ok {
		return g, errors.New("profiles.GetGroup: group not found")
	}

	return g, nil
}
