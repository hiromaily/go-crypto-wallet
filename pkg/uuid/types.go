package uuid

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

type UUIDHandler interface {
	GenerateV4() (uuid.UUID, error)
	GenerateV7() (uuid.UUID, error)
	Parse(s string) (UUID, error)
}
