// Package eth defines interfaces for Ethereum blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer. All types used in these interfaces are domain types,
// avoiding circular dependencies with the infrastructure layer.
//
// All layers MUST use the small, focused interfaces defined in this file.
package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/p2p"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// TxCreateParams contains the parameters returned from transaction creation
// that are needed by the use case layer to construct the domain entity.
// This DTO allows the use case to maintain full responsibility for domain entity creation.
type TxCreateParams struct {
	UUID        string // Transaction UUID for tracing
	FromAddress string // Sender's wallet address
	ToAddress   string // Receiver's wallet address
	Amount      uint64 // Final amount being sent (Wei)
	Fee         uint64 // Transaction fee (Wei)
	GasLimit    uint32 // Gas limit for transaction
	Nonce       uint64 // Transaction nonce

	// Transaction type and fee fields for JSON file construction
	EthTxType            uint8  // 0=legacy (LegacyTx), 2=EIP-1559 (DynamicFeeTx)
	ChainID              uint64 // EIP-155 chain identifier (EIP-1559 only)
	GasPrice             uint64 // Wei (legacy Type 0 only)
	MaxFeePerGas         uint64 // Wei (EIP-1559 Type 2 only)
	MaxPriorityFeePerGas uint64 // Wei (EIP-1559 Type 2 only)
}

// ERC20er defines the interface for ERC-20 token and native ETH transaction creation.
// Used by the Watch wallet create-transaction use case for both native ETH and ERC-20 tokens.
// Implementations handle token contract interactions, balance queries, and EIP-1559 support detection.
type ERC20er interface {
	ValidateAddr(addr string) error
	FloatToBigInt(v float64) *big.Int
	GetBalance(ctx context.Context, hexAddr string, quantityTag domainETH.QuantityTag) (*big.Int, error)
	CreateRawTransaction(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainETH.RawTx, *TxCreateParams, error)
	// SupportsEIP1559 detects whether the connected node supports EIP-1559 transactions.
	// Returns false for ERC-20 tokens (which use legacy transactions).
	SupportsEIP1559(ctx context.Context) bool
	// CreateRawTransactionEIP1559 creates an EIP-1559 (Type 2) transaction.
	// For ERC-20 tokens, this delegates to CreateRawTransaction (legacy fallback).
	CreateRawTransactionEIP1559(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainETH.RawTx, *TxCreateParams, error)
}

// EtherTxMonitor defines the interface for monitoring Ethereum transactions.
// Used by watch wallet monitor-transaction use case.
type EtherTxMonitor interface {
	GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainETH.UserAmount)
	GetConfirmation(ctx context.Context, hashTx string) (uint64, error)
	// GetTxReceipt retrieves a transaction receipt as a clean domain type.
	// Returns (nil, nil) when the transaction has not yet been included in a block.
	// Returns (nil, error) on node connectivity or parsing failures.
	GetTxReceipt(ctx context.Context, txHash string) (*domainETH.TransactionReceipt, error)
}

// =============================================================================
// Interfaces for specific use - Interface Segregation Principle (ISP)
// =============================================================================
//
// These interfaces are designed to be minimal and focused on specific use cases.
// Use cases and adapters MUST depend on these small interfaces instead of the
// large Ethereumer interface.
//
// DI Layer Integration:
// - DI layer continues to inject the full Ethereumer implementation
// - Go's implicit interface satisfaction handles the type conversion automatically
//
// =============================================================================

// ETHLifecycle manages the Ethereum client lifecycle.
// Used by wallet adapters (keygen, sign, watch) that need to close the connection
// and identify the coin type.
type ETHLifecycle interface {
	Close()
	CoinTypeCode() domainCoin.CoinTypeCode
}

// ETHKeyAccessor provides keystore access for private key import operations.
// Used by keygen import-private-key use case.
type ETHKeyAccessor interface {
	GetKeyDir() string
	ToECDSA(privKey string) (*ecdsa.PrivateKey, error)
}

// ETHTransactionSigner signs raw Ethereum transactions using the local keystore.
// Used by keygen and sign sign-transaction use cases.
type ETHTransactionSigner interface {
	SignOnRawTransaction(rawTx *domainETH.RawTx, passphrase string) (*domainETH.RawTx, error)
}

// ETHTransactionSender broadcasts signed Ethereum transactions to the network.
// Used by watch wallet send-transaction use case.
//
// Deprecated: Use TxSender instead (defined in the EIP-1559 Transaction Flow
// Interfaces section below). ETHTransactionSender is retained only for backward
// compatibility; new code should depend on TxSender.
type ETHTransactionSender = TxSender

// ETHRawKeyImporter imports raw private keys via the Ethereum RPC.
// Used by keygen API CLI commands.
type ETHRawKeyImporter interface {
	ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error)
}

