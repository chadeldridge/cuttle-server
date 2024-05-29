package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/chadeldridge/cuttle/connections"
	"github.com/chadeldridge/cuttle/profiles"
)

var (
	results bytes.Buffer
	logs    bytes.Buffer
)

func main() {
	profile := profiles.NewProfile("My Profile")
	defer connections.Pool.CloseAll()

	// Setup server
	server, err := connections.NewServer("test.home", 0, &results, &logs)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := connections.NewSSH(&server, "debian")
	if err != nil {
		log.Fatal(err)
	}

	server.SetHandler(&conn)
	conn.SetPassword(getPassword())

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

func getPassword() string {
	if p, ok := os.LookupEnv("PASSWORD"); ok {
		return p
	}

	log.Fatal("failed to get env PASSWORD")
	return ""
}
