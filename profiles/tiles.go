package profiles

import (
	"errors"
	"strings"
)

// Command Tiles

// Command Variables
// {{Group(groupName)}}		"server1, server2, server3"...
// {{Server(serverName)}}	"server1"
// {{IP(serverName)}}		"192.168.1.1"
// {{IPs(groupName)}}		"192.168.1.1, 192.168.1.2, 192.168.1.3"...

const (
	smallestTileSize      = 20 // Size in pixels.
	defaultSizeMultiplier = 2  // Default multiplier to determine DisplaySize.
)

type Tile struct {
	HideCmd     bool   // Whether or not to send the cmd value to the client.
	HideExp     bool   // Whether or not to send the exp value to the client.
	Name        string // Tile name. ("Ping", "Check Connectivity", etc.)
	Cmd         string // Command to run on a remote server.
	Exp         string // String to match with the results of cmd.
	DisplaySize int    // Size is a multiple of the smallest button size. Default 40.
}

// DefaultTile creates a new Tile object with several default settings.
func DefaultTile() Tile {
	return Tile{
		HideCmd:     false,
		HideExp:     false,
		DisplaySize: smallestTileSize * defaultSizeMultiplier,
	}
}

// NewTile creates a new Tile object with a name, cmd(command), and exp(expect) string.
func NewTile(name string, cmd string, exp string) Tile {
	t := DefaultTile()
	// INCOMPLETE: Add html safe validation for name, cmd, and exp here.
	t.Name = name
	t.Cmd = cmd
	t.Exp = exp

	return t
}

// SetHideCmd sets the Tile.hideCmd field.
func (t *Tile) SetHideCmd(hide bool) { t.HideCmd = hide }

// SetHideExp sets the Tile.hideExp field.
func (t *Tile) SetHideExp(hide bool) { t.HideExp = hide }

// SetName validates and sets Tile.name.
func (t *Tile) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Tile.SetName: name cannot be empty")
	}

	// INCOMPLETE: Add html safe validation for name here.
	t.Name = name
	return nil
}

// SetSize takes a multiplier and sets the Tile.DisplaySize to smallestTileSize * multiplier.
// Passing 0 will default to defaultSizeMultiplier.
func (t *Tile) SetSize(multiplier int) {
	if multiplier <= 0 {
		multiplier = defaultSizeMultiplier
	}

	t.DisplaySize = smallestTileSize * multiplier
}

// SetCmd sets a command to be ran on a server.
func (t *Tile) SetCmd(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return errors.New("profiles.Tile.SetCmd: cmd cannot be empty or whitespace only")
	}

	// INCOMPLETE: Add html safe validation for cmd here.
	t.Cmd = cmd
	return nil
}

// SetExp sets the expect string which will be matches against the output of Tile.cmd after being
// ran on a server.
func (t *Tile) SetExp(exp string) error {
	exp = strings.TrimSpace(exp)
	if exp == "" {
		return errors.New("profiles.Tile.SetExp: exp cannot be empty or whitespace only")
	}

	// INCOMPLETE: Add html safe validation for exp here.
	t.Exp = exp
	return nil
}
