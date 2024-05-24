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
	smallestTileSize      = 20
	defaultSizeMultiplier = 2
)

type Tile struct {
	hideCmd     bool   // Whether or not to send the cmd value to the client.
	hideExp     bool   // Whether or not to send the exp value to the client.
	name        string // Tile name. ("Ping", "Check Connectivity", etc.)
	cmd         string // Command to run on a remote server.
	exp         string // String to match with the results of cmd.
	displaySize int    // Size is a multiple of the smallest button size. Default = 2
}

func DefaultTile() Tile {
	return Tile{
		hideCmd:     false,
		hideExp:     false,
		displaySize: smallestTileSize * defaultSizeMultiplier,
	}
}

func NewTile(name string, cmd string) Tile {
	t := DefaultTile()
	t.name = name
	t.cmd = cmd

	return t
}

func (t *Tile) HideCmd() bool    { return t.hideCmd }
func (t *Tile) HideExp() bool    { return t.hideCmd }
func (t *Tile) Name() string     { return t.name }
func (t *Tile) Cmd() string      { return t.cmd }
func (t *Tile) Exp() string      { return t.exp }
func (t *Tile) DisplaySize() int { return t.displaySize }

func (t *Tile) SetHideCmd(hide bool) { t.hideCmd = hide }
func (t *Tile) SetHideExp(hide bool) { t.hideCmd = hide }

func (t *Tile) SetName(name string) {
	// Add validation here
	t.name = name
}

// SetSize takes a multiplier and sets the Tile.DisplaySize to smallestTileSize * multiplier.
// Passing 0 will default to defaultSizeMultiplier.
func (t *Tile) SetSize(multiplier int) {
	if multiplier == 0 {
		multiplier = defaultSizeMultiplier
	}

	t.displaySize = smallestTileSize * multiplier
}

func (t *Tile) SetCmd(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return errors.New("tile.SetCmd: cmd cannot be empty or whitespace only")
	}

	// Add command validation here.
	t.cmd = cmd
	return nil
}

func (t *Tile) SetExp(exp string) error {
	exp = strings.TrimSpace(exp)
	if exp == "" {
		return errors.New("tile.SetExp: exp cannot be empty or whitespace only")
	}

	// Add exp validation here.
	t.exp = exp
	return nil
}
