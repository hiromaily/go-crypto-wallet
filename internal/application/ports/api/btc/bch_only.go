package btc

// BCHer is the interface for Bitcoin Cash (BCH) operations.
// Currently it is identical to BitcoinCompatible, but exists as a separate type
// to allow BCH-specific methods to be added in the future without affecting BTC.
//
// BCH-specific characteristics:
//   - Uses Raw Transaction Hex format (not PSBT)
//   - Uses SIGHASH_FORKID (0x41) for transaction signing
//   - Uses CashAddr format for addresses (bitcoincash:q...)
//   - No SegWit support (no P2WPKH, P2WSH, bech32)
//   - No Taproot support (no P2TR, bech32m, Schnorr)
//   - No Descriptor support (uses ImportAddress instead)
//
// Future BCH-specific methods could include:
//   - CashAddr encoding/decoding utilities
//   - CashToken support (when implemented)
//   - SLP (Simple Ledger Protocol) token support
//
//nolint:iface // BCHer will have BCH-specific methods in the future
type BCHer interface {
	BitcoinCompatible

	// Future BCH-specific methods can be added here
	// Example:
	// EncodeCashAddr(addr string) (string, error)
	// DecodeCashAddr(cashAddr string) (string, error)
}
