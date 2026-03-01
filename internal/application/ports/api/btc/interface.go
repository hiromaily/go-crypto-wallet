package btc

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/quagmt/udecimal"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	btcdescriptor "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/descriptor"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// Bitcoiner is the full Bitcoin/BitcoinCash interface.
//
// Usage restriction: Bitcoiner MUST only be referenced in the DI layer (internal/di/).
// All other layers (use cases, adapters, CLI) MUST use the small, focused interfaces
// defined below to satisfy the Interface Segregation Principle.
//
// Use the focused interfaces instead:
//
//   - ChainConfigProvider: Chain configuration and coin type
//   - AmountConverter: Amount format conversion
//   - UTXOProvider: UTXO queries
//   - RawTransactionCreator: Raw transaction creation
//   - AddressOperator: Address operations
//   - BalanceChecker: Balance queries
//   - TransactionSender: Transaction broadcast
//   - TransactionMonitor: Transaction monitoring
//   - AddressImporter: Address import
//   - PrivateKeyImporter: Private key import
//   - DescriptorManager: Descriptor operations (BTC only)
//   - PSBTCreator: Create PSBT (BTC only)
//   - PSBTSigner: Sign PSBT (BTC only)
//   - PSBTFinalizer: Finalize PSBT (BTC only)
//   - RawTransactionSigner: Raw transaction signing
//   - MultisigManager: Multisig address management
type Bitcoiner interface {
	// public_account.go -> wrapper of GetAddressInfo to return account
	GetAccount(addr string) (string, error)

	// address.go
	GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error)
	GetAddressesByLabel(labelName string) ([]btcutil.Address, error)
	ValidateAddress(addr string) (*btcrpc.ValidateAddressResult, error)
	DecodeAddress(addr string) (btcutil.Address, error)

	// amount.go
	AmountString(amt btcutil.Amount) string
	AmountToDecimal(amt btcutil.Amount) (udecimal.Decimal, error)
	FloatToDecimal(f float64) (udecimal.Decimal, error)
	FloatToAmount(f float64) (btcutil.Amount, error)
	StrToAmount(s string) (btcutil.Amount, error)
	StrSatoshiToAmount(s string) (btcutil.Amount, error)

	// balance.go
	GetBalance() (btcutil.Amount, error)
	GetBalanceByListUnspent(confirmationNum uint64) (btcutil.Amount, error)
	GetBalanceByAccount(accountType domainAccount.AccountType, confirmationNum uint64) (btcutil.Amount, error)

	// block.go
	GetBlockCount() (int64, error)

	// bitcoin.go
	Close()
	GetChainConf() *chaincfg.Params
	SetChainConf(conf *chaincfg.Params)
	SetChainConfNet(btcNet wire.BitcoinNet)
	ConfirmationBlock() uint64
	FeeRangeMax() float64
	FeeRangeMin() float64
	Version() domainBTC.Version
	CoinTypeCode() domainCoin.CoinTypeCode
	GetPkgRPC() btcrpc.BTCRPC

	// fee.go
	GetTransactionFee(tx *wire.MsgTx) (btcutil.Amount, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)

	// import.go
	ImportPrivKey(privKeyWIF *btcutil.WIF) error
	ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error
	ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error
	ImportAddress(pubkey string) error
	ImportAddressWithoutReScan(pubkey string) error
	ImportAddressWithLabel(addr, label string, rescan bool) error
	ImportDescriptors(requests []btcrpc.ImportDescriptorsRequest) ([]btcrpc.ImportDescriptorsResponse, error)
	ImportMulti(
		requests []btcrpc.ImportMultiRequest,
		options *btcrpc.ImportMultiOptions,
	) ([]btcrpc.ImportMultiResponse, error)

	// descriptor_info.go
	GetDescriptorInfo(descriptor string) (*btcrpc.DescriptorInfo, error)
	ListDescriptors(privateDescriptors bool) (*btcrpc.ListDescriptorsResult, error)

	// label.go
	SetLabel(addr, label string) error
	// GetReceivedByLabelAndMinConf(accountName string, minConf int) (btcutil.Amount, error)

	// logging.go
	Logging() (*btcrpc.LoggingResult, error)

	// multisig.go
	AddMultisigAddress(
		requiredSigs int, addresses []string, accountName string, addressType domainBTC.AddressType,
	) (*btcrpc.AddMultisigAddressResult, error)

	// network.go
	GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error)
	GetBlockchainInfo() (*btcrpc.GetBlockchainInfoResult, error)

	// transaction.go
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
	GetTransactionByTxID(txID string) (*btcrpc.GetTransactionResult, error)
	GetTxOutByTxID(txID string, index uint32) (*btcjson.GetTxOutResult, error)
	DecodeRawTransaction(hexTx string) (*btcrpc.TxRawResult, error)
	GetRawTransactionByHex(strHashTx string) (*btcutil.Tx, error)
	CreateRawTransaction(
		inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
	) (*wire.MsgTx, error)
	FundRawTransaction(hex string) (*btcrpc.FundRawTransactionResult, error)
	SignRawTransaction(tx *wire.MsgTx, prevtxs []dtobtc.PreviousTx) (*wire.MsgTx, bool, error)
	SignRawTransactionWithKey(
		tx *wire.MsgTx, privKeysWIF []string, prevtxs []dtobtc.PreviousTx,
	) (*wire.MsgTx, bool, error)
	SendTransactionByHex(hex string) (*chainhash.Hash, error)
	SendTransactionByByte(rawTx []byte) (*chainhash.Hash, error)
	Sign(tx *wire.MsgTx, strPrivateKey string) (string, error)

	// psbt.go (BIP174 Partially Signed Bitcoin Transaction support)
	CreatePSBT(msgTx *wire.MsgTx, prevTxs []dtobtc.PreviousTx, senderAccount domainAccount.AccountType) (string, error)
	ParsePSBT(psbtBase64 string) (*dtobtc.ParsedPSBT, error)
	ValidatePSBT(psbtBase64 string) error
	SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)
	WalletProcessPsbt(psbtBase64 string, sign bool) (string, bool, error) // RPC-based signing for descriptor wallets
	FinalizePSBT(psbtBase64 string) (string, error)
	ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
	IsPSBTComplete(psbtBase64 string) (bool, error)
	GetPSBTFee(psbtBase64 string) (int64, error)

	// unspent.go
	ListUnspent(confirmationNum uint64) ([]dtobtc.UnspentOutput, error)
	ListUnspentByAccount(
		accountType domainAccount.AccountType, confirmationNum uint64,
	) ([]dtobtc.UnspentOutput, error)
	GetUnspentListAddrs(
		unspentList []dtobtc.UnspentOutput, accountType domainAccount.AccountType,
	) []string
	LockUnspent(tx *dtobtc.UnspentOutput) error
	UnlockUnspent() error
}

