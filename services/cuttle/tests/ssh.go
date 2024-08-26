package tests

import (
	"fmt"
	"strings"
	"time"

	"github.com/chadeldridge/cuttle-server/services/cuttle/connections"
)

type SSHTest struct {
	HideCmd bool   // Whether or not to send the cmd value to the client.
	HideExp bool   // Whether or not to send the exp value to the client.
	Cmd     string // Command to run on a remote server.
	Exp     string // String to match with the results of cmd.
}

// NewSSHTest creates a new SSH test with the given parameters.
// name: The name of the test.
// mustSucceed: If false, the Tile will continue with the test stack if this test fails.
// cmd: The command to run on the server.
// exp: The expected output of the command.
//
// These TestArg will be evaluated:
// "hide_cmd": bool. If true, the cmd will not be sent to the client. Default is true.
// "hide_exp": bool. If true, the exp will not be sent to the client. Default is true.
func NewSSHTest(name string, mustSucceed bool, cmd string, exp string, args ...TestArg) Test {
	return Test{
		Name:        name,
		MustSucceed: mustSucceed,
		Tester: &SSHTest{
			HideCmd: getSSHHideCmd(args),
			HideExp: getSSHHideExp(args),
			Cmd:     cmd,
			Exp:     exp,
		},
	}
}

func getSSHHideCmd(args []TestArg) bool {
	v := FindArg(args, "hide_cmd")
	if v == nil {
		return true
	}

	return v.(bool)
}

func getSSHHideExp(args []TestArg) bool {
	v := FindArg(args, "hide_exp")
	if v == nil {
		return true
	}

	return v.(bool)
}

// Run runs the SSHTest on the given server.
func (t SSHTest) Run(server connections.Server, args ...TestArg) error {
	err := server.Open(connections.SSH)
	if err != nil {
		server.Buffers.Log(time.Now(), fmt.Sprintf("SSHTest.Run: %s", err))
		return ErrTestFailed
	}
	defer server.Close(connections.SSH, false)

	err = server.Run(connections.SSH, t.Cmd, t.Exp)
	if err != nil {
		server.Buffers.Log(time.Now(), fmt.Sprintf("SSHTest.Run: %s", err))
		return ErrTestFailed
	}

	return nil
}

// SetHideCmd sets whether or not to hide SSHTest.Cmd from non-admin users.
func (t *SSHTest) SetHideCmd(hide bool) { t.HideCmd = hide }

// SetHideExp sets whether or not to hide SSHTest.Exp from non-admin users.
func (t *SSHTest) SetHideExp(hide bool) { t.HideExp = hide }

// SetCmd sets a command to be ran on a server.
func (t *SSHTest) SetCmd(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return fmt.Errorf("profiles.Tile.SetCmd: cmd cannot be empty or whitespace only")
	}

	// INCOMPLETE: Add html safe validation for cmd here.
	t.Cmd = cmd
	return nil
}

// SetExp sets the expect string which will be matches against the output of Tile.cmd after being
// ran on a server.
func (t *SSHTest) SetExp(exp string) error {
	exp = strings.TrimSpace(exp)
	if exp == "" {
		return fmt.Errorf("profiles.Tile.SetExp: exp cannot be empty or whitespace only")
	}

	// INCOMPLETE: Add html safe validation for exp here.
	t.Exp = exp
	return nil
}
