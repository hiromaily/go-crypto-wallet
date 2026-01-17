// Package bch provides Bitcoin Cash specific implementations.
//
// This file contains overridden methods that are NOT supported by BCH.
// These methods will return an error if called, preventing accidental use
// of BTC-only features (PSBT, Descriptor, MuSig2) in BCH context.
package bch

import (
	"errors"

	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// ErrBCHUnsupported is the base error for unsupported BCH operations.
var ErrBCHUnsupported = errors.New("BCH unsupported")

// errPSBTNotSupported is returned when PSBT methods are called on BCH.
// BCH uses Raw Transaction Hex format instead of PSBT (BIP174).
var errPSBTNotSupported = errors.New("PSBT (BIP174) is not supported by BCH - use Raw Transaction Hex format instead")

// errDescriptorNotSupported is returned when Descriptor methods are called on BCH.
// BCH does not support Output Descriptors.
var errDescriptorNotSupported = errors.New(
	"output descriptors are not supported by BCH - use ImportAddress methods instead",
)

// =============================================================================
// PSBT Methods - NOT SUPPORTED BY BCH
// uses Raw Transaction Hex format, not PSBT (BIP174)
// =============================================================================

// CreatePSBT is not supported by BCH.
// BCH should use CreateRawTransaction instead.
func (*BitcoinCash) CreatePSBT(
	_ *wire.MsgTx,
	_ []dtobtc.PreviousTx,
	_ domainAccount.AccountType,
) (string, error) {
	return "", errPSBTNotSupported
}

// ParsePSBT is not supported by BCH.
func (*BitcoinCash) ParsePSBT(_ string) (*dtobtc.ParsedPSBT, error) {
	return nil, errPSBTNotSupported
}

// ValidatePSBT is not supported by BCH.
func (*BitcoinCash) ValidatePSBT(_ string) error {
	return errPSBTNotSupported
}

// SignPSBTWithKey is not supported by BCH.
// BCH should use SignRawTransactionWithKey with SIGHASH_FORKID (0x41) instead.
func (*BitcoinCash) SignPSBTWithKey(_ string, _ []string) (string, bool, error) {
	return "", false, errPSBTNotSupported
}

// WalletProcessPsbt is not supported by BCH.
func (*BitcoinCash) WalletProcessPsbt(_ string, _ bool) (string, bool, error) {
	return "", false, errPSBTNotSupported
}

// FinalizePSBT is not supported by BCH.
func (*BitcoinCash) FinalizePSBT(_ string) (string, error) {
	return "", errPSBTNotSupported
}

// ExtractTransaction is not supported by BCH.
// BCH transactions are already in raw hex format, no extraction needed.
func (*BitcoinCash) ExtractTransaction(_ string) (*wire.MsgTx, error) {
	return nil, errPSBTNotSupported
}

// IsPSBTComplete is not supported by BCH.
func (*BitcoinCash) IsPSBTComplete(_ string) (bool, error) {
	return false, errPSBTNotSupported
}

// GetPSBTFee is not supported by BCH.
// BCH should use GetTransactionFee with the raw transaction instead.
func (*BitcoinCash) GetPSBTFee(_ string) (int64, error) {
	return 0, errPSBTNotSupported
}

// =============================================================================
// Descriptor Methods - NOT SUPPORTED BY BCH
// does not support Output Descriptors
// =============================================================================

// ImportDescriptors is not supported by BCH.
// BCH should use ImportAddress or ImportAddressWithLabel instead.
func (*BitcoinCash) ImportDescriptors(
	_ []dtobtc.ImportDescriptorsRequest,
) ([]dtobtc.ImportDescriptorsResponse, error) {
	return nil, errDescriptorNotSupported
}

// GetDescriptorInfo is not supported by BCH.
func (*BitcoinCash) GetDescriptorInfo(_ string) (*dtobtc.DescriptorInfo, error) {
	return nil, errDescriptorNotSupported
}

// ListDescriptors is not supported by BCH.
func (*BitcoinCash) ListDescriptors(_ bool) (*dtobtc.ListDescriptorsResult, error) {
	return nil, errDescriptorNotSupported
}