// =============================================================================
// Interfaces for specific - Interface Segregation Principle (ISP)
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
//
// =============================================================================
// Usage Pattern
// =============================================================================
//
// Each usecase defines a local interface by embedding the small interfaces it needs:
//
//	// In usecase file (e.g., watch/btc/create_transaction.go)
//	type createTxBTCClient interface {
//	    apibtc.ChainConfigProvider
//	    apibtc.AmountConverter
//	    apibtc.UTXOProvider
//	    apibtc.RawTransactionCreator
//	    apibtc.AddressOperator
//	    apibtc.BalanceChecker
//	    apibtc.PSBTCreator  // BTC only
//	}
//
//	type createTransactionUseCase struct {
//	    btcClient createTxBTCClient  // Use local interface, not Bitcoiner
//	    // ...
//	}
//
// DI Layer Integration:
// - DI layer continues to inject the full Bitcoiner implementation
// - Go's implicit interface satisfaction handles the type conversion automatically
// - No changes needed in DI layer when usecases use these small interfaces
//
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

// RawTransactionConverter converts between transaction formats.
// Used by signing usecases that need to convert hex <-> MsgTx.
type RawTransactionConverter interface {
	ToHex(tx *wire.MsgTx) (string, error)
	ToMsgTx(txHex string) (*wire.MsgTx, error)
}