// ETHNodeAPIClient provides Ethereum node API operations for watch wallet CLI commands.
type ETHNodeAPIClient interface {
	ClientVersion(ctx context.Context) (string, error)
	NetVersion(ctx context.Context) (uint16, error)
	NodeInfo(ctx context.Context) (*p2p.NodeInfo, error)
	Syncing(ctx context.Context) (*domainETH.ResponseSyncing, bool, error)
}

// =============================================================================
// Composed Interfaces for Wallet Adapters
// =============================================================================

// ETHKeygenSignClient combines lifecycle and raw key import for keygen/sign wallet adapters.
// Used by ETHKeygen and ETHSign wallet adapters, which need both lifecycle management
// and the ability to expose the importrawkey CLI command.
type ETHKeygenSignClient interface {
	ETHLifecycle
	ETHRawKeyImporter
}

// ETHWatchClient combines lifecycle and node API operations for the watch wallet adapter.
// Used by ETHWatch wallet adapter, which needs lifecycle management and watch node CLI commands.
type ETHWatchClient interface {
	ETHLifecycle
	ETHNodeAPIClient
}

// =============================================================================
// EIP-1559 Transaction Flow Interfaces (Task 2.1b)
// =============================================================================
//
// These interfaces define the port boundary for the EIP-1559 transaction flow.
// They use clean domain types to prevent infrastructure type leakage into the
// use case layer.
//
// Note: Some method signatures differ from the monolithic Ethereumer to use
// domain-friendly types (e.g., EstimateGas uses plain parameters instead of
// ethereum.CallMsg). The infrastructure implementations will be updated in
// Tasks 4.x to satisfy these interfaces.
// =============================================================================

// ChainConfigProvider provides chain-level configuration.
// Used by keygen signing and watch transaction creation use cases.
type ChainConfigProvider interface {
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// BalanceChecker retrieves account balances.
// Used by watch wallet monitor-transaction use case.
type BalanceChecker interface {
	GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, error)
	BalanceAt(ctx context.Context, addr string) (*big.Int, error)
}

// TxCreator creates unsigned transactions for both legacy and EIP-1559 formats.
// Used by watch wallet create-transaction use case.
type TxCreator interface {
	CreateRawTransaction(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainETH.RawTx, *TxCreateParams, error)
	CreateRawTransactionEIP1559(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainETH.RawTx, *TxCreateParams, error)
	SupportsEIP1559(ctx context.Context) bool
}

// GasEstimator estimates gas and fees for transaction creation.
// Used by watch wallet create-transaction use case.
// Note: EstimateGas uses domain-friendly parameters instead of ethereum.CallMsg.
type GasEstimator interface {
	GasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, from, to string, value *big.Int) (uint64, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

// TxSigner signs raw transactions offline using a private key directly.
// Used by keygen wallet sign-transaction use case for air-gapped offline signing.
// This differs from ETHTransactionSigner which uses a keystore passphrase:
// SignTxWithPrivateKey accepts *ecdsa.PrivateKey directly, enabling true offline
// operation without requiring an Ethereum node or keystore.
type TxSigner interface {
	SignTxWithPrivateKey(
		rawTx *domainETH.RawTx, privKey *ecdsa.PrivateKey, chainID *big.Int,
	) (*domainETH.RawTx, error)
}

// TxSender broadcasts signed transactions to the network.
// Used by watch wallet send-transaction use case.
// ETHTransactionSender is a type alias for this interface for backward compatibility.
type TxSender interface {
	SendSignedRawTransaction(ctx context.Context, signedTxHex string) (string, error)
}

// TxMonitor retrieves transaction status and confirmation count.
// Used by watch wallet monitor-transaction use case.
type TxMonitor interface {
	GetTransactionReceipt(ctx context.Context, txHash string) (*domainETH.TransactionReceipt, error)
	GetConfirmation(ctx context.Context, txHash string) (uint64, error)
}

// AddressValidator validates Ethereum addresses.
// Used by watch wallet create-transaction use case.
type AddressValidator interface {
	ValidateAddr(addr string) error
}

// =============================================================================
// Composed Interfaces for EIP-1559 Transaction Flow Use Cases
// =============================================================================

// WatchTxCreationDeps is the composed interface for the Watch wallet's
// create-transaction use case. It combines all dependencies needed to
// create both legacy and EIP-1559 transactions.
type WatchTxCreationDeps interface {
	ChainConfigProvider
	TxCreator
	GasEstimator
	AddressValidator
}

// KeygenSignTxDeps is the composed interface for the Keygen wallet's
// sign-transaction use case. ETH is single-sig EOA; there is no Sign wallet.
// This combines chain configuration with offline private-key signing capability.
type KeygenSignTxDeps interface {
	ChainConfigProvider
	TxSigner
}
