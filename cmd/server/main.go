package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/chadeldridge/cuttle/profiles"
	"github.com/chadeldridge/cuttle/tests"
)

var (
	results bytes.Buffer
	logs    bytes.Buffer

	remoateHost = "test.home"
	remoteUser  = "bob"
	pass        = "testUserP@ssw0rd"
	// encPass     = []byte("Myv3ryGo0dandsupersecureP@sswordTM$")
)

func main() {
	/*
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
	*/
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

	/*
		if err := server.SetIP("192.168.50.105"); err != nil {
			log.Fatal(err)
		}
	*/

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
}
