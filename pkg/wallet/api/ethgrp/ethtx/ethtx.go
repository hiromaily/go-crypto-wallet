package ethtx

import (
	"math/big"
)

// RawTx is raw transaction
type RawTx struct {
	UUID  string  `json:"uuid"`
	From  string  `json:"from"`
	To    string  `json:"to"`
	Value big.Int `json:"value"`
	Nonce uint64  `json:"nonce"`
	TxHex string  `json:"txhex"`
	Hash  string  `json:"hash"`
}
