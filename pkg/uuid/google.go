package uuid

import (
	"github.com/google/uuid"
)

//
// googleUUIDHandler implements UUIDHandler interface
//

type googleUUIDHandler struct{}

// Ensure googleUUIDHandler implements UUIDHandler at compile time.
var _ UUIDHandler = (*googleUUIDHandler)(nil)

func NewGoogleUUIDHandler() *googleUUIDHandler {
	return &googleUUIDHandler{}
}

func (*googleUUIDHandler) GenerateV4() (UUID, error) {
	return uuid.NewRandom()
}

func (*googleUUIDHandler) GenerateV7() (UUID, error) {
	return uuid.NewV7()
}

func (*googleUUIDHandler) Parse(s string) (UUID, error) {
	return uuid.Parse(s)
}
