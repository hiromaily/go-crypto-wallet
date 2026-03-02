package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
)

// ETHRPC defines all Ethereum JSON-RPC operations exposed by rpcClient.
type ETHRPC interface {
	// admin.go
	AddPeer(ctx context.Context, nodeURL string) error
	AdminDataDir(ctx context.Context) (string, error)
	NodeInfo(ctx context.Context) (*p2p.NodeInfo, error)
	AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error)

	// eth.go
	Syncing(ctx context.Context) (*ResponseSyncing, bool, error)
	ProtocolVersion(ctx context.Context) (uint64, error)
	Coinbase(ctx context.Context) (string, error)
	Accounts(ctx context.Context) ([]string, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	EnsureBlockNumber(ctx context.Context, loopCount int) (*big.Int, error)
	GetBalance(ctx context.Context, hexAddr string, quantityTag QuantityTag) (*big.Int, error)
	GetTransactionCount(ctx context.Context, hexAddr string, quantityTag QuantityTag) (*big.Int, error)
	GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	GetUncleCountByBlockNumber(ctx context.Context, blockNumber uint64) (*big.Int, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*BlockInfo, error)

	// gas.go
	GasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (*big.Int, error)

	// miner.go
	StartMining(ctx context.Context) error
	StopMining(ctx context.Context) error
	Mining(ctx context.Context) (bool, error)
	HashRate(ctx context.Context) (*big.Int, error)

	// net.go
	NetVersion(ctx context.Context) (uint16, error)
	NetListening(ctx context.Context) (bool, error)
	NetPeerCount(ctx context.Context) (*big.Int, error)

	// personal.go
	ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error)
	ListAccounts(ctx context.Context) ([]string, error)
	NewAccount(ctx context.Context, passphrase string) (string, error)
	LockAccount(ctx context.Context, hexAddr string) error
	UnlockAccount(ctx context.Context, hexAddr, passphrase string, duration uint64) (bool, error)

	// transaction.go
	Sign(ctx context.Context, hexAddr, message string) (string, error)
	SendTransaction(ctx context.Context, msg *ethereum.CallMsg) (string, error)
	SendRawTransaction(ctx context.Context, signedTx string) (string, error)
	SendRawTransactionWithTypesTx(ctx context.Context, tx *types.Transaction) (string, error)
	GetTransactionByHash(ctx context.Context, hashTx string) (*ResponseGetTransaction, error)
	GetTransactionReceipt(ctx context.Context, hashTx string) (*ResponseGetTransactionReceipt, error)
	GetTxReceipt(ctx context.Context, txHash string) (*TransactionReceipt, error)

	// web3.go
	ClientVersion(ctx context.Context) (string, error)
	SHA3(ctx context.Context, data string) (string, error)

	// rpc.go
	Close()
}
