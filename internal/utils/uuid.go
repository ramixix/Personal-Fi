package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

// -------------------------------------------------
// GenerateUUID generates a UUID v4 (random)
// -------------------------------------------------
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, uuid); err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	// The Version: By manipulating the 6th byte, the code clears the high bits and sets them to 0100 (binary for 4). This marks it as a "Randomly Generated" UUID.
	// The Variant: The 8th byte is modified to set the variant to RFC 4122 (the standard). This ensures compatibility across different systems like Java, Python, or databases.
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	// Format as string: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// -----------------------------------------------------------------------------------
// MustGenerateUUID generates a UUID and panics on error. It is used specially at initialization, where panicing and stoping program is OK.
// -----------------------------------------------------------------------------------
func MustGenerateUUID() string {
	uuid, err := GenerateUUID()
	if err != nil {
		panic(err)
	}
	return uuid
}

// -------------------------------------------------------------
// IsValidUUID checks if a string is a valid UUID format
// -------------------------------------------------------------
func IsValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// Check format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// Check if all other characters are hex digits
	for i, c := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !isHexDigit(c) {
			return false
		}
	}

	return true
}

func isHexDigit(character rune) bool {
	return (character >= '0' && character <= '9') || (character >= 'a' && character <= 'f') || (character >= 'A' && character <= 'F')
}
