package main

import (
	"github.com/chadeldridge/cuttle-server/core"
)

var (
	version = "Cuttle v0.1.0"
	help    = `
Usage:
	cuttle [options] [args]
Options:
	--help				Print this help message.
	--version			Print the version.
	-c, --config-file <path>	Path to the configuration file.
	-C, --cert-file <path>		Path to the TLS certificate.
	-d, --db-root <path>		Path to the database.
	-e, --env <env>			Environment to run the server in.
	-h, --host <host>		Host to bind the API to.
	-k, --key-file <path>		Path to the TLS key.
	-p, --port <port>		Port to bind the API to.
	-v, --verbose			Enable verbose logging.`
)

// Skip the app name and return all flags and arguments.
func parseFlags(logger *core.Logger, args []string) (map[string]string, []string) {
	flags := map[string]string{}
	a := []string{}
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--help":
			logger.Println(help)
			return nil, nil
		case "--version":
			logger.Println(version)
			return nil, nil
		case "-c", "--config-file":
			flags["config_file"] = args[i+1]
			i++
		case "-C", "--cert-file":
			flags["tls_cert_file"] = args[i+1]
		case "-d", "--db-root":
			flags["db_root"] = args[i+1]
			i++
		case "-e", "--env":
			flags["env"] = args[i+1]
			i++
		case "-h", "--host":
			flags["api_host"] = args[i+1]
			i++
		case "-k", "--key-file":
			flags["tls_key_file"] = args[i+1]
		case "-p", "--port":
			flags["api_port"] = args[i+1]
			i++
		case "-v", "--verbose":
			flags["debug"] = "true"
		default:
			a = append(a, args[i])
		}
	}

	return flags, a
}
