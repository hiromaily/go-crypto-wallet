// Package eth provides Data Transfer Objects for Ethereum transactions.
package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// ETHMPCTransactionFile is the JSON file exchanged between wallets during a P4 MPC-TSS flow.
//
// State machine: TxType transitions from "unsigned" to "signed" when SignedTxHex is populated.
//
// File naming:
//
//	Unsigned: {action_type}_mpc_{uuid}.json
//	Signed:   {action_type}_mpc_{uuid}_signed.json
//
// This format is distinct from ETHTransactionFile (single-sig EOA) and
// ETHMultisigTransactionFile (Safe). No format mixing is permitted.
type ETHMPCTransactionFile struct {
	Version    int    `json:"version"`     // File format version (1)
	TxType     string `json:"tx_type"`     // "unsigned" or "signed"
	UUID       string `json:"uuid"`        // UUIDv4 generated at proposal time
	ActionType string `json:"action_type"` // "deposit", "payment", or "transfer"

	// Transaction parameters
	From     string `json:"from"`      // EIP-55 checksummed sender address (joint distributed EOA)
	To       string `json:"to"`        // EIP-55 checksummed recipient address
	Value    string `json:"value"`     // Wei as decimal string
	Nonce    uint64 `json:"nonce"`     // Ethereum account nonce
	GasLimit uint64 `json:"gas_limit"` // Gas limit for the transaction
	ChainID  uint64 `json:"chain_id"`  // EIP-155 chain identifier

	// EIP-1559 (Type 2) fee fields
	MaxFeePerGas         string `json:"max_fee_per_gas"`          // Wei as decimal string
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"` // Wei as decimal string

	// Signing material
	//
	// TxHash is the 0x-prefixed Keccak256 hash of the unsigned transaction.
	// This is the exact 32-byte pre-image that TSS nodes sign.
	// Nodes verify hash(RawTxHex) == TxHash before participating.
	TxHash   string `json:"tx_hash"`    // 0x-prefixed Keccak256 hash of the unsigned tx
	RawTxHex string `json:"raw_tx_hex"` // 0x-prefixed RLP-encoded unsigned transaction bytes

	// Filled after TSS signing completes
	SignedTxHex string `json:"signed_tx_hex,omitempty"` // 0x-prefixed signed transaction bytes

	// TSS configuration
	Threshold int      `json:"threshold"` // T — minimum signers required
	PartyIDs  []string `json:"party_ids"` // All N party IDs from the DKG ceremony
}

// Sentinel errors for ETHMPCTransactionFile validation failures.
var (
	// ErrInvalidMPCVersion is returned when version is less than 1.
	ErrInvalidMPCVersion = errors.New("version must be >= 1")
	// ErrInvalidMPCTxType is returned when tx_type is not "unsigned" or "signed".
	ErrInvalidMPCTxType = errors.New("tx_type must be \"unsigned\" or \"signed\"")
	// ErrInvalidMPCChainID is returned when chain_id is 0.
	ErrInvalidMPCChainID = errors.New("chain_id must be > 0")
	// ErrEmptyMPCUUID is returned when uuid is empty.
	ErrEmptyMPCUUID = errors.New("uuid is required")
	// ErrEmptyMPCFrom is returned when from address is empty.
	ErrEmptyMPCFrom = errors.New("from address is required")
	// ErrEmptyMPCTo is returned when to address is empty.
	ErrEmptyMPCTo = errors.New("to address is required")
	// ErrInvalidMPCAddress is returned when an address is not a valid EIP-55 checksummed hex address.
	ErrInvalidMPCAddress = errors.New("address must be a valid EIP-55 checksummed Ethereum address")
	// ErrEmptyMPCValue is returned when value is empty.
	ErrEmptyMPCValue = errors.New("value is required")
	// ErrZeroMPCGasLimit is returned when gas_limit is 0.
	ErrZeroMPCGasLimit = errors.New("gas_limit must be > 0")
	// ErrEmptyMPCTxHash is returned when tx_hash is empty.
	ErrEmptyMPCTxHash = errors.New("tx_hash is required")
	// ErrEmptyMPCRawTxHex is returned when raw_tx_hex is empty.
	ErrEmptyMPCRawTxHex = errors.New("raw_tx_hex is required")
	// ErrInvalidMPCThreshold is returned when threshold is less than 1.
	ErrInvalidMPCThreshold = errors.New("threshold must be >= 1")
	// ErrMPCPartyIDsBelowThreshold is returned when len(PartyIDs) < Threshold.
	ErrMPCPartyIDsBelowThreshold = errors.New("len(party_ids) must be >= threshold")
	// ErrMissingMPCSignedTxHex is returned when tx_type is "signed" but SignedTxHex is empty.
	ErrMissingMPCSignedTxHex = errors.New("signed_tx_hex is required when tx_type is \"signed\"")
	// ErrNotUnsigned is returned when a send operation is attempted on a file that is not unsigned.
	ErrNotUnsigned = errors.New("transaction file must have tx_type \"unsigned\" to initiate signing")
)

// Validate validates all invariants of the ETHMPCTransactionFile.
func (f *ETHMPCTransactionFile) Validate() error {
	if err := f.validateHeader(); err != nil {
		return err
	}
	if err := f.validateRequiredFields(); err != nil {
		return err
	}
	if f.TxType == string(domainTx.TxTypeSigned) && f.SignedTxHex == "" {
		return ErrMissingMPCSignedTxHex
	}
	return nil
}

func (f *ETHMPCTransactionFile) validateHeader() error {
	if f.Version < 1 {
		return ErrInvalidMPCVersion
	}
	if f.TxType != string(domainTx.TxTypeUnsigned) && f.TxType != string(domainTx.TxTypeSigned) {
		return ErrInvalidMPCTxType
	}
	if f.ChainID == 0 {
		return ErrInvalidMPCChainID
	}
	return nil
}

func (f *ETHMPCTransactionFile) validateRequiredFields() error {
	if f.UUID == "" {
		return ErrEmptyMPCUUID
	}
	if f.From == "" {
		return ErrEmptyMPCFrom
	}
	if !common.IsHexAddress(f.From) {
		return ErrInvalidMPCAddress
	}
	if f.To == "" {
		return ErrEmptyMPCTo
	}
	if !common.IsHexAddress(f.To) {
		return ErrInvalidMPCAddress
	}
	if f.Value == "" {
		return ErrEmptyMPCValue
	}
	if f.GasLimit == 0 {
		return ErrZeroMPCGasLimit
	}
	if f.TxHash == "" {
		return ErrEmptyMPCTxHash
	}
	if f.RawTxHex == "" {
		return ErrEmptyMPCRawTxHex
	}
	if f.Threshold < 1 {
		return ErrInvalidMPCThreshold
	}
	if len(f.PartyIDs) < f.Threshold {
		return ErrMPCPartyIDsBelowThreshold
	}
	return nil
}
