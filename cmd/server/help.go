package main

import "fmt"

func printVersion() {
	fmt.Println("Cuttle v0.1.0")
}

func printHelp() {
	fmt.Print(`
Usage:
	cuttle [options] [args]
Options:
	-c, --config-file <path>	Path to the configuration file.
	-d, --db-root <path>		Path to the database.
	-h, --help			Print this help message.
	-H, --host <host>		Host to bind the API to.
	-p, --port <port>		Port to bind the API to.
	-v, --version			Print the version.
`)
}
