package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/chadeldridge/cuttle/connections"
)

var (
	results bytes.Buffer
	logs    bytes.Buffer
)

func main() {
	server, err := connections.NewServer("test.home", 0, &results, &logs)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.SetIP("192.168.50.105"); err != nil {
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
		log.Println(err)
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
