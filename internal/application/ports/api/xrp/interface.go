// Package xrp defines interfaces for XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
//
// Small, focused interfaces follow the Interface Segregation Principle.
// Use cases should depend only on the specific interfaces they need through local interface
// compositions.
package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
)

// XRPer defines the main interface for XRP blockchain operations.
// It embeds specialized interfaces for admin, public, and WebSocket transaction operations.
type XRPer interface {
	XRPAdminer
	XRPPublicer
	TransactionSubmitter
	TransactionSigner

	// account
	GetAccountInfo(ctx context.Context, address string) (*xrplgo.AccountInfo, error)

	// balance
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64

	// transaction
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)

	// xrp
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// SignerEntryInput represents a signer for creating SignerListSet transactions
type SignerEntryInput struct {
	Account string
	Weight  uint32
}

// XRPPublicer defines the interface for XRP public node operations.
// These operations query public information from the XRP network.
type XRPPublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*xrprpc.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*xrprpc.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*xrprpc.ResponseServerInfo, error)
}

// XRPAdminer defines the interface for XRP admin node operations.
// These operations typically require admin access to the XRP node.
type XRPAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*xrprpc.ResponseValidationCreate, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType dtoxrp.XRPKeyType,
	) (*xrprpc.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*xrprpc.ResponseWalletPropose, error)
}

// Small, focused interfaces following the Interface Segregation Principle.

// CoinTypeProvider provides coin type information.
type CoinTypeProvider interface {
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// BalanceChecker provides balance query operations.
type BalanceChecker interface {
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64
}

// TransactionPreparer prepares raw transactions.
type TransactionPreparer interface {
	CreateRawTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.TxInput, string, error)
}

// TransactionCombiner combines multiple signed transactions (for multisig).
type TransactionCombiner interface {
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
}

// RegularKeyPreparer prepares SetRegularKey transactions.
type RegularKeyPreparer interface {
	PrepareSetRegularKeyTransaction(
		ctx context.Context,
		senderAccount, regularKey string,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SetRegularKeyTxInput, string, error)
}

// SignerListPreparer prepares SignerListSet transactions.
type SignerListPreparer interface {
	PrepareSignerListSetTransaction(
		ctx context.Context,
		senderAccount string,
		signerQuorum uint32,
		signerEntries []SignerEntryInput,
		instructions *dtoxrp.Instructions,
	) (*dtoxrp.SignerListSetTxInput, string, error)
}

// KeyGenerator generates XRP keys/wallets.
type KeyGenerator interface {
	WalletPropose(ctx context.Context, passphrase string) (*xrprpc.ResponseWalletPropose, error)
}

// Closer provides cleanup operations.
type Closer interface {
	Close() error
}
