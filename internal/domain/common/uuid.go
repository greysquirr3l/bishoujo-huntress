package common

import (
	"github.com/google/uuid"
)

// UUID is a wrapper around uuid.UUID to provide domain-specific functionality
type UUID uuid.UUID

// NewUUID creates a new random UUID
func NewUUID() UUID {
	return UUID(uuid.New())
}

// ParseUUID parses a UUID string
func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID(id), err
}

// String returns the string representation of the UUID
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// IsZero returns true if the UUID is the zero value
func (u UUID) IsZero() bool {
	return uuid.UUID(u) == uuid.Nil
}

// Equal returns true if the UUIDs are equal
func (u UUID) Equal(other UUID) bool {
	return uuid.UUID(u) == uuid.UUID(other)
}
