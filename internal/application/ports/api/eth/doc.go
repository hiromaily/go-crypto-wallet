// Package eth defines interfaces for Ethereum blockchain operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/api/eth/).
//
// # Interface Categories
//
// ## Token / Legacy Monitoring Interfaces
//
//   - ERC20er: Interface for ERC-20 token operations
//   - EtherTxMonitor: Legacy interface for monitoring transaction confirmations
//
// ## ISP-Compliant Lifecycle / Key Interfaces (PR #575)
//
//   - ETHLifecycle: Client lifecycle (Close, CoinTypeCode)
//   - ETHKeyAccessor: Keystore access for private key import
//   - ETHTransactionSigner: Sign raw transactions via keystore passphrase
//   - ETHRawKeyImporter: Import raw keys via RPC
//   - ETHNodeAPIClient: Node API operations (clientVersion, netVersion, etc.)
//   - ETHKeygenSignClient: Composed — ETHLifecycle + ETHRawKeyImporter
//   - ETHWatchClient: Composed — ETHLifecycle + ETHNodeAPIClient
//
// ## ISP-Compliant EIP-1559 Transaction Flow Interfaces (Task 2.1b)
//
//   - ChainConfigProvider: Chain-level configuration (CoinTypeCode, GetChainConf)
//   - BalanceChecker: Account balance queries
//   - TxCreator: Create legacy and EIP-1559 unsigned transactions
//   - GasEstimator: Gas price and tip cap estimation
//   - TxSigner: Offline signing using *ecdsa.PrivateKey (no keystore required)
//   - TxSender: Broadcast signed transactions (ETHTransactionSender is a type alias)
//   - TxMonitor: Transaction receipt and confirmation count
//   - AddressValidator: Ethereum address validation
//   - WatchTxCreationDeps: Composed — ChainConfigProvider + TxCreator + GasEstimator + AddressValidator
//   - KeygenSignTxDeps: Composed — ChainConfigProvider + TxSigner
//
// # Types
//
//   - TxCreateParams: DTO for transaction creation parameters (port-level, no infrastructure types)
//
// # Usage
//
// Use cases depend on the small, focused interfaces — never on Ethereumer directly:
//
//	import apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
//
//	// Watch wallet: create transaction
//	type createTxUseCase struct {
//	    eth apieth.WatchTxCreationDeps
//	}
//
//	// Keygen wallet: sign transaction offline
//	type signTxUseCase struct {
//	    eth apieth.KeygenSignTxDeps
//	}
//
// # Related Packages
//
//   - internal/infrastructure/api/eth/eth/: ETH implementation
//   - internal/infrastructure/api/eth/erc20/: ERC-20 implementation
//   - internal/domain/ethereum/: Domain types used in interface methods
package eth
