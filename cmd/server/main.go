package main

// Credit: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
// HTTP server initially modeled after Mat Ryer's article. Perfect timing hearing about it on
// the Go Time podcast right when I was thinking about dropping gin-gonic and going with a more
// bare bones solution for this project.

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/chadeldridge/cuttle/api"
	"github.com/chadeldridge/cuttle/core"
	"github.com/chadeldridge/cuttle/db"
	"github.com/chadeldridge/cuttle/router"
	"github.com/chadeldridge/cuttle/web"
)

/*
var (
	results bytes.Buffer
	logs    bytes.Buffer

	remoateHost = "test.home"
	remoteUser  = "bob"
	pass        = "testUserP@ssw0rd"
	// encPass     = []byte("Myv3ryGo0dandsupersecureP@sswordTM$")
)
*/

// run allows us to setup and implement in testing and production.
func run(ctx context.Context, out io.Writer, args []string, env map[string]string) error {
	// Capture the interrupt signal to gracefully shutdown the server.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Setup logger.
	logger := core.NewLogger(out, "cuttle: ", log.LstdFlags, false)

	// Get flags.
	flags, args := parseFlags(logger, args)
	if flags == nil && args == nil {
		return nil
	}

	// Setup the configuration.
	config, err := core.NewConfig(flags, args, env)
	if err != nil {
		return err
	}

	// Update logger with config value.
	logger.DebugMode = config.Debug

	// Print config if in debug mode.
	logger.Debugf("Config: %+v\n", config)

	// Setup the database.
	db.SetDBRoot(config.DBRoot)
	mainDB, err := db.NewSqliteDB("cuttle.db")
	if err != nil {
		return err
	}

	err = mainDB.Open()
	if err != nil {
		return err
	}
	defer mainDB.Close()

	users, err := db.NewUsers(mainDB)
	if err != nil {
		return err
	}

	// Setup the HTTP server.
	srv := router.NewHTTPServer(logger, config)
	srv.Users = &users
	// Add routes and do anything else we need to do before starting the server.

	// Add web routes.
	err = web.AddRoutes(&srv)
	if err != nil {
		return err
	}

	// Add API routes.
	err = api.AddRoutes(&srv)
	if err != nil {
		return err
	}

	// Start API Server.
	return srv.Start(ctx, config.ShutdownTimeout)
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args, getEnv()); err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}

func getEnv() map[string]string {
	env := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}

	return env
}

/*
	defer connections.Pool.CloseAll()

	// Create a base profile.
	profile, err := profiles.NewProfile("My Profile")
	if err != nil {
		log.Fatal(err)
	}

	// Setup our test server.
	server, err := connections.NewServer(remoateHost, 0, &results, &logs)
	if err != nil {
		log.Fatal(err)
	}

	// Create a connector and add it to the server.
	conn, err := connections.NewSSHConnector("my connector", remoteUser)
	if err != nil {
		log.Fatal(err)
	}

	conn.AddPasswordAuth(pass)
	server.SetConnector(&conn)

	// Test the connections to the server.
	err = server.TestConnection()
	if err != nil {
		log.Fatalf("TestConnection error: %s\n", err)
	}

	// Create group with the server and add to the profile.
	group := profiles.NewGroup("Test Group", server)
	profile.AddGroups(group)
*/

/*
	if err := server.SetIP("192.168.50.105"); err != nil {
		log.Fatal(err)
	}
*/

/*
	// Create some test to add to a tile.
	ping := tests.NewPingTest("Ping Test", false, 1, tests.Quiet())
	tcpHalfOpen := tests.NewTCPPortHalfOpen("TCP Half Open Test", false, 22)
	tcpOpen := tests.NewTCPPortOpen("TCP Open Test", false, 22)
	sshTest := tests.NewSSHTest("SSH Test", false, "echo 'Hello, World!'", "Hello, World!")

	// Create a tile with the tests and add it to the profile.
	tile := profiles.NewTile("Full Test", ping, tcpHalfOpen, tcpOpen, sshTest)
	err = profile.AddTiles(tile)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the tile with the group.
	err = profile.Execute(tile.Name, group.Name)
	if err != nil {
		log.Printf("Execute error: %s\n", err)
	}

	fmt.Printf("--- Results ---\n%s\n", results.String())
	fmt.Printf("\n--- Logs ---\n%s\n", logs.String())
*/
