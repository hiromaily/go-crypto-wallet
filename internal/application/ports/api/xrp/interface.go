// Package xrp defines interfaces for Ripple/XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
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
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.TxInput, string, error)

	// ripple
	Close() error
	CoinTypeCode() domainCoin.CoinTypeCode
	GetChainConf() *chaincfg.Params
}

// RippleAPIer defines the interface for Ripple API operations.
// Implementations handle account management, address generation, and transaction operations.
type RippleAPIer interface {
	// RippleAccountAPI
	GetAccountInfo(ctx context.Context, address string) (*dtoRipple.ResponseGetAccountInfo, error)
	// RippleAddressAPI
	GenerateAddress(ctx context.Context) (*dtoRipple.ResponseGenerateAddress, error)
	GenerateXAddress(ctx context.Context) (*dtoRipple.ResponseGenerateXAddress, error)
	IsValidAddress(ctx context.Context, addr string) (bool, error)
	// RippleTxAPI
	PrepareTransaction(
		ctx context.Context,
		senderAccount, receiverAccount string,
		amount float64,
		instructions *dtoRipple.Instructions,
	) (*dtoRipple.TxInput, string, error)
	SignTransaction(ctx context.Context, txJSON *dtoRipple.TxInput, secret string) (string, string, error)
	CombineTransaction(ctx context.Context, signedTxs []string) (string, string, error)
	SubmitTransaction(ctx context.Context, signedTx string) (*dtoRipple.SentTx, uint64, error)
	WaitValidation(ctx context.Context, targetledgerVarsion uint64) (uint64, error)
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoRipple.TxInfo, error)
}

// RipplePublicer defines the interface for Ripple public node operations.
// These operations query public information from the Ripple network.
type RipplePublicer interface {
	// public_account
	AccountChannels(ctx context.Context, sender, receiver string) (*dtoRipple.ResponseAccountChannels, error)
	AccountInfo(ctx context.Context, address string) (*dtoRipple.ResponseAccountInfo, error)
	// public_server_info
	ServerInfo(ctx context.Context) (*dtoRipple.ResponseServerInfo, error)
}

// RippleAdminer defines the interface for Ripple admin node operations.
// These operations typically require admin access to the Ripple node.
type RippleAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*dtoRipple.ResponseValidationCreate, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType dtoRipple.XRPKeyType,
	) (*dtoRipple.ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*dtoRipple.ResponseWalletPropose, error)
}
