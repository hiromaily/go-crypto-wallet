package xrp

import (
	"fmt"
	"math/big"

	addresscodec "github.com/Peersyst/xrpl-go/address-codec"
)

// First byte is the network
// Second byte is the version
// Remaining bytes are the payload
type hash []byte

var hashTypes = [...]struct {
	Description       string
	Prefix            byte
	Payload           int
	MaximumCharacters int
}{
	RippleAccountID:      {"Short name for sending funds to an account.", 'r', 20, 35},
	RippleNodePublic:     {"Validation public key for node.", 'n', 33, 53},
	RippleNodePrivate:    {"Validation private key for node.", 'p', 32, 52},
	RippleFamilySeed:     {"Family seed.", 's', 16, 29},
	RippleAccountPrivate: {"Account private key.", 'p', 32, 52},
	RippleAccountPublic:  {"Account public key.", 'a', 33, 53},
}

func NewRippleHash(s string) (Hash, error) {
	// Special case which will deal short addresses
	switch s {
	case "0":
		return newHashFromString(AccountZero)
	case "1":
		return newHashFromString(AccountOne)
	default:
		return newHashFromString(s)
	}
}

// NewRippleHashCheck checks hash matches expected version
func NewRippleHashCheck(s string, version HashVersion) (Hash, error) {
	hash, err := NewRippleHash(s)
	if err != nil {
		return nil, err
	}
	if hash.Version() != version {
		want := hashTypes[version].Description
		got := hashTypes[hash.Version()].Description
		return nil, fmt.Errorf("bad version for: %s expected: %s got: %s ", s, want, got)
	}
	return hash, nil
}

func NewAccountID(b []byte) (Hash, error) {
	return newHash(b, RippleAccountID)
}

func NewAccountPublicKey(b []byte) (Hash, error) {
	return newHash(b, RippleAccountPublic)
}

func NewAccountPrivateKey(b []byte) (Hash, error) {
	return newHash(b, RippleAccountPrivate)
}

func NewNodePublicKey(b []byte) (Hash, error) {
	return newHash(b, RippleNodePublic)
}

func NewNodePrivateKey(b []byte) (Hash, error) {
	return newHash(b, RippleNodePrivate)
}

func NewFamilySeed(b []byte) (Hash, error) {
	return newHash(b, RippleFamilySeed)
}

// EncodeEd25519Seed encodes 16 seed bytes as an ed25519 family seed string
// using the XRPL ed25519 prefix [0x01, 0xe1, 0x4b].
//
// This format is recognized by Peersyst/xrpl-go's wallet.FromSeed() and
// addresscodec.DecodeSeed() as an ed25519 seed, enabling correct keypair derivation.
// The resulting string starts with "sEd" in XRPL's Base58 alphabet.
func EncodeEd25519Seed(b []byte) (string, error) {
	if len(b) != 16 {
		return "", fmt.Errorf("ed25519 seed must be exactly 16 bytes, got %d", len(b))
	}
	// Prefix [0x01, 0xe1, 0x4b] followed by 16 seed bytes, then Base58Check encoded
	return addresscodec.Base58CheckEncode(b, 0x01, 0xe1, 0x4b), nil
}

func AccountID(key Key, sequence *uint32) (Hash, error) {
	return NewAccountID(key.Id(sequence))
}

func AccountPublicKey(key Key, sequence *uint32) (Hash, error) {
	return NewAccountPublicKey(key.Public(sequence))
}

func AccountPrivateKey(key Key, sequence *uint32) (Hash, error) {
	return NewAccountPrivateKey(key.Private(sequence))
}

func NodePublicKey(key Key) (Hash, error) {
	return NewNodePublicKey(key.Public(nil))
}

func NodePrivateKey(key Key) (Hash, error) {
	return NewNodePrivateKey(key.Private(nil))
}

func GenerateFamilySeed(password string) (Hash, error) {
	return NewFamilySeed(Sha512Quarter([]byte(password)))
}

func newHash(b []byte, version HashVersion) (Hash, error) {
	n := hashTypes[version].Payload
	if len(b) > n {
		return nil, fmt.Errorf("Hash is wrong size, expected: %d got: %d", n, len(b))
	}
	return append(hash{byte(version)}, b...), nil
}

func newHashFromString(s string) (Hash, error) {
	decoded, err := addresscodec.Base58CheckDecode(s)
	if err != nil {
		return nil, err
	}
	return hash(decoded), nil
}

func (h hash) String() string {
	return addresscodec.Base58CheckEncode(h.Payload(), byte(h.Version()))
}

func (h hash) Version() HashVersion {
	return HashVersion(h[0])
}

func (h hash) Payload() []byte {
	return h[1:]
}

// Return a slice of the payload with leading zeroes omitted
func (h hash) PayloadTrimmed() []byte {
	payload := h.Payload()
	for i := range payload {
		if payload[i] != 0 {
			return payload[i:]
		}
	}
	return payload[len(payload)-1:]
}

func (h hash) Value() *big.Int {
	return big.NewInt(0).SetBytes(h.Payload())
}

func (h hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h hash) Clone() Hash {
	c := make(hash, len(h))
	copy(c, h)
	return c
}
