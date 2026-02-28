package rpc

// RequestCommand is the minimal request body for commands with no parameters.
type RequestCommand struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

// KeyType represents the cryptographic key algorithm used by an XRP Ledger account.
type KeyType string

const (
	// KeyTypeSECP256K1 is the secp256k1 elliptic curve algorithm (default).
	KeyTypeSECP256K1 KeyType = "secp256k1"
	// KeyTypeED25519 is the Ed25519 algorithm.
	KeyTypeED25519 KeyType = "ed25519"
)

// String returns the string representation of the KeyType.
func (k KeyType) String() string {
	return string(k)
}
