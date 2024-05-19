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
	server, err := profiles.NewServer("test.home", "sshpwd")
	if err != nil {
		log.Fatal(err)
	}

	/*
		if err := server.SetHostname("test.home"); err != nil {
			log.Fatal(err)
		}
	*/

	if err := server.SetIP("192.168.50.105"); err != nil {
		log.Fatal(err)
	}

	conn, err := connections.NewSSH(server, &results, &logs)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetUser("debian")
	conn.SetPassword(getPassword())

	testConnection(conn)
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

func testConnection(conn connections.Handler) {
	err := conn.TestConnection()
	if err != nil {
		log.Println(err)
	}
}
