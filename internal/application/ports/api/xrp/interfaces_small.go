// Package xrp defines interfaces for XRP blockchain operations.
//
// This file contains small, focused interfaces following the Interface Segregation Principle.
// Use cases should depend only on the specific interfaces they need through local interface
// compositions.
package xrp

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

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

// LedgerWaiter waits for ledger validation.
//
// Deprecated: Use TransactionSubmitter.WaitValidation instead.
// This interface is kept for backward compatibility.
type LedgerWaiter interface {
	WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error)
}

// TransactionGetter retrieves transaction information.
//
// Deprecated: Use TransactionSubmitter.GetTransaction instead.
// This interface is kept for backward compatibility.
type TransactionGetter interface {
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*dtoxrp.TxInfo, error)
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
	WalletPropose(ctx context.Context, passphrase string) (*dtoxrp.ResponseWalletPropose, error)
}

// Closer provides cleanup operations.
type Closer interface {
	Close() error
}
