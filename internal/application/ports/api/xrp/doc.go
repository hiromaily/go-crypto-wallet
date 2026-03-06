// Package xrp defines interfaces for XRP blockchain operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/xrp/).
//
// # Interfaces
//
// Composed interfaces for use cases:
//   - XRPPublicClient: All public node operations (primary interface for new code)
//   - XRPAdminClient: Admin node operations (keygen wallets, standalone mode)
//   - XRPer: Deprecated monolithic interface (kept for backward compatibility)
//
// Small focused interfaces (Interface Segregation Principle):
//   - XRPPublicer, XRPAdminer: Node-specific operation groups
//   - AccountInfoProvider: Account info and balance queries
//   - BalanceChecker: Balance queries (single and aggregate)
//   - TransactionPreparer: Raw transaction creation
//   - TransactionSigner: Offline signing (air-gapped capable)
//   - TransactionSubmitter: Submission and retrieval
//   - LedgerPoller: Current ledger index queries
//   - LedgerAdvancer: Ledger advancement (admin, standalone mode)
//   - CoinTypeProvider, Closer: Supporting operations
//
// # Usage
//
// Use cases depend on the smallest interface needed:
//
//	type myUseCase struct {
//	    xrp apixrp.AccountInfoProvider  // not the full XRPPublicClient
//	}
package xrp