// RawTransactionCreator creates and manipulates raw transactions.
// Used by watch wallet transaction creation.
type RawTransactionCreator interface {
	RawTransactionConverter
	CreateRawTransaction(
		inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount,
	) (*wire.MsgTx, error)
	GetFee(tx *wire.MsgTx, adjustmentFee float64) (btcutil.Amount, error)
}

// AddressOperator provides address-related operations.
// Used by transaction creation and import usecases.
type AddressOperator interface {
	DecodeAddress(addr string) (btcutil.Address, error)
	GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error)
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
	GetTransactionByTxID(txID string) (*btcrpc.GetTransactionResult, error)
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
		requests []btcrpc.ImportMultiRequest, options *btcrpc.ImportMultiOptions,
	) ([]btcrpc.ImportMultiResponse, error)
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
	GetDescriptorInfo(descriptor string) (*btcrpc.DescriptorInfo, error)
	ImportDescriptors(requests []btcrpc.ImportDescriptorsRequest) ([]btcrpc.ImportDescriptorsResponse, error)
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
		requiredSigs int, addresses []string, accountName string, addressType domainBTC.AddressType,
	) (*btcrpc.AddMultisigAddressResult, error)
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
	RawTransactionConverter // ToHex, ToMsgTx
	RawTransactionSigner    // SignRawTransactionWithKey
}

// BTCTransactionSigner combines interfaces needed for BTC transaction signing.
// Use this for BTC keygen/sign transaction usecases.
type BTCTransactionSigner interface {
	ChainConfigProvider
	PSBTSigner
	PSBTHandler
}

// PrivateKeyImportClient combines interfaces needed for private key import.
// Use this for BTC/BCH keygen/sign import-private-key usecases.
type PrivateKeyImportClient interface {
	ChainConfigProvider // CoinTypeCode
	PrivateKeyImporter  // ImportPrivKeyWithoutReScan
	AddressOperator     // GetAccount, GetAddressInfo
}

// AddressImportClient combines interfaces needed for address import.
// Use this for BTC/BCH watch import-address usecases.
type AddressImportClient interface {
	ChainConfigProvider   // CoinTypeCode
	LegacyAddressImporter // ImportAddressWithLabel, ImportMulti
	AddressOperator       // GetAddressInfo
}

// -----------------------------------------------------------------------------
// Additional Small Interfaces (for CLI and lifecycle management)
// -----------------------------------------------------------------------------

// NodeBalanceChecker checks total wallet balance across all accounts.
type NodeBalanceChecker interface {
	GetBalance() (btcutil.Amount, error)
}

// PKGRPCProvider provides direct access to the underlying RPC client.
// Used by callers that need raw RPC operations (fee estimation, wallet security, etc.)
// without going through delegation wrappers.
type PKGRPCProvider interface {
	GetPkgRPC() btcrpc.BTCRPC
}

// NetworkInformer provides network and blockchain information.
type NetworkInformer interface {
	GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error)
	GetBlockchainInfo() (*btcrpc.GetBlockchainInfoResult, error)
}

// FullUTXOLister lists all unspent outputs without account filter.
type FullUTXOLister interface {
	ListUnspent(confirmationNum uint64) ([]dtobtc.UnspentOutput, error)
}

// NodeLogger provides node-level logging operations.
type NodeLogger interface {
	Logging() (*btcrpc.LoggingResult, error)
}

// AddressValidator validates Bitcoin addresses.
type AddressValidator interface {
	ValidateAddress(addr string) (*btcrpc.ValidateAddressResult, error)
}

// UTXOLocker locks and unlocks UTXOs.
type UTXOLocker interface {
	LockUnspent(tx *dtobtc.UnspentOutput) error
	UnlockUnspent() error
}

// BTCLifecycle manages the Bitcoin client lifecycle.
type BTCLifecycle interface {
	Close()
}

// -----------------------------------------------------------------------------
// Composed Interfaces for Watch Wallet API CLI
// -----------------------------------------------------------------------------

