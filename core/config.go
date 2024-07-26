package core

import (
	"fmt"
	"maps"
	"strings"

	"github.com/chadeldridge/cuttle/db"
)

//
// cli args > environment vars > config file > defaults
//

const (
	DefaultAPIHost         = "0.0.0.0"
	DefaultAPIPort         = "8080"
	DefaultDBRoot          = db.DefaultDBFolder
	DefaultShutdownTimeout = 5
)

var (
	FileNotFound  = "file not found"
	ErrInvalidEnv = fmt.Errorf("invalid environment")
	ErrUnknownOpt = fmt.Errorf("unknown option")

	supportedEnv = []string{"prod", "dev"}
	/*
		supportedVars = []string{
			"CUTTLE_API_HOST",
			"CUTTLE_API_PORT",
			"CUTTLE_CONFIG_FILE",
			"CUTTLE_DB_ROOT",
			"CUTTLE_DEBUG",
			"CUTLE_ENV",
		}
	*/
)

type Config struct {
	Env             string `yaml:"env,omitempty"`
	Debug           bool   `yaml:"debug,omitempty"`
	TLSCertFile     string `yaml:"tls_cert_file,omitempty"`
	TLSKeyFile      string `yaml:"tls_key_file,omitempty"`
	APIHost         string `yaml:"api_host,omitempty"`
	APIPort         string `yaml:"api_port,omitempty"`
	DBRoot          string `yaml:"db_root,omitempty"`                      // DBRoot is the root path for the database.
	ShutdownTimeout int    `default:"5" yaml:"shutdown_timeout,omitempty"` // in seconds
}

func NewConfig(
	flags map[string]string,
	args []string,
	env map[string]string,
) (*Config, error) {
	// Create a default config.
	c := &Config{
		Env:             "dev",
		Debug:           false,
		TLSCertFile:     "",
		TLSKeyFile:      "",
		APIHost:         DefaultAPIHost,
		APIPort:         DefaultAPIPort,
		DBRoot:          DefaultDBRoot,
		ShutdownTimeout: DefaultShutdownTimeout,
	}

	// Parse the environment variables into a normalized format.
	flags = parseEnvVars(flags, env)

	// If there's a config file, parse it and set the values in the config.
	var file string
	if v, ok := flags["config_file"]; ok {
		file = v
	}

	if err := c.parseConfigFile(file); err != nil && err.Error() != FileNotFound {
		return c, err
	}

	// Try to find each supported variable in the flags and env. If found, set the value in the config.
	for k, v := range flags {
		// Skip config_file and empty values.
		if k == "config_file" || v == "" {
			continue
		}

		err := c.setConfigValue(k, v)
		if err != nil {
			return c, err
		}
	}

	// If ShutdownTimeout is 0 it will cause a shutdown error and hang the server when using TLS.
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = DefaultShutdownTimeout
	}

	// Fail if we cannot find the TLS files.
	if err := c.setTLSFiles(); err != nil {
		return c, err
	}

	return c, nil
}

func validateEnv(env string) bool {
	for _, v := range supportedEnv {
		if env == v {
			return true
		}
	}

	return false
}

func (c *Config) setConfigValue(k, v string) error {
	switch k {
	case "api_host":
		c.APIHost = v
	case "api_port":
		c.APIPort = v
	case "db_root":
		c.DBRoot = v
	case "debug":
		if v == "true" {
			c.Debug = true
			break
		}
		c.Debug = false
	case "env":
		v = strings.ToLower(v)
		if !validateEnv(v) {
			return ErrInvalidEnv
		}
		c.Env = v
	case "tls_cert_file":
		c.TLSCertFile = v
	case "tls_key_file":
		c.TLSKeyFile = v
	default:
		return ErrUnknownOpt
	}

	return nil
}

// Parse all supported environment variables into a map.
func parseEnvVars(flags, env map[string]string) map[string]string {
	f := map[string]string{}
	// Copy the environment variables into the map with normalized keys.
	for k, v := range env {
		if !strings.HasPrefix(k, "CUTTLE_") {
			continue
		}

		f[strings.ToLower(strings.TrimPrefix(k, "CUTTLE_"))] = v
	}

	maps.Copy(f, flags)
	return f
}

/*
func getConfigLocation(file string) (string, error) {
	// Check for a flag or env provided file.
	if file != "" {
		if _, err := os.Stat(file); err != nil {
			return "", err
		}
		return file, nil
	}

	// If the config file is not in a dedicated cuttle folder, prefer a filename of cuttle.yaml.
	var locations []string

	// Add the user's home directory.
	if dir, err := os.UserHomeDir(); err == nil {
		locations = append(locations, dir+"/cuttle.yaml", dir+"/.config/cuttle/config.yaml")
	}

	// Add the current working directory. Assume the filename contains the app name so we do
	// not load in another app's config by mistake.
	if dir, err := os.Getwd(); err == nil {
		locations = append(locations, dir+"/cuttle.yaml", dir+"/config.yaml")
	}

	for _, l := range locations {
		if _, err := os.Stat(l); err == nil {
			return l, nil
		}
	}

	return FileNotFound, nil
}
*/

func (c *Config) setTLSFiles() error {
	if err := c.setTLSCertFile(); err != nil {
		return err
	}

	return c.setTLSKeyFile()
}

func (c *Config) setTLSCertFile() error {
	certFiles := []string{
		"cuttle.pem",
		"cuttle_cert.pem",
		"cuttle.cert",
		"cuttle_cert.cert",
		"certs/cuttle.pem",
		"certs/cuttle_cert.pem",
		"certs/cuttle.cert",
		"cetrs/cuttle_cert.crt",
	}

	// If the cert file location is set, return the results of tester.
	if c.TLSCertFile != "" {
		return tester(c.TLSCertFile)
	}

	// Try to find the cert file location and return if none can be found.
	cert, err := FindFiles(certFiles...)
	if err != nil {
		return fmt.Errorf("tls cert: %w", err)
	}
	c.TLSCertFile = cert

	return nil
}

func (c *Config) setTLSKeyFile() error {
	keyFiles := []string{
		"cuttle.key",
		"cuttle_key.pem",
		"certs/cuttle.key",
		"certs/cuttle_key.pem",
	}

	// If the key file location is set, return the results of tester.
	if c.TLSKeyFile != "" {
		return tester(c.TLSKeyFile)
	}

	// Try to find the key file location and return if none can be found.
	key, err := FindFiles(keyFiles...)
	if err != nil {
		return fmt.Errorf("tls key: %w", err)
	}
	c.TLSKeyFile = key

	return nil
}

// parseConfigFile reads the config file and unmarshals the data into the config struct.
func (c *Config) parseConfigFile(file string) error {
	if file != "" {
		err := tester(file)
		if err != nil {
			return err
		}

		// If the file was found and readable, parse the data and return.
		return ParseYAML(file, c)
	}

	// If no config file was set, return and use cli, env, or default values.
	file, err := FindFiles("cuttle.yaml", "config.yaml")
	if err != nil {
		return err
	}

	return ParseYAML(file, c)
}
