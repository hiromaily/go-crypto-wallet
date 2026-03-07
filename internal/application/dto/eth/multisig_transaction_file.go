// Package eth provides Data Transfer Objects for Ethereum transactions.
package eth

import (
	"errors"
	"strings"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// SignEntry holds a single signer's address and their Safe-compatible EIP-712 signature.
type SignEntry struct {
	SignerAddress string `json:"signer_address"` // EIP-55 checksummed EOA address
	SignatureHex  string `json:"signature_hex"`  // 0x-prefixed 65-byte hex; v adjusted to 27/28
}

// ETHMultisigTransactionFile represents the JSON file exchanged between wallets during
// a Safe (Gnosis Safe v1.4.1) multi-signature ETH transaction workflow.
//
// State machine: TxType transitions from "unsigned" to "signed" when len(Signatures) >= Threshold.
//
// File naming convention: {actionType}_multisig_{uuid}_{signedCount}.json
// Example: payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json
//
// All gas fields (SafeTxGas, BaseGas, GasPrice, GasToken, RefundReceiver) default to zero/empty
// for simple ETH transfers. The safe nonce is fetched at proposal time via the Safe contract.
type ETHMultisigTransactionFile struct {
	Version int    `json:"version"` // File format version (1)
	TxType  string `json:"tx_type"` // "unsigned" or "signed"
	UUID    string `json:"uuid"`    // Unique proposal identifier (UUIDv4); used in file names

	// Safe contract identity
	SafeAddress string `json:"safe_address"` // EIP-55 checksummed Safe proxy address

	// Safe execTransaction parameters
	To             string `json:"to"`              // Recipient address (EIP-55 checksum)
	Value          string `json:"value"`           // Wei amount as decimal string
	Data           string `json:"data"`            // 0x-prefixed hex call data (empty = "0x" for ETH transfer)
	Operation      uint8  `json:"operation"`       // 0=Call, 1=DelegateCall
	SafeTxGas      string `json:"safe_tx_gas"`     // Wei decimal string (0 for ETH transfer)
	BaseGas        string `json:"base_gas"`        // Wei decimal string (0 for ETH transfer)
	GasPrice       string `json:"gas_price"`       // Wei decimal string (0 for ETH transfer)
	GasToken       string `json:"gas_token"`       // Token address for gas (zero address for ETH)
	RefundReceiver string `json:"refund_receiver"` // Refund receiver address (zero address)
	Nonce          string `json:"nonce"`           // Safe contract nonce as decimal string

	// Chain and hash
	ChainID    uint64 `json:"chain_id"`     // EIP-155 chain identifier
	SafeTxHash string `json:"safe_tx_hash"` // EIP-712 safeTxHash as 0x-prefixed hex (authoritative)

	// Signature accumulation
	Threshold  int         `json:"threshold"`  // Required signer count (m in m-of-n)
	Signatures []SignEntry `json:"signatures"` // Accumulated signer entries; empty until first signing
}

// Sentinel errors for multisig file validation failures.
var (
	// ErrSafeTxHashMismatch is returned when the recomputed EIP-712 hash differs from the stored hash.
	ErrSafeTxHashMismatch = errors.New("recomputed safeTxHash does not match stored safe_tx_hash")
	// ErrDuplicateSigner is returned when the same signer address appears more than once.
	ErrDuplicateSigner = errors.New("duplicate signer address in signatures")
	// ErrNotFullySigned is returned when a send is attempted but TxType is still "unsigned".
	ErrNotFullySigned = errors.New("transaction is not fully signed (tx_type must be \"signed\")")

	// ErrInvalidMultisigVersion is returned when version is less than 1.
	ErrInvalidMultisigVersion = errors.New("version must be >= 1")
	// ErrInvalidMultisigTxType is returned when tx_type is not "unsigned" or "signed".
	ErrInvalidMultisigTxType = errors.New("tx_type must be \"unsigned\" or \"signed\"")
	// ErrInvalidMultisigChainID is returned when chain_id is 0.
	ErrInvalidMultisigChainID = errors.New("chain_id must be > 0")
	// ErrEmptyMultisigUUID is returned when uuid is empty.
	ErrEmptyMultisigUUID = errors.New("uuid is required")
	// ErrEmptyMultisigSafeAddress is returned when safe_address is empty.
	ErrEmptyMultisigSafeAddress = errors.New("safe_address is required")
	// ErrEmptyMultisigTo is returned when to address is empty.
	ErrEmptyMultisigTo = errors.New("to address is required")
	// ErrEmptyMultisigValue is returned when value is empty.
	ErrEmptyMultisigValue = errors.New("value is required")
	// ErrEmptyMultisigNonce is returned when nonce is empty.
	ErrEmptyMultisigNonce = errors.New("nonce is required")
	// ErrEmptyMultisigSafeTxHash is returned when safe_tx_hash is empty.
	ErrEmptyMultisigSafeTxHash = errors.New("safe_tx_hash is required")
	// ErrInvalidMultisigThreshold is returned when threshold is less than 1.
	ErrInvalidMultisigThreshold = errors.New("threshold must be >= 1")
	// ErrSignatureCountExceedsThreshold is returned when more signatures exist than the threshold.
	ErrSignatureCountExceedsThreshold = errors.New("signature count cannot exceed threshold")
	// ErrSignedTxTypeInconsistent is returned when tx_type is "signed" but signature count < threshold.
	ErrSignedTxTypeInconsistent = errors.New("tx_type is \"signed\" but signature count is below threshold")
	// ErrUnsignedTxTypeInconsistent is returned when tx_type is "unsigned" but signature count >= threshold.
	ErrUnsignedTxTypeInconsistent = errors.New("tx_type is \"unsigned\" but signature count meets threshold")
)

// Validate validates the ETHMultisigTransactionFile and enforces all invariants.
func (f *ETHMultisigTransactionFile) Validate() error {
	if err := f.validateHeader(); err != nil {
		return err
	}
	if err := f.validateRequiredFields(); err != nil {
		return err
	}
	if err := f.validateSignatures(); err != nil {
		return err
	}
	return f.validateTxTypeConsistency()
}

func (f *ETHMultisigTransactionFile) validateHeader() error {
	if f.Version < 1 {
		return ErrInvalidMultisigVersion
	}
	if f.TxType != string(domainTx.TxTypeUnsigned) && f.TxType != string(domainTx.TxTypeSigned) {
		return ErrInvalidMultisigTxType
	}
	if f.ChainID == 0 {
		return ErrInvalidMultisigChainID
	}
	return nil
}

func (f *ETHMultisigTransactionFile) validateRequiredFields() error {
	if f.UUID == "" {
		return ErrEmptyMultisigUUID
	}
	if f.SafeAddress == "" {
		return ErrEmptyMultisigSafeAddress
	}
	if f.To == "" {
		return ErrEmptyMultisigTo
	}
	if f.Value == "" {
		return ErrEmptyMultisigValue
	}
	if f.Nonce == "" {
		return ErrEmptyMultisigNonce
	}
	if f.SafeTxHash == "" {
		return ErrEmptyMultisigSafeTxHash
	}
	if f.Threshold < 1 {
		return ErrInvalidMultisigThreshold
	}
	return nil
}

func (f *ETHMultisigTransactionFile) validateSignatures() error {
	if len(f.Signatures) > f.Threshold {
		return ErrSignatureCountExceedsThreshold
	}
	// Check for duplicate signer addresses (case-insensitive)
	seen := make(map[string]struct{}, len(f.Signatures))
	for _, entry := range f.Signatures {
		key := strings.ToLower(entry.SignerAddress)
		if _, exists := seen[key]; exists {
			return ErrDuplicateSigner
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (f *ETHMultisigTransactionFile) validateTxTypeConsistency() error {
	sigCount := len(f.Signatures)
	isSigned := f.TxType == string(domainTx.TxTypeSigned)
	if isSigned && sigCount < f.Threshold {
		return ErrSignedTxTypeInconsistent
	}
	if !isSigned && sigCount >= f.Threshold {
		return ErrUnsignedTxTypeInconsistent
	}
	return nil
}
