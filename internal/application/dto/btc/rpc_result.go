package btc

// AddressInfo contains address information returned from getaddressinfo RPC.
// Replaces btcrpc.GetAddressInfoResult in port interface signatures.
type AddressInfo struct {
	Address             string
	ScriptPubKey        string
	IsMine              bool
	Solvable            bool
	Desc                string
	IsWatchOnly         bool
	IsScript            bool
	IsWitness           bool
	Pubkey              string
	IsCompressed        bool
	IsChange            bool
	HDKeyPath           string
	HDMasterFingerprint string
	Labels              []string
	// Label is a legacy single-string field returned by older BCH nodes.
	// Use Labels instead; this field exists only for BCH fallback logic.
	Label string
}

// ValidateAddressResult contains address validation results from validateaddress RPC.
// Replaces btcrpc.ValidateAddressResult in port interface signatures.
type ValidateAddressResult struct {
	IsValid           bool
	Address           string
	ScriptPubKey      string
	IsScript          bool
	IsWitness         bool
	WitnessVersion    uint32
	WitnessProgramHex string
}

// MultisigAddress contains the result of addmultisigaddress RPC.
// Replaces btcrpc.AddMultisigAddressResult in port interface signatures.
type MultisigAddress struct {
	Address      string
	RedeemScript string
}

// ImportMultiRequest is a single request for the importmulti RPC.
// Replaces btcrpc.ImportMultiRequest in port interface signatures.
type ImportMultiRequest struct {
	ScriptPubKey any
	Timestamp    any
	RedeemScript string
	PubKeys      []string
	Keys         []string
	Internal     bool
	WatchOnly    bool
	Label        string
}

// ImportMultiOptions are options for the importmulti RPC.
// Replaces btcrpc.ImportMultiOptions in port interface signatures.
type ImportMultiOptions struct {
	Rescan bool
}

// ImportMultiResponse is the response for a single importmulti request.
// Replaces btcrpc.ImportMultiResponse in port interface signatures.
type ImportMultiResponse struct {
	Success  bool
	Warnings []string
	Error    *ImportMultiError
}

// ImportMultiError contains error details for a failed importmulti request.
type ImportMultiError struct {
	Code    int
	Message string
}

// DescriptorInfo contains the result of getdescriptorinfo RPC or a single entry from listdescriptors.
// Replaces btcrpc.DescriptorInfo and btcrpc.DescriptorEntry in port interface signatures.
type DescriptorInfo struct {
	Descriptor     string
	Checksum       string
	IsRange        bool
	IsSolvable     bool
	HasPrivateKeys bool
	// Internal is set for descriptors from listdescriptors; nil means not applicable.
	Internal *bool
}

// ImportDescriptorsRequest is a single descriptor import request for importdescriptors.
// Replaces btcrpc.ImportDescriptorsRequest in port interface signatures.
type ImportDescriptorsRequest struct {
	Descriptor string
	Timestamp  any
	Active     bool
	Range      any
	Label      string
	Internal   bool
	Watchonly  bool
}

// ImportDescriptorsResponse is the response for a single importdescriptors request.
// Replaces btcrpc.ImportDescriptorsResponse in port interface signatures.
type ImportDescriptorsResponse struct {
	Success  bool
	Warnings []string
	Error    *ImportDescriptorsError
}

// ImportDescriptorsError contains error details for a failed descriptor import.
type ImportDescriptorsError struct {
	Code    int
	Message string
}

// TransactionInput is a transaction input for createrawtransaction.
// Replaces btcjson.TransactionInput in port interface signatures.
type TransactionInput struct {
	TxID     string
	Vout     uint32
	Sequence *uint32
}

// TransactionResult contains transaction details from gettransaction RPC.
// Replaces btcrpc.GetTransactionResult in port interface signatures.
type TransactionResult struct {
	Amount        float64
	Fee           float64
	Confirmations uint64
	BlockHash     string
	BlockHeight   uint64
	Txid          string
	Time          int64
	Details       []TransactionDetail
	Hex           string
}

// TransactionDetail is a detail entry in TransactionResult.
type TransactionDetail struct {
	Address  string
	Category string
	Amount   float64
	Label    string
	Vout     uint32
}

// NetworkInfo contains network information from getnetworkinfo RPC.
// Replaces btcrpc.GetNetworkInfoResult in port interface signatures.
type NetworkInfo struct {
	Version         int32
	Subversion      string
	ProtocolVersion uint32
	Connections     uint32
	Networks        []NetworkInfoNetwork
	RelayFee        float64
	LocalAddresses  []NetworkInfoLocalAddress
	Warnings        []string
}

// NetworkInfoNetwork is a network entry in NetworkInfo.
type NetworkInfoNetwork struct {
	Name      string
	Limited   bool
	Reachable bool
	Proxy     string
}

// NetworkInfoLocalAddress is a local address entry in NetworkInfo.
type NetworkInfoLocalAddress struct {
	Address string
	Port    int
	Score   int
}

// LoggingStatus contains logging configuration from the logging RPC.
// Replaces btcrpc.LoggingResult in port interface signatures.
type LoggingStatus struct {
	Net         bool
	Tor         bool
	Mempool     bool
	HTTP        bool
	Bench       bool
	Zmq         bool
	Walletdb    bool
	RPC         bool
	Estimatefee bool
	Addrman     bool
	Selectcoins bool
	Reindex     bool
	Cmpctblock  bool
	Rand        bool
	Prune       bool
	Proxy       bool
	Mempoolrej  bool
	Libevent    bool
	Coindb      bool
	Qt          bool
	Leveldb     bool
	Validation  bool
}
