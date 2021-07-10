package uuid

import (
	googleUuid "github.com/google/uuid"
)

// Generate creates and returns a new UUID as a string.
func Generate() string {
	return googleUuid.NewString()
}
