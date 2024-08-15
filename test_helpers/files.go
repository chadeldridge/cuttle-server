package test_helpers

import (
	"log"
	"os"
)

func DeleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
		log.Fatalf("deleteDB: %s", err)
	}
}
