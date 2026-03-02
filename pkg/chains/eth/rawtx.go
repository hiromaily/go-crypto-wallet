package eth

import "math/big"

// RawTx represents a raw Ethereum transaction before being submitted to the network.
type RawTx struct {
	UUID  string  `json:"uuid"`
	From  string  `json:"from"`
	To    string  `json:"to"`
	Value big.Int `json:"value"`
	Nonce uint64  `json:"nonce"`
	TxHex string  `json:"txhex"`
	Hash  string  `json:"hash"`
}
