// Package ripple defines interfaces for Ripple/XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
//
// Note: This package imports XRP infrastructure types to define the interface.
// This is acceptable because the interface is the abstraction and the infrastructure
// implements it. The dependency direction is: infrastructure -> ports (interface).
package ripple

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
)

// Rippler defines the main interface for Ripple/XRP blockchain operations.
// It embeds specialized interfaces for admin, public, and API operations.
type Rippler interface {
	RippleAdminer
	RipplePublicer
	RippleAPIer

	// balance
	GetBalance(ctx context.Context, addr string) (float64, error)
	GetTotalBalance(ctx context.Context, addrs []string) float64

	// transaction
	CreateRawTransaction(
		ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions,
	) (*xrp.TxInput, string, error)

	// ripple
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// RippleAPIer defines the interface for Ripple API operations.
// Implementations handle account management, address generation, and transaction operations.
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(ctx context.Context, address string) (*xrp.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress(ctx context.Context) (*xrp.ResponseGenerateAddress, error)
	GenerateXAddress(ctx context.Context) (*xrp.ResponseGenerateXAddress, error)
	IsValidAddress(ctx context.Context, addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(
		ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *xrp.Instructions,
	) (*xrp.TxInput, string, error)
	SignTransaction(ctx context.Context, txJSON *xrp.TxInput, secret string) (string, string, error)
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
	SubmitTransaction(ctx context.Context, signedTx string) (*xrp.SentTx, uint64, error)
	WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*xrp.TxInfo, error)
}

// RipplePublicer defines the interface for Ripple public node operations.
// These operations query public information from the Ripple network.
type RipplePublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*xrp.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*xrp.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*xrp.ResponseServerInfo, error)
}

// RippleAdminer defines the interface for Ripple admin node operations.
// These operations typically require admin access to the Ripple node.
type RippleAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*xrp.ResponseValidationCreate, error)
	WalletProposeWithKey(ctx context.Context, seed string, keyType xrp.XRPKeyType) (*xrp.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*xrp.ResponseWalletPropose, error)
}
