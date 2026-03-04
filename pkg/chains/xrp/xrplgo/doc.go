// Package xrplgo provides a WebSocket client for XRP Ledger (XRPL) node communication.
//
// It wraps the github.com/xrpscan/xrpl-go library and implements the interfaces
// defined in internal/application/ports/api/xrp, including:
//   - AccountInfoProvider: fetching account info and XRP balances
//   - TransactionSubmitter: submitting signed transactions and waiting for ledger validation
//
// # Usage
//
// Create a client with [NewClient] and a [ClientConfig]:
//
//	cfg := xrplgo.DefaultConfig("wss://s.altnet.rippletest.net:51233")
//	client, err := xrplgo.NewClient(cfg)
//	if err != nil {
//	    // handle error
//	}
//	defer client.Close()
//
// # Amount Representation
//
// XRP amounts are represented as XRP (not drops) in all public API responses.
// Internally, the XRPL protocol uses drops (1 XRP = 1,000,000 drops); conversion
// is handled transparently by this package.
package xrplgo
