package ddd

import (
	"github.com/google/uuid"
)

type IDomainEvent interface {
	ID() uuid.UUID
	Name() string
}
