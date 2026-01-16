// Package ethereum provides Ethereum blockchain API implementations.
//
// Deprecated: The interfaces in this package have been moved to
// internal/application/ports/ethereum following Clean Architecture principles.
// Import from there instead:
//
//	import apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
//
// The implementations remain in this package but now implement the interfaces
// defined in the application/ports layer.
package ethereum

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
)

// Ethereumer defines the interface for Ethereum blockchain operations.
//
// Deprecated: This interface has been moved to internal/application/ports/ethereum.
// Use apieth.Ethereumer instead.
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
	// GetStoreageAt(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (string, error)
	GetTransactionCount(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
	// GetBlockTransactionCountByBlockHash(ctx context.Context, blockHash string) (*big.Int, error)
	GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	// GetUncleCountByBlockHash(ctx context.Context, blockHash string) (*big.Int, error)
	GetUncleCountByBlockNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	// GetCode(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
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
	) (*domainEthereum.RawTx, *apieth.TxCreateParams, error)
	SignOnRawTransaction(rawTx *domainEthereum.RawTx, passphrase string) (*domainEthereum.RawTx, error)
	SendSignedRawTransaction(ctx context.Context, signedTxHex string) (string, error)
	GetConfirmation(ctx context.Context, hashTx string) (uint64, error)
	// util
	DecodeBig(input string) (*big.Int, error)
	ValidateAddr(addr string) error
	FromWei(v int64) *big.Int
	FromGWei(v int64) *big.Int
	// FromEther(v int64) *big.Int
	FromFloatEther(v float64) *big.Int
	FloatToBigInt(v float64) *big.Int
}

// ERC20er defines the interface for ERC-20 token operations.
//
// Deprecated: This interface has been moved to internal/application/ports/ethereum.
// Use apieth.ERC20er instead.
type ERC20er interface {
	ValidateAddr(addr string) error
	FloatToBigInt(v float64) *big.Int
	GetBalance(ctx context.Context, hexAddr string, quantityTag domainEthereum.QuantityTag) (*big.Int, error)
	CreateRawTransaction(
		ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
	) (*domainEthereum.RawTx, *apieth.TxCreateParams, error)
}

// EtherTxCreator is a type alias for ERC20er used in transaction creation contexts.
//
// Deprecated: This type alias has been moved to internal/application/ports/ethereum.
// Use apieth.EtherTxCreator instead.
type EtherTxCreator = ERC20er

// EtherTxMonitor defines the interface for monitoring Ethereum transactions.
//
// Deprecated: This interface has been moved to internal/application/ports/ethereum.
// Use apieth.EtherTxMonitor instead.
type EtherTxMonitor interface {
	GetTotalBalance(ctx context.Context, addrs []string) (*big.Int, []domainEthereum.UserAmount)
	GetConfirmation(ctx context.Context, hashTx string) (uint64, error)
}
