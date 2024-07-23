package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func defaultConfig() *Config {
	return &Config{
		Env:             "dev",
		Debug:           false,
		APIHost:         DefaultAPIHost,
		APIPort:         DefaultAPIPort,
		DBRoot:          DefaultDBRoot,
		ShutdownTimeout: DefaultShutdownTimeout,
	}
}

func TestConfigNewConfig(t *testing.T) {
	require := require.New(t)
	want := map[string]string{
		"api_host": "127.0.0.1",
		"api_port": "9090",
		"db_root":  "/tmp/db",
		"debug":    "true",
		"env":      "prod",
	}

	t.Run("no flags", func(t *testing.T) {
		flags := map[string]string{}
		args := []string{}
		env := map[string]string{}

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		require.Equal(defaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("flags only", func(t *testing.T) {
		flags := map[string]string{
			"api_host": "127.0.0.1",
			"api_port": "9090",
			"db_root":  "/tmp/db",
			"debug":    "true",
			"env":      "prod",
		}
		args := []string{}
		env := map[string]string{}

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
	})

	t.Run("env only", func(t *testing.T) {
		flags := map[string]string{}
		env := map[string]string{
			"CUTTLE_API_HOST": "127.0.0.1",
			"CUTTLE_API_PORT": "9090",
			"CUTTLE_DB_ROOT":  "/tmp/db",
			"CUTTLE_DEBUG":    "true",
			"CUTTLE_ENV":      "prod",
		}
		args := []string{}

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
	})

	t.Run("from file", func(t *testing.T) {
		flags := map[string]string{"config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{}
		args := []string{}
		config := &Config{
			Env:     want["env"],
			Debug:   true,
			APIHost: want["api_host"],
			APIPort: want["api_port"],
			DBRoot:  want["db_root"],
		}

		file := "/tmp/cuttle.yaml"
		writeFile(file, config)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
		removeFile(file)
	})

	t.Run("missing file", func(t *testing.T) {
		flags := map[string]string{"config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{}
		args := []string{}

		file := "/tmp/cuttle.yaml"
		removeFile(file)

		c, err := NewConfig(flags, args, env)
		require.Error(err, "NewConfig() did not return an error")
		require.Equal(
			"stat /tmp/cuttle.yaml: no such file or directory",
			err.Error(),
			"NewConfig() did not return the correct error",
		)
		require.Equal(defaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("invalid option", func(t *testing.T) {
		flags := map[string]string{"invalid": "value"}
		env := map[string]string{}
		args := []string{}

		c, err := NewConfig(flags, args, env)
		require.Error(err, "NewConfig() did not return an error")
		require.Equal(ErrUnknownOpt, err, "NewConfig() did not return the correct error")
		require.Equal(defaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("override", func(t *testing.T) {
		flags := map[string]string{"api_host": want["api_host"], "config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{"CUTTLE_API_HOST": "63.57.123.41", "CUTTLE_API_PORT": want["api_port"]}
		args := []string{}
		config := &Config{
			Env:     want["env"],
			Debug:   true,
			APIHost: "192.168.0.1",
			APIPort: "3000",
			DBRoot:  want["db_root"],
		}

		file := "/tmp/cuttle.yaml"
		writeFile(file, config)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
		removeFile(file)
	})
}

func testConfig(t *testing.T, want map[string]string, c *Config) {
	require := require.New(t)
	require.Equal(want["api_host"], c.APIHost, "NewConfig() did not set the APIHost")
	require.Equal(want["api_port"], c.APIPort, "NewConfig() did not set the APIPort")
	require.Equal(want["db_root"], c.DBRoot, "NewConfig() did not set the DBRoot")
	require.True(c.Debug, "NewConfig() did not set the debug")
	require.Equal(want["env"], c.Env, "NewConfig() did not set the env")
}

func TestConfigValidateEnv(t *testing.T) {
	require := require.New(t)

	for _, env := range supportedEnv {
		t.Run(env, func(t *testing.T) {
			ok := validateEnv(env)
			require.True(ok, "validateEnv() returned false")
		})
	}

	t.Run("invalid", func(t *testing.T) {
		ok := validateEnv("invalid")
		require.False(ok, "validateEnv() returned true")
	})
}

func TestConfigSetConfigValue(t *testing.T) {
	require := require.New(t)
	c := defaultConfig()

	t.Run("set one", func(t *testing.T) {
		err := setConfigValue(c, "api_host", "127.0.0.1")
		require.NoError(err, "setConfigValue() returned an error")
		require.Equal("127.0.0.1", c.APIHost, "setConfigValue() did not set the value")
	})

	t.Run("set all", func(t *testing.T) {
		m := map[string]string{
			"api_host": "127.0.0.1",
			"api_port": "9090",
			"db_root":  "/tmp/db",
			"debug":    "true",
			"env":      "prod",
		}
		tc := Config{
			Env:             "prod",
			Debug:           true,
			APIHost:         "127.0.0.1",
			APIPort:         "9090",
			DBRoot:          "/tmp/db",
			ShutdownTimeout: DefaultShutdownTimeout,
		}

		for k, v := range m {
			err := setConfigValue(c, k, v)
			require.NoError(err, "setConfigValue() returned an error")
		}

		require.Equal(tc, *c, "setConfigValue() did not set the value")
	})

	t.Run("invalid key", func(t *testing.T) {
		err := setConfigValue(c, "invalid", "value")
		require.Error(err, "setConfigValue() did not return an error")
		require.Equal(ErrUnknownOpt, err, "setConfigValue() did not return the correct error")
	})

	t.Run("invalid env", func(t *testing.T) {
		err := setConfigValue(c, "env", "invalid")
		require.Error(err, "setConfigValue() did not return an error")
		require.Equal(ErrInvalidEnv, err, "setConfigValue() did not return the correct error")
	})

	t.Run("debug false", func(t *testing.T) {
		err := setConfigValue(c, "debug", "false")
		require.NoError(err, "setConfigValue() returned an error")
		require.False(c.Debug, "setConfigValue() did not set the value")
	})
}

func TestConfigParseEnvVars(t *testing.T) {
	require := require.New(t)
	want := map[string]string{
		"api_host": "127.0.0.1",
		"api_port": "9090",
		"db_root":  "/tmp/db",
		"debug":    "true",
		"env":      "prod",
	}

	t.Run("flags only", func(t *testing.T) {
		tf := map[string]string{
			"api_host": "127.0.0.1",
			"api_port": "9090",
			"db_root":  "/tmp/db",
			"debug":    "true",
			"env":      "prod",
		}
		env := map[string]string{}

		flags := parseEnvVars(tf, env)
		for k, v := range want {
			require.Equal(v, flags[k], "parseEnvVars() did not set the value")
		}
	})

	t.Run("env only", func(t *testing.T) {
		tf := map[string]string{}
		env := map[string]string{
			"random":          "value",
			"CUTTLE_API_HOST": "127.0.0.1",
			"CUTTLE_API_PORT": "9090",
			"CUTTLE_DB_ROOT":  "/tmp/db",
			"CUTTLE_DEBUG":    "true",
			"CUTTLE_ENV":      "prod",
		}

		flags := parseEnvVars(tf, env)
		for k, v := range want {
			require.Equal(v, flags[k], "parseEnvVars() did not set the value")
		}
	})

	t.Run("flag override", func(t *testing.T) {
		// The flags should take precedence over the env.
		tf := map[string]string{
			"api_host": "127.0.0.1",
			"api_port": "9090",
			"debug":    "true",
		}
		env := map[string]string{
			"CUTTLE_API_HOST": "192.168.0.1",
			"CUTTLE_API_PORT": "3000",
			"CUTTLE_DB_ROOT":  "/tmp/db",
			"CUTTLE_DEBUG":    "false",
			"CUTTLE_ENV":      "prod",
		}

		flags := parseEnvVars(tf, env)
		for k, v := range want {
			require.Equal(v, flags[k], "parseEnvVars() did not set the value")
		}
	})
}

func writeFile(file string, config *Config) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o700); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0o600)
}

func removeFile(file string) error {
	return os.Remove(file)
}

func TestConfigGetConfigLocation(t *testing.T) {
	require := require.New(t)
	c := defaultConfig()

	t.Run("no file", func(t *testing.T) {
		file, err := getConfigLocation("")
		require.NoError(err, "getConfigLocation() returned an error: %s", err)
		require.Equal(FileNotFound, file, "getConfigLocation() did not return the default location")
	})

	t.Run("missing file", func(t *testing.T) {
		f := "/tmp/cuttle.yaml"
		removeFile(f)
		file, err := getConfigLocation(f)
		require.Error(err, "getConfigLocation() did not return an error")
		require.Equal("", file, "getConfigLocation() did not return the correct location")
	})

	t.Run("file set", func(t *testing.T) {
		f := "/tmp/cuttle.yaml"
		writeFile(f, c)
		file, err := getConfigLocation(f)
		require.NoError(err, "getConfigLocation() returned an error: %s", err)
		require.Equal(f, file, "getConfigLocation() did not return the correct location")
		removeFile(f)
	})

	t.Run("default files", func(t *testing.T) {
		hdir, err := os.UserHomeDir()
		require.NoError(err, "os.UserHomeDir() returned an error: %s", err)

		cwd, err := os.Getwd()
		require.NoError(err, "os.Getwd() returned an error: %s", err)

		files := []string{
			hdir + "/cuttle.yaml",
			hdir + "/.config/cuttle/config.yaml",
			cwd + "/cuttle.yaml",
			cwd + "/config.yaml",
		}

		for _, file := range files {
			writeFile(file, c)
			got, err := getConfigLocation("")
			require.NoError(err, "getConfigLocation() returned an error: %s", err)
			require.Equal(file, got, "getConfigLocation() did not return the correct location")
			removeFile(file)
		}
	})
}

func TestConfigParseConfigFile(t *testing.T) {
	require := require.New(t)
	c := defaultConfig()
	want := &Config{
		Env:     "prod",
		Debug:   true,
		APIHost: "127.0.0.1",
		APIPort: "9090",
		DBRoot:  "/tmp/db",
	}

	t.Run("no file", func(t *testing.T) {
		err := parseConfigFile(c, "")
		require.NoError(err, "parseConfigFile() returned an error: %s", err)
		require.Equal(defaultConfig(), c, "parseConfigFile() did not return the default config")
	})

	t.Run("missing file", func(t *testing.T) {
		file := "/tmp/cuttle.yaml"
		removeFile(file)
		err := parseConfigFile(c, file)
		require.Error(err, "parseConfigFile() did not return an error")
		require.Contains(err.Error(), "no such file or directory", "parseConfigFile() did not return the correct error")
	})

	t.Run("invalid config", func(t *testing.T) {
		file := "/tmp/cuttle.yaml"
		os.WriteFile(file, []byte("invalid"), 0o600)
		err := parseConfigFile(c, file)
		require.Error(err, "parseConfigFile() did not return an error")
		require.Equal(defaultConfig(), c, "parseConfigFile() did not return the default config")
		removeFile(file)
	})

	t.Run("valid config", func(t *testing.T) {
		file := "/tmp/cuttle.yaml"
		writeFile(file, want)
		err := parseConfigFile(c, file)
		require.NoError(err, "parseConfigFile() returned an error: %s", err)
		require.Equal(want, c, "parseConfigFile() did not return the correct config")
		removeFile(file)
	})
}
