package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConfigNewConfig(t *testing.T) {
	require := require.New(t)
	tester = MockTester
	reader = MockReader
	want := map[string]string{
		"api_host":      "127.0.0.1",
		"api_port":      "9090",
		"db_root":       "/tmp/db",
		"debug":         "true",
		"env":           "prod",
		"tls_cert_file": "/tmp/cuttle.pem",
		"tls_key_file":  "/tmp/cuttle.key",
	}

	t.Run("no flags", func(t *testing.T) {
		flags := map[string]string{}
		args := []string{}
		env := map[string]string{}
		MockClearFS()

		c, err := NewConfig(flags, args, env)
		require.Error(err, "NewConfig() returned an error: %s", err)
		require.Equal("tls cert: file not found", err.Error(), "NewConfig() did not return the correct error")
		require.Equal(DefaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("flags only", func(t *testing.T) {
		flags := map[string]string{
			"api_host":      "127.0.0.1",
			"api_port":      "9090",
			"db_root":       "/tmp/db",
			"debug":         "true",
			"env":           "prod",
			"tls_cert_file": "/tmp/cuttle_cert.cert",
			"tls_key_file":  "/tmp/cuttle_key.pem",
		}
		args := []string{}
		env := map[string]string{}
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
	})

	t.Run("env only", func(t *testing.T) {
		flags := map[string]string{}
		env := map[string]string{
			"CUTTLE_API_HOST":      "127.0.0.1",
			"CUTTLE_API_PORT":      "9090",
			"CUTTLE_DB_ROOT":       "/tmp/db",
			"CUTTLE_DEBUG":         "true",
			"CUTTLE_ENV":           "prod",
			"CUTTLE_TLS_CERT_FILE": "/tmp/cuttle_cert.cert",
			"CUTTLE_TLS_KEY_FILE":  "/tmp/cuttle_key.pem",
		}
		args := []string{}
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
	})

	t.Run("from file", func(t *testing.T) {
		flags := map[string]string{"config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{}
		args := []string{}
		config := &Config{
			Env:         want["env"],
			Debug:       true,
			TLSCertFile: "/tmp/cuttle_cert.cert",
			TLSKeyFile:  "/tmp/cuttle_key.pem",
			APIHost:     want["api_host"],
			APIPort:     want["api_port"],
			DBRoot:      want["db_root"],
		}

		data, err := yaml.Marshal(config)
		require.NoError(err, "yaml.Marshal() returned an error: %s", err)
		MockWriteFile("/tmp/cuttle.yaml", data, true, nil)
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
	})

	t.Run("missing file", func(t *testing.T) {
		flags := map[string]string{"config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{}
		args := []string{}
		MockClearFS()

		c, err := NewConfig(flags, args, env)
		require.Error(err, "NewConfig() did not return an error")
		require.Equal(
			"MockTester /tmp/cuttle.yaml: no such file or directory",
			err.Error(),
			"NewConfig() did not return the correct error",
		)
		require.Equal(DefaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("invalid flag", func(t *testing.T) {
		flags := map[string]string{"invalid": "value"}
		env := map[string]string{}
		args := []string{}

		c, err := NewConfig(flags, args, env)
		require.Error(err, "NewConfig() did not return an error")
		require.Equal(ErrUnknownOpt, err, "NewConfig() did not return the correct error")
		require.Equal(DefaultConfig(), c, "NewConfig() did not return the default config")
	})

	t.Run("override", func(t *testing.T) {
		flags := map[string]string{"api_host": want["api_host"], "config_file": "/tmp/cuttle.yaml"}
		env := map[string]string{"CUTTLE_API_HOST": "63.57.123.41", "CUTTLE_API_PORT": want["api_port"]}
		args := []string{}
		config := &Config{
			Env:         want["env"],
			Debug:       true,
			TLSCertFile: "/tmp/cuttle_cert.cert",
			TLSKeyFile:  "/tmp/cuttle_key.pem",
			APIHost:     "192.168.0.1",
			APIPort:     "3000",
			DBRoot:      want["db_root"],
		}

		data, err := yaml.Marshal(config)
		require.NoError(err, "yaml.Marshal() returned an error: %s", err)
		MockWriteFile("/tmp/cuttle.yaml", data, true, nil)
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		c, err := NewConfig(flags, args, env)
		require.NoError(err, "NewConfig() returned an error: %s", err)
		testConfig(t, want, c)
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
	c := DefaultConfig()

	t.Run("set one", func(t *testing.T) {
		err := c.setConfigValue("api_host", "127.0.0.1")
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
			err := c.setConfigValue(k, v)
			require.NoError(err, "setConfigValue() returned an error")
		}

		require.Equal(tc, *c, "setConfigValue() did not set the value")
	})

	t.Run("invalid key", func(t *testing.T) {
		err := c.setConfigValue("invalid", "value")
		require.Error(err, "setConfigValue() did not return an error")
		require.Equal(ErrUnknownOpt, err, "setConfigValue() did not return the correct error")
	})

	t.Run("invalid env", func(t *testing.T) {
		err := c.setConfigValue("env", "invalid")
		require.Error(err, "setConfigValue() did not return an error")
		require.Equal(ErrInvalidEnv, err, "setConfigValue() did not return the correct error")
	})

	t.Run("debug false", func(t *testing.T) {
		err := c.setConfigValue("debug", "false")
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

func TestConfigSetTLSFiles(t *testing.T) {
	require := require.New(t)
	c := DefaultConfig()
	tester = MockTester

	t.Run("not set", func(t *testing.T) {
		err := c.setTLSFiles()
		require.Error(err, "setTLSFiles() did not return an error")
		require.Equal("tls cert: file not found", err.Error(), "setTLSFiles() did not set the cert file")
		require.Equal("", c.TLSKeyFile, "setTLSFiles() did not set the key file")
	})

	t.Run("cert set", func(t *testing.T) {
		c.TLSCertFile = "/tmp/cuttle_cert.cert"
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)

		err := c.setTLSFiles()
		require.Error(err, "setTLSFiles() did not return an error")
		require.Equal("tls key: file not found", err.Error(), "setTLSFiles() did not set the key file")
		require.Equal("", c.TLSKeyFile, "setTLSFiles() did not set the key file")
	})

	t.Run("both set", func(t *testing.T) {
		c.TLSCertFile = "/tmp/cuttle_cert.cert"
		c.TLSKeyFile = "/tmp/cuttle_key.pem"
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		err := c.setTLSFiles()
		require.NoError(err, "setTLSFiles() returned an error: %s", err)
		require.Equal("/tmp/cuttle_cert.cert", c.TLSCertFile, "setTLSFiles() did not set the cert file")
		require.Equal("/tmp/cuttle_key.pem", c.TLSKeyFile, "setTLSFiles() did not set the key file")
	})
}

func TestConfigSetTLSCertFile(t *testing.T) {
	require := require.New(t)
	c := DefaultConfig()
	tester = MockTester

	t.Run("not set", func(t *testing.T) {
		err := c.setTLSCertFile()
		require.Error(err, "setTLSFiles() did not return an error")
		require.Equal("tls cert: file not found", err.Error(), "setTLSFiles() did not set the cert file")
		require.Equal("", c.TLSKeyFile, "setTLSFiles() did not set the key file")
	})

	t.Run("specified", func(t *testing.T) {
		c.TLSCertFile = "/tmp/cuttle_cert.cert"
		MockWriteFile("/tmp/cuttle_cert.cert", []byte(""), true, nil)

		err := c.setTLSCertFile()
		require.NoError(err, "setTLSFiles() returned an error: %s", err)
		require.Equal("/tmp/cuttle_cert.cert", c.TLSCertFile, "setTLSFiles() set the cert file")
	})

	t.Run("found", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(err, "os.UserHomeDir() returned an error: %s", err)

		c.TLSCertFile = ""
		MockWriteFile(dir+"/cuttle_cert.cert", []byte(""), true, nil)

		err = c.setTLSCertFile()
		require.NoError(err, "setTLSFiles() returned an error: %s", err)
		require.Equal(dir+"/cuttle_cert.cert", c.TLSCertFile, "setTLSFiles() set the cert file")
	})
}

func TestConfigSetTLSKeyFile(t *testing.T) {
	require := require.New(t)
	c := DefaultConfig()
	tester = MockTester

	t.Run("not set", func(t *testing.T) {
		err := c.setTLSKeyFile()
		require.Error(err, "setTLSFiles() did not return an error")
		require.Equal("tls key: file not found", err.Error(), "setTLSFiles() did not set the key file")
		require.Equal("", c.TLSKeyFile, "setTLSFiles() did not set the key file")
	})

	t.Run("specified", func(t *testing.T) {
		c.TLSKeyFile = "/tmp/cuttle_key.pem"
		MockWriteFile("/tmp/cuttle_key.pem", []byte(""), true, nil)

		err := c.setTLSKeyFile()
		require.NoError(err, "setTLSFiles() returned an error: %s", err)
		require.Equal("/tmp/cuttle_key.pem", c.TLSKeyFile, "setTLSFiles() set the key file")
	})

	t.Run("found", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(err, "os.UserHomeDir() returned an error: %s", err)

		MockWriteFile(dir+"/cuttle_key.pem", []byte(""), true, nil)
		c.TLSKeyFile = ""

		err = c.setTLSKeyFile()
		require.NoError(err, "setTLSFiles() returned an error: %s", err)
		require.Equal(dir+"/cuttle_key.pem", c.TLSKeyFile, "setTLSFiles() set the key file")
	})
}

func TestConfigParseConfigFile(t *testing.T) {
	require := require.New(t)
	tester = MockTester
	reader = MockReader
	c := DefaultConfig()
	want := &Config{
		Env:     "prod",
		Debug:   true,
		APIHost: "127.0.0.1",
		APIPort: "9090",
		DBRoot:  "/tmp/db",
	}

	t.Run("no files", func(t *testing.T) {
		err := c.parseConfigFile("")
		require.Error(err, "parseConfigFile() did not return an error")
		require.Equal("file not found", err.Error(), "parseConfigFile() did not return the correct error")
		require.Equal(DefaultConfig(), c, "parseConfigFile() did not return the default config")
	})

	t.Run("missing file", func(t *testing.T) {
		file := "/tmp/cuttle.yaml"
		MockClearFS()

		err := c.parseConfigFile(file)
		require.Error(err, "parseConfigFile() did not return an error")
		require.Equal(
			"MockTester /tmp/cuttle.yaml: no such file or directory",
			err.Error(),
			"parseConfigFile() did not return the correct error",
		)
	})

	t.Run("invalid config", func(t *testing.T) {
		MockWriteFile("/tmp/cuttle.yaml", []byte("invalid"), true, nil)

		err := c.parseConfigFile("/tmp/cuttle.yaml")
		require.Error(err, "parseConfigFile() did not return an error")
		require.Equal(DefaultConfig(), c, "parseConfigFile() did not return the default config")
	})

	t.Run("valid config", func(t *testing.T) {
		data, err := yaml.Marshal(want)
		require.NoError(err, "yaml.Marshal() returned an error: %s", err)
		MockWriteFile("/tmp/cuttle.yaml", data, true, nil)

		err = c.parseConfigFile("/tmp/cuttle.yaml")
		require.NoError(err, "parseConfigFile() returned an error: %s", err)
		require.Equal(want, c, "parseConfigFile() did not return the correct config")
	})

	t.Run("found config", func(t *testing.T) {
		dir, err := os.UserHomeDir()
		require.NoError(err, "os.UserHomeDir() returned an error: %s", err)

		data, err := yaml.Marshal(want)
		require.NoError(err, "yaml.Marshal() returned an error: %s", err)
		MockWriteFile(dir+"/cuttle.yaml", data, true, nil)

		err = c.parseConfigFile("")
		require.NoError(err, "parseConfigFile() returned an error: %s", err)
		require.Equal(want, c, "parseConfigFile() did not return the correct config")
	})
}
