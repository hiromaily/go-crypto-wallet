package ethtx

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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

func EncodeTx(tx *types.Transaction) (*string, error) {
	txb, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	txHex := hexutil.Encode(txb)
	return &txHex, nil
}

func DecodeTx(txHex string) (*types.Transaction, error) {
	txc, err := hexutil.Decode(txHex)
	if err != nil {
		return nil, err
	}

	var txde types.Transaction
	err = rlp.Decode(bytes.NewReader(txc), &txde)
	if err != nil {
		return nil, err
	}

	return &txde, nil
}
