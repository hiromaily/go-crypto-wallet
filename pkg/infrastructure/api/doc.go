// Package api provides external blockchain API client implementations.
//
// This package contains:
//   - bitcoin/: Bitcoin and Bitcoin Cash RPC client implementations
//   - ethereum/: Ethereum JSON-RPC and ERC-20 client implementations
//   - ripple/: Ripple (XRP) gRPC client implementation
//
// API clients are responsible for:
//   - Communicating with external blockchain node APIs
//   - Serializing and deserializing API requests/responses
//   - Handling network errors and retries
//   - Managing API connections
//
// API clients do NOT:
//   - Contain business logic
//   - Validate transactions (use domain validators)
//   - Make domain decisions about wallet operations
//
// API clients may implement domain service interfaces (if defined) or be wrapped
// by adapters to align with domain abstractions.
package api
