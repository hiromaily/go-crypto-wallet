// Package network provides network connection infrastructure.
//
// This package contains:
//   - websocket/: WebSocket connection management for real-time communication
//
// Network infrastructure is responsible for:
//   - Establishing and managing network connections
//   - Handling connection lifecycle (connect, disconnect, reconnect)
//   - Managing connection errors and timeouts
//   - Providing connection objects to API clients
//
// Network infrastructure does NOT:
//   - Contain business logic
//   - Handle application-level protocols (that's for API clients)
//   - Make routing or message handling decisions
//
// Network components are utilities used by API clients for establishing
// connections to external systems.
package network
