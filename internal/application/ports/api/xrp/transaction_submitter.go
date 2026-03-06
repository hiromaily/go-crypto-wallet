package xrp

import (
	"context"

	xrprpcadmin "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
)

// LedgerPoller queries the current ledger index.
// Implemented by the public XRP client. Used by use cases that need to poll
// for ledger advancement (e.g., waiting for transaction validation).
type LedgerPoller interface {
	LedgerCurrent(ctx context.Context) (*xrprpcpublic.ResponseLedgerCurrent, error)
}

// LedgerAdvancer advances the ledger (admin operation, standalone mode only).
// Implemented by the admin XRP client. May be nil when admin is not configured.
type LedgerAdvancer interface {
	LedgerAccept(ctx context.Context) (*xrprpcadmin.ResponseLedgerAccept, error)
}

// TransactionSubmitter defines an interface for submitting and tracking XRP transactions.
//
// This interface follows the Interface Segregation Principle by providing only
// submission and retrieval operations. Ledger polling is handled separately via
// LedgerPoller and LedgerAdvancer — use cases coordinate ledger advancement
// themselves (Plan C), rather than delegating to a single WaitValidation method.
//
// # Constraints
//   - Online-only operations (requires network connectivity)
//   - No account query methods (those belong to AccountInfoProvider)
//   - No signing methods (those belong to TransactionSigner)
//   - Idempotent submission (submitting same blob returns same result)
type TransactionSubmitter interface {
	// SubmitTransaction broadcasts a signed transaction to the XRP Ledger.
	SubmitTransaction(ctx context.Context, signedTx string) (*xrplgo.SentTx, uint64, error)

	// GetTransaction retrieves detailed information about a submitted transaction.
	GetTransaction(ctx context.Context, txID string, targetLedgerVersion uint64) (*xrplgo.TxInfo, error)
}
