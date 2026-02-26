// Package eth provides Data Transfer Objects for Ethereum transactions.
package eth

import (
	"errors"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// ETHTransactionFile represents a JSON file structure for ETH transaction exchange
// between the Watch wallet (online) and Keygen wallet (offline).
// This format supports both legacy (Type 0) and EIP-1559 (Type 2) transactions.
//
// File naming convention: {action}_{txID}_{type}_{signedCount}_{timestamp}.json
// Example: deposit_8_unsigned_0_1534744535097796209.json
//
// RawTxHex uses tx.MarshalBinary() (EIP-2718 envelope) — not RLP — to correctly
// encode both legacy and EIP-1559 transactions.
type ETHTransactionFile struct {
	Version   int    `json:"version"`     // File format version (1)
	TxType    string `json:"tx_type"`     // "unsigned" or "signed"
	EthTxType uint8  `json:"eth_tx_type"` // 0=legacy, 2=EIP-1559

	// Transaction identity
	UUID    string `json:"uuid"`     // Internal UUID for database record lookup after broadcast
	ChainID uint64 `json:"chain_id"` // EIP-155 chain identifier
	Nonce   uint64 `json:"nonce"`    // Sender transaction count
	From    string `json:"from"`     // Sender address
	To      string `json:"to"`       // Recipient address
	Value   string `json:"value"`    // Wei amount as decimal string
	Gas     uint64 `json:"gas"`      // Gas limit

	// Legacy (Type 0) fee fields
	GasPrice string `json:"gas_price,omitempty"` // Wei as decimal string

	// EIP-1559 (Type 2) fee fields
	MaxFeePerGas         string `json:"max_fee_per_gas,omitempty"`          // Wei as decimal string
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas,omitempty"` // Wei as decimal string

	// Call data
	Data string `json:"data,omitempty"` // Hex-encoded input data

	// Encoded transaction (EIP-2718 binary format via tx.MarshalBinary())
	RawTxHex    string `json:"raw_tx_hex"`              // Always present; unsigned tx binary
	SignedTxHex string `json:"signed_tx_hex,omitempty"` // Present only when tx_type == "signed"
}

// EthTxType constants for Ethereum transaction type identifiers.
const (
	EthTxTypeLegacy  uint8 = 0 // Legacy (Type 0) transaction
	EthTxTypeEIP1559 uint8 = 2 // EIP-1559 (Type 2) dynamic fee transaction
)

// Sentinel errors for validation failures.
var (
	// ErrInvalidETHVersion is returned when version is less than 1
	ErrInvalidETHVersion = errors.New("version must be >= 1")
	// ErrInvalidETHTxType is returned when tx_type is not "unsigned" or "signed"
	ErrInvalidETHTxType = errors.New("tx_type must be \"unsigned\" or \"signed\"")
	// ErrInvalidETHTxTypeCode is returned when eth_tx_type is not 0 (legacy) or 2 (EIP-1559)
	ErrInvalidETHTxTypeCode = errors.New("eth_tx_type must be 0 (legacy) or 2 (EIP-1559)")
	// ErrInvalidETHChainID is returned when chain_id is 0
	ErrInvalidETHChainID = errors.New("chain_id must be > 0")
	// ErrEmptyETHFrom is returned when from address is empty
	ErrEmptyETHFrom = errors.New("from address is required")
	// ErrEmptyETHTo is returned when to address is empty
	ErrEmptyETHTo = errors.New("to address is required")
	// ErrEmptyETHValue is returned when value is empty
	ErrEmptyETHValue = errors.New("value is required")
	// ErrZeroETHGas is returned when gas limit is 0
	ErrZeroETHGas = errors.New("gas must be > 0")
	// ErrEmptyETHRawTxHex is returned when raw_tx_hex is empty
	ErrEmptyETHRawTxHex = errors.New("raw_tx_hex is required")
	// ErrMissingETHGasPrice is returned when legacy tx is missing gas_price
	ErrMissingETHGasPrice = errors.New("gas_price is required for legacy (Type 0) transactions")
	// ErrMissingETHMaxFeePerGas is returned when EIP-1559 tx is missing max_fee_per_gas
	ErrMissingETHMaxFeePerGas = errors.New("max_fee_per_gas is required for EIP-1559 (Type 2) transactions")
	// ErrMissingETHMaxPriorityFeePerGas is returned when EIP-1559 tx is missing max_priority_fee_per_gas
	ErrMissingETHMaxPriorityFeePerGas = errors.New(
		"max_priority_fee_per_gas is required for EIP-1559 (Type 2) transactions",
	)
	// ErrMissingETHSignedTxHex is returned when signed tx is missing signed_tx_hex
	ErrMissingETHSignedTxHex = errors.New("signed_tx_hex is required when tx_type is \"signed\"")
)

// Validate validates the ETHTransactionFile and enforces all invariants.
func (f *ETHTransactionFile) Validate() error {
	if err := f.validateHeader(); err != nil {
		return err
	}
	if err := f.validateTransactionFields(); err != nil {
		return err
	}
	if err := f.validateFeeFields(); err != nil {
		return err
	}
	if f.TxType == string(domainTx.TxTypeSigned) && f.SignedTxHex == "" {
		return ErrMissingETHSignedTxHex
	}
	return nil
}

// validateHeader checks version, tx type, eth tx type, and chain ID.
func (f *ETHTransactionFile) validateHeader() error {
	if f.Version < 1 {
		return ErrInvalidETHVersion
	}
	if f.TxType != string(domainTx.TxTypeUnsigned) && f.TxType != string(domainTx.TxTypeSigned) {
		return ErrInvalidETHTxType
	}
	if f.EthTxType != EthTxTypeLegacy && f.EthTxType != EthTxTypeEIP1559 {
		return ErrInvalidETHTxTypeCode
	}
	if f.ChainID == 0 {
		return ErrInvalidETHChainID
	}
	return nil
}

// validateTransactionFields checks required address, value, gas, and raw tx fields.
func (f *ETHTransactionFile) validateTransactionFields() error {
	if f.From == "" {
		return ErrEmptyETHFrom
	}
	if f.To == "" {
		return ErrEmptyETHTo
	}
	if f.Value == "" {
		return ErrEmptyETHValue
	}
	if f.Gas == 0 {
		return ErrZeroETHGas
	}
	if f.RawTxHex == "" {
		return ErrEmptyETHRawTxHex
	}
	return nil
}

// validateFeeFields checks required fee fields based on the transaction type.
func (f *ETHTransactionFile) validateFeeFields() error {
	switch f.EthTxType {
	case EthTxTypeLegacy:
		if f.GasPrice == "" {
			return ErrMissingETHGasPrice
		}
	case EthTxTypeEIP1559:
		if f.MaxFeePerGas == "" {
			return ErrMissingETHMaxFeePerGas
		}
		if f.MaxPriorityFeePerGas == "" {
			return ErrMissingETHMaxPriorityFeePerGas
		}
	}
	return nil
}
