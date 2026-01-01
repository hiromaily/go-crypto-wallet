package bitcoindto

import (
	"github.com/btcsuite/btcd/btcutil"
)

// AddressInfo contains information about a Bitcoin address
type AddressInfo struct {
	Address        string
	ScriptPubKey   string
	IsWitness      bool
	WitnessVersion *int
	WitnessProg    string
	Labels         []string
	IsCompressed   bool
	IsChange       bool
	Timestamp      int64
	HDKeyPath      string
	HDMasterKeyID  string
	IsMine         bool
	IsWatchOnly    bool
	IsScript       bool
	PubKey         string
	Embedded       *AddressInfo
}

// ValidateAddressResult contains validation result for a Bitcoin address
type ValidateAddressResult struct {
	IsValid      bool
	Address      string
	ScriptPubKey string
	IsScript     bool
	IsWitness    bool
}

// NetworkInfo contains Bitcoin network information
type NetworkInfo struct {
	Version         int32
	SubVersion      string
	ProtocolVersion int32
	LocalServices   string
	LocalRelay      bool
	TimeOffset      int64
	NetworkActive   bool
	Connections     int32
	Networks        []NetworkAddress
	RelayFee        float64
	IncrementalFee  float64
	LocalAddresses  []LocalAddress
	Warnings        string
}

// NetworkAddress represents a network address
type NetworkAddress struct {
	Name                      string
	Limited                   bool
	Reachable                 bool
	Proxy                     string
	ProxyRandomizeCredentials bool
}

// LocalAddress represents a local address
type LocalAddress struct {
	Address string
	Port    uint16
	Score   int32
}

// BlockchainInfo contains blockchain information
type BlockchainInfo struct {
	Chain                string
	Blocks               int64
	Headers              int64
	BestBlockHash        string
	Difficulty           float64
	MedianTime           int64
	VerificationProgress float64
	InitialBlockDownload bool
	Pruned               bool
	SoftForks            map[string]any
	Warnings             string
}

// TransactionResult contains transaction information
type TransactionResult struct {
	Amount          btcutil.Amount
	Fee             btcutil.Amount
	Confirmations   int64
	BlockHash       string
	BlockIndex      int64
	BlockTime       int64
	TxID            string
	WalletConflicts []string
	Time            int64
	TimeReceived    int64
	Details         []TransactionDetail
	Hex             string
}

// TransactionDetail contains transaction detail information
type TransactionDetail struct {
	Address   string
	Category  string
	Amount    btcutil.Amount
	Label     string
	Vout      uint32
	Fee       *btcutil.Amount
	Abandoned bool
}

// RawTransaction contains raw transaction information
type RawTransaction struct {
	Hex           string
	TxID          string
	Hash          string
	Size          int32
	VSize         int32
	Weight        int32
	Version       int32
	LockTime      uint32
	Vin           []RawTransactionInput
	Vout          []RawTransactionOutput
	BlockHash     string
	Confirmations uint64
	Time          int64
	BlockTime     int64
}

// RawTransactionInput contains raw transaction input information
type RawTransactionInput struct {
	TxID      string
	Vout      uint32
	ScriptSig ScriptSig
	Sequence  uint32
	Witness   []string
}

// ScriptSig contains script signature information
type ScriptSig struct {
	Asm string
	Hex string
}

// RawTransactionOutput contains raw transaction output information
type RawTransactionOutput struct {
	Value        btcutil.Amount
	Index        uint32
	ScriptPubKey ScriptPubKey
}

// ScriptPubKey contains script public key information
type ScriptPubKey struct {
	Asm     string
	Hex     string
	ReqSigs int32
	Type    string
	Address string
}

// FundRawTransactionResult contains fund raw transaction result
type FundRawTransactionResult struct {
	Hex       string
	Fee       btcutil.Amount
	ChangePos int32
}

// PreviousTx contains previous transaction information for signing
type PreviousTx struct {
	TxID          string
	Vout          uint32
	ScriptPubKey  string
	RedeemScript  string
	WitnessScript string
	Amount        btcutil.Amount
}

// ParsedPSBT contains parsed PSBT information
type ParsedPSBT struct {
	Tx         ParsedPSBTTx
	Unknown    map[string]string
	Inputs     []ParsedPSBTInput
	Outputs    []ParsedPSBTOutput
	Fee        btcutil.Amount
	IsComplete bool
}

// ParsedPSBTTx contains parsed PSBT transaction information
type ParsedPSBTTx struct {
	TxID     string
	Hash     string
	Version  int32
	LockTime uint32
	Vin      []ParsedPSBTVin
	Vout     []ParsedPSBTVout
}

// ParsedPSBTVin contains parsed PSBT input information
type ParsedPSBTVin struct {
	TxID     string
	Vout     uint32
	Sequence uint32
}

// ParsedPSBTVout contains parsed PSBT output information
type ParsedPSBTVout struct {
	Value        btcutil.Amount
	ScriptPubKey string
}

// ParsedPSBTInput contains parsed PSBT input details
type ParsedPSBTInput struct {
	WitnessUTXO        *ParsedPSBTUTXO
	NonWitnessUTXO     *RawTransaction
	PartialSignatures  map[string]string
	SigHashType        uint32
	RedeemScript       string
	WitnessScript      string
	BIP32Derivation    []BIP32Derivation
	FinalScriptSig     string
	FinalScriptWitness []string
	Unknown            map[string]string
}

// ParsedPSBTOutput contains parsed PSBT output details
type ParsedPSBTOutput struct {
	RedeemScript    string
	WitnessScript   string
	BIP32Derivation []BIP32Derivation
	Unknown         map[string]string
}

// ParsedPSBTUTXO contains parsed PSBT UTXO information
type ParsedPSBTUTXO struct {
	Amount       btcutil.Amount
	ScriptPubKey string
}

// BIP32Derivation contains BIP32 derivation path information
type BIP32Derivation struct {
	PubKey      string
	MasterKeyID string
	Path        string
}

// UnspentOutput contains unspent output information
type UnspentOutput struct {
	TxID          string
	Vout          uint32
	Address       string
	Account       string
	ScriptPubKey  string
	Amount        btcutil.Amount
	Confirmations int64
	RedeemScript  string
	WitnessScript string
	Spendable     bool
	Solvable      bool
	Safe          bool
	Label         string
}

// LoggingResult contains logging information
type LoggingResult struct {
	Net                 bool
	TorControl          bool
	MemPool             bool
	RPC                 bool
	EstimateFee         bool
	Coindb              bool
	SelectCoins         bool
	Bench               bool
	ZMQPubRawBlock      bool
	ZMQPubRawTx         bool
	ZMQPubHashBlock     bool
	ZMQPubHashTx        bool
	ValidationInterface bool
	Reindex             bool
	Proxy               bool
}

// MultisigAddress contains multisig address information
type MultisigAddress struct {
	Address      string
	RedeemScript string
	Descriptor   string
}
