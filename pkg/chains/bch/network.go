package bch

import "github.com/btcsuite/btcd/wire"

// refer to original [github.com/cpacia/bchutil](https://github.com/cpacia/bchutil/blob/master/protocol.go)

const (
	// MainnetMagic represents the main bitcoin network.
	MainnetMagic wire.BitcoinNet = 0xe8f3e1e3

	// TestnetMagic represents the test network (version 3).
	TestnetMagic wire.BitcoinNet = 0xf4f3e5f4

	// RegtestMagic represents the regression test network.
	RegtestMagic wire.BitcoinNet = 0xfabfb5da
)
