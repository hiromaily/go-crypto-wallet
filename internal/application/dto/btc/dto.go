package btc

import (
	"github.com/btcsuite/btcd/btcutil"
)

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
