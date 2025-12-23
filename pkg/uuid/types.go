package uuid

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

type UUIDHandler interface {
	GenerateV4() (UUID, error)
	GenerateV7() (UUID, error)
	Parse(s string) (UUID, error)
}
