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
		return fmt.Errorf("core.ValidatePasswordHash: hash - %w", ErrParamEmpty)
	}

	match := regexp.MustCompile(`^\$2[ayb]\$[0-9]{2}\$[0-9a-zA-Z./]{53}$`)
	if ok := match.MatchString(hash); !ok {
		return fmt.Errorf("core.ValidatePasswordHash: hash - %w", ErrInvalidFormat)
	}

	return nil
}
