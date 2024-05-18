package main

import (
	"fmt"
	"log"

	"github.com/chadeldridge/cuttle"
)

func main() {
	server, err := cuttle.NewServer("test.home", "sshpwd")
	if err != nil {
		log.Fatal(err)
	}

	if err := server.SetHostname("test.home"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("set hostname to test.home")

	if err := server.SetHostname("192.168.50.105"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("set hostname to 192.168.50.105")
	for _, b := range server.IP() {
		fmt.Printf(" %v,", b)
	}

	if err := server.SetHostname("89ey*(#@F*)89023r"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("set hostname to 89ey*(#@F*)89023r")
}

/*
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
*/
