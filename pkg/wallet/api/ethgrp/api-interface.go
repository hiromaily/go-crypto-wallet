package ethgrp

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Ethereumer Ethereum Interface
type Ethereumer interface {
	// balance
	GetTotalBalance(addrs []string) (*big.Int, []eth.UserAmount)
	// client
	BalanceAt(hexAddr string) (*big.Int, error)
	SendRawTx(tx *types.Transaction) error
	// ethereum
	Close()
	CoinTypeCode() coin.CoinTypeCode
	GetChainConf() *chaincfg.Params
	// key
	ToECDSA(privKey string) (*ecdsa.PrivateKey, error)
	GetKeyDir(accountType account.AccountType) string
	GetPrivKey(hexAddr, password string, accountType account.AccountType) (*keystore.Key, error)
	RenameParityKeyFile(hexAddr string, accountType account.AccountType) error
	// rpc_admin
	AddPeer(nodeURL string) error
	AdminDataDir() (string, error)
	NodeInfo() (*p2p.NodeInfo, error)
	AdminPeers() ([]*p2p.PeerInfo, error)
	// rpc_eth
	Syncing() (*eth.ResponseSyncing, bool, error)
	ProtocolVersion() (uint64, error)
	Coinbase() (string, error)
	Accounts() ([]string, error)
	BlockNumber() (*big.Int, error)
	EnsureBlockNumber(loopCount int) (*big.Int, error)
	GetBalance(hexAddr string, quantityTag eth.QuantityTag) (*big.Int, error)
	// GetStoreageAt(hexAddr string, quantityTag eth.QuantityTag) (string, error)
	GetTransactionCount(hexAddr string, quantityTag eth.QuantityTag) (*big.Int, error)
	// GetBlockTransactionCountByBlockHash(blockHash string) (*big.Int, error)
	GetBlockTransactionCountByNumber(blockNumber uint64) (*big.Int, error)
	// GetUncleCountByBlockHash(blockHash string) (*big.Int, error)
	GetUncleCountByBlockNumber(blockNumber uint64) (*big.Int, error)
	// GetCode(hexAddr string, quantityTag eth.QuantityTag) (*big.Int, error)
	GetBlockByNumber(blockNumber uint64) (*eth.BlockInfo, error)
	// rpc_eth_gas
	GasPrice() (*big.Int, error)
	EstimateGas(msg *ethereum.CallMsg) (*big.Int, error)
	// rpc_eth_tx
	Sign(hexAddr, message string) (string, error)
	SendTransaction(msg *ethereum.CallMsg) (string, error)
	SendRawTransaction(signedTx string) (string, error)
	SendRawTransactionWithTypesTx(tx *types.Transaction) (string, error)
	GetTransactionByHash(hashTx string) (*eth.ResponseGetTransaction, error)
	GetTransactionReceipt(hashTx string) (*eth.ResponseGetTransactionReceipt, error)
	// rpc_miner
	StartMining() error
	StopMining() error
	Mining() (bool, error)
	HashRate() (*big.Int, error)
	// rpc_net
	NetVersion() (uint16, error)
	NetListening() (bool, error)
	NetPeerCount() (*big.Int, error)
	// rpc_personal
	ImportRawKey(hexKey, passPhrase string) (string, error)
	ListAccounts() ([]string, error)
	NewAccount(passphrase string, accountType account.AccountType) (string, error)
	LockAccount(hexAddr string) error
	UnlockAccount(hexAddr, passphrase string, duration uint64) (bool, error)
	// rpc_web3
	ClientVersion() (string, error)
	SHA3(data string) (string, error)
	// transaction
	CreateRawTransaction(fromAddr, toAddr string, amount uint64, additionalNonce int) (*eth.RawTx, *models.EthDetailTX, error)
	SignOnRawTransaction(rawTx *eth.RawTx, passphrase string, accountType account.AccountType) (*eth.RawTx, error)
	SendSignedRawTransaction(signedTxHex string) (string, error)
	GetConfirmation(hashTx string) (uint64, error)
	// util
	DecodeBig(input string) (*big.Int, error)
	ValidationAddr(addr string) error
	FromWei(v int64) *big.Int
	FromGWei(v int64) *big.Int
	// FromEther(v int64) *big.Int
	FromFloatEther(v float64) *big.Int
}
