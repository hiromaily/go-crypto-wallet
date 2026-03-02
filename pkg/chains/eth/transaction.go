package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// EncodeTx encodes a transaction to hex using EIP-2718 binary format.
// This correctly handles both legacy (Type 0) and typed (Type 2 EIP-1559) transactions.
// Unlike rlp.EncodeToBytes, MarshalBinary produces the raw EIP-2718 envelope
// expected by eth_sendRawTransaction.
func EncodeTx(tx *types.Transaction) (*string, error) {
	txb, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("fail to marshal transaction to binary: %w", err)
	}
	txHex := hexutil.Encode(txb)
	return &txHex, nil
}

// DecodeTx decodes a hex-encoded transaction using EIP-2718 binary format.
// This correctly handles both legacy (Type 0) and typed (Type 2 EIP-1559) transactions.
func DecodeTx(txHex string) (*types.Transaction, error) {
	txc, err := hexutil.Decode(txHex)
	if err != nil {
		return nil, fmt.Errorf("fail to decode hex transaction: %w", err)
	}

	var txde types.Transaction
	if err = txde.UnmarshalBinary(txc); err != nil {
		return nil, fmt.Errorf("fail to unmarshal binary transaction: %w", err)
	}

	return &txde, nil
}
