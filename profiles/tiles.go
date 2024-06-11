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
	hideCmd     bool   // Whether or not to send the cmd value to the client.
	hideExp     bool   // Whether or not to send the exp value to the client.
	name        string // Tile name. ("Ping", "Check Connectivity", etc.)
	cmd         string // Command to run on a remote server.
	exp         string // String to match with the results of cmd.
	displaySize int    // Size is a multiple of the smallest button size. Default 40.
}

// DefaultTile creates a new Tile object with several default settings.
func DefaultTile() Tile {
	return Tile{
		hideCmd:     false,
		hideExp:     false,
		displaySize: smallestTileSize * defaultSizeMultiplier,
	}
}

// NewTile creates a new Tile object with a name, cmd(command), and exp(expect) string.
func NewTile(name string, cmd string, exp string) Tile {
	t := DefaultTile()
	// INCOMPLETE: Add html safe validation for name, cmd, and exp here.
	t.name = name
	t.cmd = cmd
	t.exp = exp

	return t
}

// HideCmd is used to determine if Tile.cmd can be shown to the client. Tile.cmd will always be
// available to admins so they can verify or update the commands.
func (t *Tile) HideCmd() bool { return t.hideCmd }

// HideExp is used to determine if Tile.exp can be shown to the client. Tile.exp will always be
// available to admins so they can verify or update the expect string.
func (t *Tile) HideExp() bool { return t.hideCmd }

// Name returns the name set for the Tile.
func (t *Tile) Name() string { return t.name }

// Cmd returns the command string to be ran on a server.
func (t *Tile) Cmd() string { return t.cmd }

// Cmd returns the expect string to be matched against the output of the command being ran
// on a server.
func (t *Tile) Exp() string { return t.exp }

// DisplaySize retuns the size in pixels to set the Tile to. Default 40 would create a
// 40 pixel x 40 pixel Tile in the UI.
func (t *Tile) DisplaySize() int { return t.displaySize }

// SetHideCmd sets the Tile.hideCmd field.
func (t *Tile) SetHideCmd(hide bool) { t.hideCmd = hide }

// SetHideExp sets the Tile.hideExp field.
func (t *Tile) SetHideExp(hide bool) { t.hideExp = hide }

// SetName validates and sets Tile.name.
func (t *Tile) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Tile.SetName: name cannot be empty")
	}

	// INCOMPLETE: Add html safe validation for name here.
	t.name = name
	return nil
}

// SetSize takes a multiplier and sets the Tile.DisplaySize to smallestTileSize * multiplier.
// Passing 0 will default to defaultSizeMultiplier.
func (t *Tile) SetSize(multiplier int) {
	if multiplier <= 0 {
		multiplier = defaultSizeMultiplier
	}

	t.displaySize = smallestTileSize * multiplier
}

// SetCmd sets a command to be ran on a server.
func (t *Tile) SetCmd(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return errors.New("profiles.Tile.SetCmd: cmd cannot be empty or whitespace only")
	}

	// INCOMPLETE: Add html safe validation for cmd here.
	t.cmd = cmd
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
	t.exp = exp
	return nil
}
