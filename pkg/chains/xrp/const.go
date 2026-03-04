package xrp

import "strconv"

// ToFloat64 converts a string amount to float64.
func ToFloat64(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0
	}
	return f
}

// XRPToDrops converts an XRP amount to drops. 1 XRP = 1,000,000 drops, so 1 drop = 0.000001 XRP.
// See: https://xrpl.org/rippleapi-reference.html#xrptodrops
func XRPToDrops(val float64) float64 {
	return val * 0.000001
}

type HashVersion byte

const (
	AccountZero = "rrrrrrrrrrrrrrrrrrrrrhoLvTp"
	AccountOne  = "rrrrrrrrrrrrrrrrrrrrBZbvji"
	NaN         = "rrrrrrrrrrrrrrrrrrrn5RM1rHd"
	ROOT        = "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
)

// MinimumReserve is the minimum XRP reserve required for an account.
// See: https://xrpl.org/reserves.html
const MinimumReserve float64 = 20.0

// MaxLedgerVersionOffset is the default offset added to the current ledger index for transaction expiry.
// A value of 15 means approximately 60 seconds (each ledger takes ~4 seconds).
const MaxLedgerVersionOffset uint64 = 15

const (
	ALPHABET = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

	RippleAccountID      HashVersion = 0
	RippleNodePublic     HashVersion = 28
	RippleNodePrivate    HashVersion = 32
	RippleFamilySeed     HashVersion = 33
	RippleAccountPrivate HashVersion = 34
	RippleAccountPublic  HashVersion = 35
)

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
