package btc

import (
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// BTCOnly defines methods specific to Bitcoin (BTC) that are NOT supported by BCH.
// These features require Bitcoin-specific protocols that BCH forked away from.
//
// Excluded from BCH:
//   - PSBT (BIP174): BCH uses raw transaction hex format
//   - Descriptors (BIP380-386): BCH uses ImportAddress for address watching
//   - SegWit (BIP141): BCH does not support SegWit
//   - Taproot (BIP340, BIP341, BIP342): BCH does not support Taproot
//   - MuSig2 (BIP327): BCH does not support Schnorr-based multisig
//
// BCH alternatives:
//   - Use SignRawTransactionWithKey instead of PSBT signing
//   - Use ImportAddress instead of ImportDescriptors
//   - Use P2SH multisig instead of P2WSH or MuSig2
//
// Usage:
//
//	For BTC use cases requiring full functionality, use the Bitcoiner interface
//	from interface.go, which includes all methods from both BitcoinCompatible and BTCOnly.
//
//	For BCH use cases, use BitcoinCompatible or BCHer interface instead.
type BTCOnly interface {
	// import.go (Descriptor-based imports - BTC only)
	ImportDescriptors(requests []dtobtc.ImportDescriptorsRequest) ([]dtobtc.ImportDescriptorsResponse, error)
	ImportMulti(
		requests []dtobtc.ImportMultiRequest,
		options *dtobtc.ImportMultiOptions,
	) ([]dtobtc.ImportMultiResponse, error)

	// descriptor_info.go (Descriptor operations - BTC only)
	GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error)
	ListDescriptors(privateDescriptors bool) ([]dtobtc.DescriptorInfo, error)

	// psbt.go (BIP174 Partially Signed Bitcoin Transaction - BTC only)
	// PSBT provides a standard format for partially signed transactions,
	// enabling multi-party signing workflows (hardware wallets, multisig, etc.)
	CreatePSBT(msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType) (string, error)
	ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error)
	ValidatePSBT(psbtBase64 string) error
	SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
	WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error) // RPC-based signing for descriptor wallets
	FinalizePSBT(psbtBase64 string) (string, error)
	ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
	IsPSBTComplete(psbtBase64 string) (bool, error)
	GetPSBTFee(psbtBase64 string) (int64, error)
}
