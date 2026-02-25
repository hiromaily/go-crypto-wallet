// Package eth defines interfaces for Ethereum blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer. All types used in these interfaces are domain types,
// avoiding circular dependencies with the infrastructure layer.
//
// Usage restriction: Ethereumer MUST only be referenced in the DI layer (internal/di/).
// All other layers MUST use the small, focused interfaces defined in this file.
package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
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
}

// Ethereumer is the full Ethereum interface.
//
// Usage restriction: Ethereumer MUST only be referenced in the DI layer (internal/di/).
// All other layers (use cases, adapters, CLI) MUST use the small, focused interfaces
// defined below to satisfy the Interface Segregation Principle.
//
// Use the focused interfaces instead:
//
//   - ETHLifecycle: Client lifecycle (Close, CoinTypeCode)
//   - ETHKeyAccessor: Keystore access for key import
//   - ETHTransactionSigner: Sign raw transactions
//   - ETHTransactionSender: Broadcast signed transactions
//   - ETHRawKeyImporter: Import raw keys via RPC (deprecated)
//   - ETHNodeAPIClient: Node API operations for CLI
//   - ERC20er: ERC-20 and ETH transaction creation
//   - EtherTxMonitor: Transaction monitoring
type Ethereumer interface {
	// balance
	GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainEthereum.UserAmount)
	// client
	BalanceAt(ctx context.Context, hexAddr string) (*big.Int, error)
	SendRawTx(ctx context.Context, tx *types.Transaction) error
	// ethereum
	Close()
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
	// key
	ToECDSA(privKey string) (*ecdsa.PrivateKey, error)
	GetKeyDir() string
	GetPrivKey(hexAddr, password string) (*keystore.Key, error)
	// rpc_admin
	AddPeer(ctx context.Context, nodeURL string) error
	AdminDataDir(ctx context.Context) (string, error)
	NodeInfo(ctx context.Context) (*p2p.NodeInfo, error)
	AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error)
	// rpc_eth
	Syncing(ctx context.Context) (*domainEthereum.ResponseSyncing, bool, error)
	ProtocolVersion(ctx context.Context) (uint64, error)
	Coinbase(ctx context.Context) (string, error)
	Accounts(ctx context.Context) ([]string, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	EnsureBlockNumber(ctx context.Context, loopCount int) (*big.Int, error)
	GetBalance(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
	GetTransactionCount(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
	GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	GetUncleCountByBlockNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*domainEthereum.BlockInfo, error)
	// rpc_eth_gas
	GasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (*big.Int, error)
	// rpc_eth_tx
	Sign(ctx context.Context, hexAddr, message string) (string, error)
	SendTransaction(ctx context.Context, msg *ethereum.CallMsg) (string, error)
	SendRawTransaction(ctx context.Context, signedTx string) (string, error)
	SendRawTransactionWithTypesTx(ctx context.Context, tx *types.Transaction) (string, error)
	GetTransactionByHash(ctx context.Context, hashTx string) (*domainEthereum.ResponseGetTransaction, error)
	GetTransactionReceipt(ctx context.Context, hashTx string) (*domainEthereum.ResponseGetTransactionReceipt, error)
	// rpc_miner
	StartMining(ctx context.Context) error
	StopMining(ctx context.Context) error
	Mining(ctx context.Context) (bool, error)
	HashRate(ctx context.Context) (*big.Int, error)
	// rpc_net
	NetVersion(ctx context.Context) (uint16, error)
	NetListening(ctx context.Context) (bool, error)
	NetPeerCount(ctx context.Context) (*big.Int, error)
	// rpc_personal
	ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error)
	ListAccounts(ctx context.Context) ([]string, error)
	NewAccount(ctx context.Context, passphrase string, accountType domainAccount.AccountType) (string, error)
	LockAccount(ctx context.Context, hexAddr string) error
	UnlockAccount(ctx context.Context, hexAddr, passphrase string, duration uint64) (bool, error)
	// rpc_web3
	ClientVersion(ctx context.Context) (string, error)
	SHA3(ctx context.Context, data string) (string, error)
	// transaction
	CreateRawTransaction(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainEthereum.RawTx, *TxCreateParams, error)
	CreateRawTransactionEIP1559(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainEthereum.RawTx, *TxCreateParams, error)
	SupportsEIP1559(ctx context.Context) bool
	SignOnRawTransaction(rawTx *domainEthereum.RawTx, passphrase string) (*domainEthereum.RawTx, error)
	SendSignedRawTransaction(ctx context.Context, signedTxHex string) (string, error)
	GetConfirmation(ctx context.Context, hashTx string) (uint64, error)
	// util
	DecodeBig(input string) (*big.Int, error)
	ValidateAddr(addr string) error
	FromWei(v int64) *big.Int
	FromGWei(v int64) *big.Int
	FromFloatEther(v float64) *big.Int
	FloatToBigInt(v float64) *big.Int
}

// ERC20er defines the interface for ERC-20 token operations.
// Implementations handle token contract interactions and balance queries.
type ERC20er interface {
	ValidateAddr(addr string) error
	FloatToBigInt(v float64) *big.Int
	GetBalance(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
	CreateRawTransaction(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainEthereum.RawTx, *TxCreateParams, error)
}

// EtherTxMonitor defines the interface for monitoring Ethereum transactions.
// Used by watch wallet monitor-transaction use case.
type EtherTxMonitor interface {
	GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainEthereum.UserAmount)
	GetConfirmation(ctx context.Context, hashTx string) (uint64, error)
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
	SignOnRawTransaction(rawTx *domainEthereum.RawTx, passphrase string) (*domainEthereum.RawTx, error)
}

// ETHTransactionSender broadcasts signed Ethereum transactions to the network.
// Used by watch wallet send-transaction use case.
type ETHTransactionSender interface {
	SendSignedRawTransaction(ctx context.Context, signedTxHex string) (string, error)
}

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
	Syncing(ctx context.Context) (*domainEthereum.ResponseSyncing, bool, error)
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
