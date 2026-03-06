package xrp

// SentTx is the application DTO for a submitted XRP transaction result.
type SentTx struct {
	ResultCode    string
	ResultMessage string
	TxBlob        string
	TxJSON        TxInput
}

// TxOutcome is part of TxInfo representing the outcome of a transaction.
type TxOutcome struct {
	Result        string
	LedgerVersion int
}

// TxInfo is the application DTO for retrieved XRP transaction details.
type TxInfo struct {
	ID      string
	Outcome TxOutcome
}

// LedgerCurrent is the application DTO for the current ledger index.
type LedgerCurrent struct {
	LedgerCurrentIndex uint64
}

// LedgerAccept is the application DTO for a ledger_accept result.
type LedgerAccept struct {
	LedgerCurrentIndex uint64
}

// AccountInfo is the application DTO for XRP account state.
type AccountInfo struct {
	Sequence                                  uint64
	XrpBalance                                string
	OwnerCount                                uint64
	PreviousAffectingTransactionID            string
	PreviousAffectingTransactionLedgerVersion uint64
}

// WalletInfo is the application DTO for a proposed XRP wallet.
type WalletInfo struct {
	AccountID     string
	KeyType       string
	MasterKey     string // deprecated by XRP Ledger; preserved for backward compatibility
	MasterSeed    string
	MasterSeedHex string
	PublicKey     string
	PublicKeyHex  string
}

// ValidationKey is the application DTO for a validation key set.
type ValidationKey struct {
	ValidationKey       string
	ValidationPublicKey string
	ValidationSeed      string
}
