package profiles

import (
	"errors"
	"fmt"
	"time"

	"github.com/chadeldridge/cuttle-server/services/cuttle/connections"
	"github.com/chadeldridge/cuttle-server/services/cuttle/tests"
)

// Command Tiles

// Command Variables
// {{Group(groupName)}}		"server1, server2, server3"...
// {{Server(serverName)}}	"server1"
// {{IP(serverName)}}		"192.168.1.1"
// {{IPs(groupName)}}		"192.168.1.1, 192.168.1.2, 192.168.1.3"...

const (
	SmallestTileSize      = 20 // Size in pixels.
	DefaultSizeMultiplier = 2  // Default multiplier to determine DisplaySize.
)

type Tile struct {
	Name        string       // Tile name. ("Ping", "Check Connectivity", etc.)
	DisplaySize int          // Size is a multiple of the smallest button size. Default 40.
	Tests       []tests.Test // List of tests to run.
	AllMustPass bool         // If true, all tests must pass for the tile to pass. Ignores Test.MustSucceed.
	InParallel  bool         // If true, all tests will be ran in parallel.
}

// DefaultTile creates a new Tile object with several default settings.
func DefaultTile() Tile {
	return Tile{
		DisplaySize: SmallestTileSize * DefaultSizeMultiplier,
		AllMustPass: false,
		InParallel:  false,
	}
}

// NewTile creates a new Tile object with a name and the given tests. Tests will be ran in order
// unless InParallel is set.
func NewTile(name string, tileTests ...tests.Test) Tile {
	t := DefaultTile()
	// INCOMPLETE: Add html safe validation for name, cmd, and exp here.
	t.SetName(name)
	t.Tests = tileTests

	return t
}

// SetName validates and sets Tile.name.
func (t *Tile) SetName(name string) error {
	if name == "" {
		return errors.New("profiles.Tile.SetName: name cannot be empty")
	}

	// INCOMPLETE: Add html safe validation for name here.
	t.Name = name
	return nil
}

// SetSize takes a multiplier and sets the Tile.DisplaySize to SmallestTileSize * multiplier.
// Passing 0 will default to DefaultSizeMultiplier.
func (t *Tile) SetSize(multiplier int) {
	if multiplier <= 0 {
		multiplier = DefaultSizeMultiplier
	}

	t.DisplaySize = SmallestTileSize * multiplier
}

// AddTest adds a test to the Tile.Tests slice at the given postion.
func (t *Tile) AddTest(position int, newTest tests.Test) {
	if position == 0 {
		t.Tests = append([]tests.Test{newTest}, t.Tests...)
		return
	}

	if position >= len(t.Tests) {
		t.Tests = append(t.Tests, newTest)
		return
	}

	t.Tests = append(t.Tests[:position], append([]tests.Test{newTest}, t.Tests[position:]...)...)
}

// AppendTest adds a test to the end of the Tile.Tests slice.
func (t *Tile) AppendTest(newTest tests.Test) { t.Tests = append(t.Tests, newTest) }

// AddTests adds a list of tests to the Tile.Tests slice.
func (t *Tile) AddTests(newTests ...tests.Test) { t.Tests = append(t.Tests, newTests...) }

// RemoveTest removes a test from the Tile.Tests slice.
func (t *Tile) RemoveTest(test tests.Test) {
	for i, tileTest := range t.Tests {
		if tileTest == test {
			t.RemoveTestAt(i)
			return
		}
	}
}

// RemoveTestAt removes a test from the Tile.Tests slice at the given position.
func (t *Tile) RemoveTestAt(position int) {
	if position == 0 {
		t.Tests = t.Tests[1:]
		return
	}

	if position == len(t.Tests)-1 {
		t.Tests = t.Tests[:len(t.Tests)-1]
		return
	}

	t.Tests = append(t.Tests[:position], t.Tests[position+1:]...)
}

// RunInParallel sets the Tile.InParallel field to true.
func (t *Tile) RunInParallel() { t.InParallel = true }

// RunInSequence sets the Tile.InParallel field to false.
func (t *Tile) RunInSequence() { t.InParallel = false }

// Run tests in the Tile.Tests slice and return an error if any tests fail.
func (t Tile) Run(server connections.Server, args ...tests.TestArg) error {
	if t.InParallel {
		return t.runInParallel(server, args)
	}

	return t.runInSequence(server, args)
}

func (t Tile) runInSequence(server connections.Server, args []tests.TestArg) error {
	for _, test := range t.Tests {
		err := test.Run(server, args...)
		if err != nil && test.MustSucceed {
			server.Buffers.PrintResults(
				time.Now(),
				fmt.Sprintf("(%s) %s - %s...fail", t.Name, test.Name, server.Hostname),
				err,
			)
			return err
		}

		server.Buffers.PrintResults(
			time.Now(),
			fmt.Sprintf("(%s) %s - %s...pass", t.Name, test.Name, server.Hostname),
			nil,
		)
	}

	server.Buffers.PrintResults(time.Now(), fmt.Sprintf("(%s) %s...pass", t.Name, server.Hostname), nil)
	return nil
}

func (t Tile) runInParallel(server connections.Server, args []tests.TestArg) error {
	errs := make(chan error, len(t.Tests))
	for _, test := range t.Tests {
		go func(test tests.Test) {
			errs <- test.Run(server, args...)
		}(test)
	}

	for i := 0; i < len(t.Tests); i++ {
		err := <-errs
		if err != nil && t.AllMustPass {
			server.Buffers.PrintResults(
				time.Now(),
				fmt.Sprintf("(%s) %s - %s...failed", t.Name, t.Tests[i].Name, server.Hostname),
				err,
			)
			return err
		}

		server.Buffers.PrintResults(
			time.Now(),
			fmt.Sprintf("(%s) %s - %s...pass", t.Name, t.Tests[i].Name, server.Hostname),
			nil,
		)
	}

	server.Buffers.PrintResults(time.Now(), fmt.Sprintf("(%s) %s...pass", t.Name, server.Hostname), nil)
	return nil
}
