package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/chadeldridge/cuttle/core"
	"github.com/stretchr/testify/require"
)

type testFlags struct {
	appname string
	flags   map[string]string
	want    map[string]string
	args    []string
}

var flagTests = []testFlags{
	{
		appname: "app",
		flags:   map[string]string{"-v": "", "-h": "127.0.0.1", "-p": "9090", "-d": "/tmp", "-c": "/tmp/config.yaml"},
		want: map[string]string{
			"debug":       "true",
			"api_host":    "127.0.0.1",
			"api_port":    "9090",
			"db_root":     "/tmp",
			"config_file": "/tmp/config.yaml",
		},
		args: []string{},
	},
	{
		appname: "app",
		flags: map[string]string{
			"--verbose":     "",
			"--host":        "127.0.0.1",
			"--port":        "9090",
			"--db-root":     "/tmp",
			"--config-file": "/tmp/config.yaml",
		},
		want: map[string]string{
			"debug":       "true",
			"api_host":    "127.0.0.1",
			"api_port":    "9090",
			"db_root":     "/tmp",
			"config_file": "/tmp/config.yaml",
		},
		args: []string{},
	},
}

func (f testFlags) toArgs() []string {
	a := []string{f.appname}
	for k, v := range f.flags {
		if v == "" {
			a = append(a, k)
			continue
		}

		a = append(a, k, v)
	}

	return append(a, f.args...)
}

func TestFlagsParseFlags(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	logger := core.NewLogger(&buf, "test: ", log.LstdFlags, false)

	t.Run("flag tests", func(t *testing.T) {
		for _, tf := range flagTests {
			log.Printf("args: %v\n", tf.toArgs())
			flags, args := parseFlags(logger, tf.toArgs())
			log.Printf("flags: %v\n", flags)
			require.Equal(tf.want, flags, "parseFlags() flags did not match")
			require.Equal(tf.args, args, "parseFlags() args did not match")
		}
	})

	t.Run("help", func(t *testing.T) {
		flags, args := parseFlags(logger, []string{"app", "--help"})
		require.Nil(flags, "parseFlags() flags did not match")
		require.Nil(args, "parseFlags() args did not match")
		require.NotEmpty(buf.String(), "parseFlags() help did not print")
		require.Contains(buf.String(), "Usage:", "parseFlags() help did not match")
		buf.Reset()
	})

	t.Run("version", func(t *testing.T) {
		flags, args := parseFlags(logger, []string{"app", "--version"})
		require.Nil(flags, "parseFlags() flags did not match")
		require.Nil(args, "parseFlags() args did not match")
		require.Contains(buf.String(), "Cuttle v", "parseFlags() version did not match")
		buf.Reset()
	})
}
