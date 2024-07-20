package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/chadeldridge/cuttle/api"
	"github.com/chadeldridge/cuttle/core"
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
func run(ctx context.Context, out io.Writer, args []string, getenv func(string) string) error {
	// Capture the interrupt signal to gracefully shutdown the server.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Get flags.
	flags, args := parseFlags(args)

	// Setup the configuration.
	config, err := core.NewConfig(flags, args, getenv)
	if err != nil {
		return err
	}

	// Setup logger.
	logger := core.NewLogger(out, "cuttle: ", log.LstdFlags, config.Debug)

	// Remove later. Debug only.
	logger.Debugf("Config: %+v\n", config)

	// Setup the HTTP server.
	srv := api.NewHTTPServer(logger, config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.APIHost, config.APIPort),
		Handler: srv,
	}

	// Start the server.
	go func() {
		logger.Printf("server listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Printf("%v\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Duration(config.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Printf("http server shutdown error: %v\n", err)
		}
	}()
	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
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
