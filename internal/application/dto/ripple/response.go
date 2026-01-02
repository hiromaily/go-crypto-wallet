package ripple

// ResponseGetAccountInfo represents the response from GetAccountInfo API call.
type ResponseGetAccountInfo struct {
	Sequence                                  uint64
	XrpBalance                                string
	OwnerCount                                uint64
	PreviousAffectingTransactionID            string
	PreviousAffectingTransactionLedgerVersion uint64
}

// ResponseGenerateAddress represents the response from GenerateAddress API call.
type ResponseGenerateAddress struct {
	XAddress       string
	ClassicAddress string
	Address        string
	Secret         string
}

// ResponseGenerateXAddress represents the response from GenerateXAddress API call.
type ResponseGenerateXAddress struct {
	XAddress string
	Secret   string
}

// ResponseAccountChannels represents the response from AccountChannels API call.
type ResponseAccountChannels struct {
	Account   string
	Channels  []AccountChannel
	Validated bool
}

// AccountChannel represents a payment channel in the AccountChannels response.
type AccountChannel struct {
	ChannelID      string
	Account        string
	Destination    string
	Amount         string
	Balance        string
	SettleDelay    uint64
	PublicKey      string
	DestinationTag uint32
	CancelAfter    uint64
	Expiration     uint64
}

// ResponseAccountInfo represents the response from AccountInfo API call.
type ResponseAccountInfo struct {
	Account            string
	Balance            string
	Sequence           uint64
	OwnerCount         uint64
	Flags              uint64
	LedgerCurrentIndex uint64
	Validated          bool
}

// ResponseServerInfo represents the response from ServerInfo API call.
type ResponseServerInfo struct {
	BuildVersion     string
	CompleteLedgers  string
	HostID           string
	IOLatencyMS      uint64
	LastClose        ServerLastClose
	LoadFactor       uint64
	Peers            uint64
	PubkeyNode       string
	ServerState      string
	ValidatedLedger  ValidatedLedger
	ValidationQuorum uint64
}

// ServerLastClose represents the last close information in ServerInfo response.
type ServerLastClose struct {
	ConvergeTimeS float64
	Proposers     uint64
}

// ValidatedLedger represents validated ledger information in ServerInfo response.
type ValidatedLedger struct {
	Age            uint64
	BaseFeeXRP     string
	Hash           string
	ReserveBaseXRP string
	ReserveIncXRP  string
	Seq            uint64
}

// ResponseValidationCreate represents the response from ValidationCreate API call.
type ResponseValidationCreate struct {
	ValidationPublicKey string
	ValidationSeed      string
	ValidationKey       string
}

// ResponseWalletPropose represents the response from WalletPropose API call.
type ResponseWalletPropose struct {
	MasterSeed    string
	MasterSeedHex string
	MasterKey     string
	AccountID     string
	PublicKey     string
	PublicKeyHex  string
	KeyType       string
	Warning       string
}
