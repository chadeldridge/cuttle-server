package core

import (
	"fmt"
	"log"
	"regexp"
)

func IsHex(s string) bool {
	if s == "" {
		return false
	}

	fail, err := regexp.MatchString("[^0-9a-fA-F]", s)
	if err != nil {
		log.Printf("IsHex: %s", err)
		return false
	}

	return !fail
}

func ValidatePasswordHash(hash string) error {
	if hash == "" {
		return fmt.Errorf("core.ValidatePasswordHash: hash was empty")
	}

	if len(hash) < 32 {
		return fmt.Errorf("core.ValidatePasswordHash: incorrect hash length: %d", len(hash))
	}

	if !IsHex(hash) {
		return fmt.Errorf("core.ValidatePasswordHash: hash is not a hex string")
	}

	return nil
}
