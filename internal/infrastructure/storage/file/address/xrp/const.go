package xrp

type HashVersion byte

const (
	AccountZero = "rrrrrrrrrrrrrrrrrrrrrhoLvTp"
	AccountOne  = "rrrrrrrrrrrrrrrrrrrrBZbvji"
	NaN         = "rrrrrrrrrrrrrrrrrrrn5RM1rHd"
	ROOT        = "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
)

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
