package tests

import (
	"testing"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/stretchr/testify/require"
)

var testUser = "bob"

// testPass = "testUserP@ssw0rd"

func TestSSHNewSSHTest(t *testing.T) {
	require := require.New(t)

	t.Run("basic", func(t *testing.T) {
		test := NewSSHTest("Test SSH echo", true, "echo Hello", "Hello")
		require.Equal("Test SSH echo", test.Name, "NewSSHTest() Name is not 'Test SSH echo'")
		require.True(test.MustSucceed, "NewSSHTest() MustSucceed is not true")
		require.Equal("echo Hello", test.Tester.(*SSHTest).Cmd, "NewSSHTest() Cmd is not 'echo Hello'")
		require.Equal("Hello", test.Tester.(*SSHTest).Exp, "NewSSHTest() Exp is not 'Hello'")
		require.True(test.Tester.(*SSHTest).HideCmd, "NewSSHTest() HideCmd is not true")
		require.True(test.Tester.(*SSHTest).HideExp, "NewSSHTest() HideExp is not true")
	})

	t.Run("show cmd", func(t *testing.T) {
		test := NewSSHTest("Test SSH echo", true, "echo Hello", "Hello", TestArg{Key: "hide_cmd", Value: false})
		require.Equal("Test SSH echo", test.Name, "NewSSHTest() Name is not 'Test SSH echo'")
		require.True(test.MustSucceed, "NewSSHTest() MustSucceed is not true")
		require.Equal("echo Hello", test.Tester.(*SSHTest).Cmd, "NewSSHTest() Cmd is not 'echo Hello'")
		require.Equal("Hello", test.Tester.(*SSHTest).Exp, "NewSSHTest() Exp is not 'Hello'")
		require.False(test.Tester.(*SSHTest).HideCmd, "NewSSHTest() HideCmd is not false")
		require.True(test.Tester.(*SSHTest).HideExp, "NewSSHTest() HideExp is not true")
	})

	t.Run("show exp", func(t *testing.T) {
		test := NewSSHTest("Test SSH echo", true, "echo Hello", "Hello", TestArg{Key: "hide_exp", Value: false})
		require.Equal("Test SSH echo", test.Name, "NewSSHTest() Name is not 'Test SSH echo'")
		require.True(test.MustSucceed, "NewSSHTest() MustSucceed is not true")
		require.Equal("echo Hello", test.Tester.(*SSHTest).Cmd, "NewSSHTest() Cmd is not 'echo Hello'")
		require.Equal("Hello", test.Tester.(*SSHTest).Exp, "NewSSHTest() Exp is not 'Hello'")
		require.True(test.Tester.(*SSHTest).HideCmd, "NewSSHTest() HideCmd is not true")
		require.False(test.Tester.(*SSHTest).HideExp, "NewSSHTest() HideExp is not false")
	})
}

func TestSSHTestRun(t *testing.T) {
	require := require.New(t)
	server := testServerSetup(t)
	defer connections.Pool.CloseAll()

	conn, err := connections.NewMockConnector("my connector", testUser)
	require.NoError(err, "connections.NewMockConnector() returned an error: %s", err)
	server.SetConnector(&conn)

	t.Run("pass", func(t *testing.T) {
		test := SSHTest{HideCmd: true, HideExp: true, Cmd: "echo Hello", Exp: "Hello"}
		err := test.Run(server)
		require.NoError(err, "SSHTest.Run() returned an error: %s", err)
	})

	t.Run("fail", func(t *testing.T) {
		test := SSHTest{HideCmd: true, HideExp: true, Cmd: "echo Hello", Exp: "Goodbye"}
		err := test.Run(server)
		require.Error(err, "SSHTest.Run() did not return an error")
	})
}
