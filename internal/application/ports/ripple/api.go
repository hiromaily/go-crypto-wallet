// Package ripple defines interfaces for Ripple/XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
//
// TODO(architecture): ARCHITECTURAL DEBT - This package currently imports infrastructure
// types (xrp package) which violates Clean Architecture dependency direction.
// The dependency should flow: infrastructure -> ports, not ports -> infrastructure.
//
// Current state (INCORRECT):
//
//	application/ports/ripple (interface)
//	    ↓ imports
//	infrastructure/api/ripple/xrp (implementation)
//
// Desired state (CORRECT):
//
//	application/ports/ripple (interface)
//	    ↑ implemented by
//	infrastructure/api/ripple/xrp (implementation)
//
// This creates inconsistency with other blockchain interfaces (BTC, ETH) which only
// import external library types and domain types, never infrastructure types.
//
// Recommended solutions:
// 1. Move XRP protocol types (Instructions, TxInput, Response*) to internal/domain/xrp
// 2. Create application-layer DTOs in internal/application/dto/ripple
// 3. Use generic types (map[string]interface{}) for protocol-specific structures
//
// This technical debt should be addressed in a future refactoring to align with
// the project's Clean Architecture principles. See PR #240 review for full analysis.
package ripple

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	// TODO: Remove this infrastructure import (architectural debt)
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