// WatchAPIClient combines interfaces needed for the watch wallet API CLI commands.
// Used by the btcWatchAPICmds adapter and btcWatchClient wallet interface.
type WatchAPIClient interface {
	NodeBalanceChecker // GetBalance
	BalanceChecker     // GetBalanceByAccount
	PKGRPCProvider     // GetPkgRPC
	FullUTXOLister     // ListUnspent
	UTXOProvider       // ListUnspentByAccount, GetUnspentListAddrs
	NodeLogger         // Logging
	AddressValidator   // ValidateAddress
	// Inline: partial overlaps with larger interfaces to avoid importing extra methods
	ConfirmationBlock() uint64
	GetAddressInfo(addr string) (*btcrpc.GetAddressInfoResult, error)
	GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error)
	UnlockUnspent() error
}

// -----------------------------------------------------------------------------
// MuSig2, Descriptor, and HD Key Interfaces
// -----------------------------------------------------------------------------

// MuSig2Servicer provides MuSig2 multi-signature functionality.
// Used by keygen/sign musig2 use cases.
type MuSig2Servicer interface {
	CreateContext(privKey *btcec.PrivateKey, allPubKeys []*btcec.PublicKey, sortKeys bool) (*musig2.Context, error)
	CreateContextWithTaproot(
		privKey *btcec.PrivateKey, allPubKeys []*btcec.PublicKey, sortKeys bool,
	) (*musig2.Context, error)
	CreateSession(ctx *musig2.Context) (*musig2.Session, error)
	GetPublicNonce(session *musig2.Session) [66]byte
	RegisterPubNonce(session *musig2.Session, pubNonce [66]byte) (bool, error)
	Sign(session *musig2.Session, messageHash [32]byte) (*musig2.PartialSignature, error)
	CombineSig(session *musig2.Session, partialSig *musig2.PartialSignature) (bool, error)
	FinalSig(session *musig2.Session) (*schnorr.Signature, error)
	GetCombinedKey(ctx *musig2.Context) (*btcec.PublicKey, error)
	VerifySignature(signature *schnorr.Signature, messageHash [32]byte, aggregatedKey *btcec.PublicKey) bool
}

// DescriptorParserr parses descriptor strings into domain descriptor objects.
// Note: double 'r' to avoid conflict with descriptorParser local variable.
type DescriptorParserr interface {
	Parse(descriptorStr string) (*btcdescriptor.Descriptor, error)
}

// DescriptorServicer generates output descriptors for various address types.
// Used by keygen generate-descriptor use case.
type DescriptorServicer interface {
	GenerateTaprootDescriptor(
		fingerprint, derivationPath string, xpub *hdkeychain.ExtendedKey, isChange bool,
	) (string, error)
	GenerateBech32Descriptor(
		fingerprint, derivationPath string, xpub *hdkeychain.ExtendedKey, isChange bool,
	) (string, error)
	GenerateP2SHSegWitDescriptor(
		fingerprint, derivationPath string, xpub *hdkeychain.ExtendedKey, isChange bool,
	) (string, error)
	GenerateP2PKHDescriptor(
		fingerprint, derivationPath string, xpub *hdkeychain.ExtendedKey, isChange bool,
	) (string, error)
	GenerateTaprootScriptPathDescriptor(signers []btcdescriptor.MultisigSigner, isChange bool) (string, error)
	GenerateMultisigDescriptor(
		requiredSigs int,
		signers []btcdescriptor.MultisigSigner,
		isChange bool,
		descriptorType btcdescriptor.DescriptorType,
	) (string, error)
}

// HDKeyOperator derives HD keys and fingerprints from seeds.
// Used by keygen generate-descriptor and sign export-fullpubkey use cases.
type HDKeyOperator interface {
	GetMasterFingerprintHex(seed []byte) (string, error)
	GetCoinLevelXPub(seed []byte, purpose uint8) (string, error)
	GetAccountXPub(seed []byte, purpose uint8, accountType domainAccount.AccountType) (string, error)
}
