package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/quagmt/udecimal"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// =============================================================================
// Small Interfaces - Interface Segregation Principle (ISP)
// =============================================================================
//
// These interfaces are designed to be minimal and focused on specific use cases.
// Usecases should depend on these small interfaces instead of the large Bitcoiner interface.
//
// Benefits:
// - Usecases depend only on methods they actually use
// - Easier to mock in tests
// - Clearer dependency requirements
// - Safer refactoring with limited blast radius
//
// Note: The full Bitcoiner interface is kept for backward compatibility.
// New code should prefer these small interfaces.
// =============================================================================

// -----------------------------------------------------------------------------
// Shared Utility Interfaces
// -----------------------------------------------------------------------------

// ChainConfigProvider provides chain configuration and coin type information.
// Used by most usecases that need network-specific parameters.
type ChainConfigProvider interface {
	GetChainConf() *chaincfg.Params
	CoinTypeCode() domainCoin.CoinTypeCode
	ConfirmationBlock() uint64
}

// AmountConverter handles amount conversion between different formats.
// Used by transaction creation and database operations.
type AmountConverter interface {
	FloatToAmount(f float64) (btcutil.Amount, error)
	FloatToDecimal(f float64) (udecimal.Decimal, error)
	AmountToDecimal(amt btcutil.Amount) (udecimal.Decimal, error)
}

// -----------------------------------------------------------------------------
// Transaction Creation Interfaces
// -----------------------------------------------------------------------------

// UTXOProvider provides UTXO (Unspent Transaction Output) information.
// Used by transaction creation usecases.
type UTXOProvider interface {
	ListUnspentByAccount(accountType domainAccount.AccountType, confirmationNum uint64) ([]dtobtc.UnspentOutput, error)
	GetUnspentListAddrs(unspentList []dtobtc.UnspentOutput, accountType domainAccount.AccountType) []string
}

// RawTransactionCreator creates and manipulates raw transactions.
// Used by watch wallet transaction creation.
type RawTransactionCreator interface {
	CreateRawTransaction(
		inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
	) (*wire.MsgTx, error)
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)
}

// AddressOperator provides address-related operations.
// Used by transaction creation and import usecases.
type AddressOperator interface {
	DecodeAddress(addr string) (btcutil.Address, error)
	GetAddressInfo(addr string) (*dtobtc.AddressInfo, error)
	GetAccount(addr string) (string, error)
}

// BalanceChecker provides balance information.
// Used by transaction creation and monitoring usecases.
type BalanceChecker interface {
	GetBalanceByAccount(accountType domainAccount.AccountType, confirmationNum uint64) (btcutil.Amount, error)
}

// -----------------------------------------------------------------------------
// Transaction Broadcast Interfaces
// -----------------------------------------------------------------------------

// TransactionSender broadcasts transactions to the network.
// Used by watch wallet send transaction usecase.
type TransactionSender interface {
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
}

// TransactionMonitor provides transaction monitoring capabilities.
// Used by watch wallet monitor usecase.
type TransactionMonitor interface {
	GetTransactionByTxID(txID string) (*dtobtc.TransactionResult, error)
}

// -----------------------------------------------------------------------------
// Import Interfaces
// -----------------------------------------------------------------------------

// AddressImporter imports addresses into the wallet.
// Used by watch wallet import usecase for both BTC and BCH.
type AddressImporter interface {
	ImportAddressWithLabel(addr, label string, rescan bool) error
}

// LegacyAddressImporter provides legacy import methods.
// Used for backward compatibility with older wallet formats.
type LegacyAddressImporter interface {
	AddressImporter
	ImportMulti(
		requests []dtobtc.ImportMultiRequest, options *dtobtc.ImportMultiOptions,
	) ([]dtobtc.ImportMultiResponse, error)
}

// PrivateKeyImporter imports private keys into the wallet.
// Used by keygen and sign wallets.
type PrivateKeyImporter interface {
	ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error
}

// -----------------------------------------------------------------------------
// BTC-Only Interfaces (PSBT, Descriptors)
// -----------------------------------------------------------------------------

// DescriptorManager manages descriptor-based wallet operations (BTC only).
// BCH does not support descriptors - use AddressImporter instead.
type DescriptorManager interface {
	GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error)
	ImportDescriptors(requests []dtobtc.ImportDescriptorsRequest) ([]dtobtc.ImportDescriptorsResponse, error)
	SetLabel(addr, label string) error
}

// PSBTCreator creates PSBT (Partially Signed Bitcoin Transactions) - BTC only.
// BCH uses raw transaction hex format instead.
type PSBTCreator interface {
	CreatePSBT(msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType) (string, error)
}

// PSBTSigner signs PSBT transactions - BTC only.
// BCH uses SignRawTransactionWithKey instead.
type PSBTSigner interface {
	ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error)
	SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
	WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error)
}

// PSBTFinalizer finalizes and extracts transactions from PSBT - BTC only.
type PSBTFinalizer interface {
	IsPSBTComplete(psbtBase64 string) (bool, error)
	FinalizePSBT(psbtBase64 string) (string, error)
	ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
}

// PSBTHandler combines all PSBT operations - BTC only.
// Use this when a usecase needs full PSBT lifecycle management.
type PSBTHandler interface {
	PSBTCreator
	PSBTSigner
	PSBTFinalizer
}

// -----------------------------------------------------------------------------
// Transaction Signing Interfaces
// -----------------------------------------------------------------------------

// RawTransactionSigner signs raw transactions.
// Used by BCH (with SIGHASH_FORKID) and legacy BTC signing.
type RawTransactionSigner interface {
	SignRawTransactionWithKey(
		tx *wire.MsgTx, privKeysWIF []string, prevtxs []dtobtc.PreviousTx,
	) (*wire.MsgTx, bool, error)
}

// -----------------------------------------------------------------------------
// Multisig Interfaces
// -----------------------------------------------------------------------------

// MultisigManager manages multisig address creation.
// Used by keygen wallet for creating multisig addresses.
type MultisigManager interface {
	AddMultisigAddress(
		requiredSigs int, addresses []string, accountName string, addressType domainBitcoin.AddressType,
	) (*dtobtc.MultisigAddress, error)
}

// -----------------------------------------------------------------------------
// Composed Interfaces for Common Use Cases
// -----------------------------------------------------------------------------

// TransactionCreationDeps combines interfaces needed for transaction creation.
// Use this for watch wallet create transaction usecase.
type TransactionCreationDeps interface {
	ChainConfigProvider
	AmountConverter
	UTXOProvider
	RawTransactionCreator
	AddressOperator
	BalanceChecker
}

// BCHTransactionSigner combines interfaces needed for BCH transaction signing.
// Use this for BCH keygen/sign transaction usecases.
type BCHTransactionSigner interface {
	ChainConfigProvider
	RawTransactionCreator // ToHex, ToMsgTx
	RawTransactionSigner  // SignRawTransactionWithKey
}

// BTCTransactionSigner combines interfaces needed for BTC transaction signing.
// Use this for BTC keygen/sign transaction usecases.
type BTCTransactionSigner interface {
	ChainConfigProvider
	PSBTSigner
}
