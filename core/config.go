package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/chadeldridge/cuttle/db"
	"gopkg.in/yaml.v3"
)

//
// cli args > environment vars > config file > defaults
//

const (
	DefaultAPIHost         = "localhost"
	DefaultAPIPort         = "8080"
	DefaultDBRoot          = db.DefaultDBFolder
	DefaultShutdownTimeout = 10
)

var (
	FileNotFound  = "file not found"
	ErrInvalidEnv = fmt.Errorf("invalid environment")
	ErrUnknownOpt = fmt.Errorf("unknown option")

	supportedEnv  = []string{"prod", "dev"}
	supportedVars = []string{
		"CUTTLE_API_HOST",
		"CUTTLE_API_PORT",
		"CUTTLE_CONFIG_FILE",
		"CUTTLE_DB_ROOT",
		"CUTTLE_DEBUG",
		"CUTLE_ENV",
	}
)

type Config struct {
	Env             string `yaml:"env"`
	Debug           bool   `yaml:"debug"`
	APIHost         string `yaml:"api_host"`
	APIPort         string `yaml:"api_port"`
	DBRoot          string `yaml:"db_root"`          // DBRoot is the root path for the database.
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // in seconds
}

func NewConfig(flags map[string]string, args []string, getenv func(string) string) (*Config, error) {
	// Create a default config.
	c := &Config{
		APIHost:         DefaultAPIHost,
		APIPort:         DefaultAPIPort,
		DBRoot:          DefaultDBRoot,
		ShutdownTimeout: DefaultShutdownTimeout,
	}

	// Parse the environment variables into a normalized format.
	flags = parseEnvVars(flags, getenv)

	// If there's a config file, parse it and set the values in the config.
	file := ""
	if v, ok := flags["config_file"]; ok {
		file = v
	}

	if err := parseConfigFile(c, file); err != nil {
		return c, err
	}

	// Try to find each supported variable in the flags and env. If found, set the value in the config.
	for k, v := range flags {
		// Skip config_file and empty values.
		if k == "config_file" || v == "" {
			continue
		}

		err := setConfigValue(c, k, v)
		if err != nil {
			return c, err
		}
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

func setConfigValue(c *Config, k, v string) error {
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
	default:
		return ErrUnknownOpt
	}

	return nil
}

// Parse all supported environment variables into a map.
func parseEnvVars(flags map[string]string, getenv func(string) string) map[string]string {
	for _, k := range supportedVars {
		newKey := strings.ToLower(strings.TrimPrefix(k, "CUTTLE_"))
		// Don't overwrite existing flags that aren't empty.
		if v, ok := flags[newKey]; ok && v != "" {
			continue
		}

		// Get the value from the environment.
		v := getenv(k)
		if v == "" {
			continue
		}

		// Add the value to the flags map.
		flags[newKey] = v
	}

	return flags
}

// getConfigLocation returns the location of the config file. If a location is set and the file
// does not exist, panic. If no config file is found and we did not panic, return fileNotFound.
func getConfigLocation(file string) (string, error) {
	fileName := "config.yaml"

	// Check for a flag or env provided file.
	if file != "" {
		if _, err := os.Stat(file); err != nil {
			return "", err
		}
		return file, nil
	}

	var locations []string
	// Add the user's home directory.
	if dir, err := os.UserHomeDir(); err == nil {
		locations = append(locations, dir+"/cuttle.yaml", dir+"/.config/cuttle/"+fileName)
	}

	// Add the current working directory. Assume the filename contains the app name so we do
	// not load in another app's config by mistake.
	if dir, err := os.Getwd(); err == nil {
		locations = append(locations, dir+"/"+fileName, dir+"/cuttle.yaml")
	}

	for _, l := range locations {
		if _, err := os.Stat(l); err == nil {
			return l, nil
		}
	}

	return FileNotFound, nil
}

// parseConfigFile reads the config file and unmarshals the data into the config struct.
func parseConfigFile(c *Config, file string) error {
	// If no config file was found, return and use cli, env, or default values.
	file, err := getConfigLocation(file)
	if err != nil {
		return err
	}

	if file == FileNotFound {
		return nil
	}

	// If a file was found but cannot be read, assume the config file is required and return err.
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Unmarshal the data into the config struct and rerturn any errors.
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
