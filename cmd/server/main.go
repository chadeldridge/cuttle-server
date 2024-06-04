package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/chadeldridge/cuttle/profiles"
)

var (
	results bytes.Buffer
	logs    bytes.Buffer

	host = "localhost"
	user = "bob"
	pass = "testUserP@ssw0rd"
)

func main() {
	profile := profiles.NewProfile("My Profile")
	defer connections.Pool.CloseAll()

	// Setup server
	server, err := connections.NewServer(host, 0, &results, &logs)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := connections.NewSSHConnector(user)
	if err != nil {
		log.Fatal(err)
	}

	server.SetConnector(&conn)
	conn.AddPasswordAuth(pass)

	err = server.TestConnection()
	if err != nil {
		log.Fatalf("TestConnection error: %s\n", err)
	}

	// Create group with server and add to profile
	group := profiles.NewGroup("Test Group", server)
	profile.AddGroups(group)

	/*
		if err := server.SetIP("192.168.50.105"); err != nil {
			log.Fatal(err)
		}
	*/

	tile := profiles.NewTile("echo Test", "echo 'my test echo'", "my test echo")
	err = profile.AddTiles(tile)
	if err != nil {
		log.Fatal(err)
	}

	err = profile.Execute(tile.Name(), group.Name)
	if err != nil {
		log.Printf("Execute error: %s\n", err)
	}

	fmt.Printf("--- Results ---\n%s\n", results.String())
	fmt.Printf("\n--- Logs ---\n%s\n", logs.String())
}
