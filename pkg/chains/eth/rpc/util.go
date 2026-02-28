package rpc

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// decodeBig decodes a hex-encoded big.Int, treating empty or "0x" as zero.
// Handles different response formats between Geth and Parity.
func decodeBig(input string) (*big.Int, error) {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	return hexutil.DecodeBig(input)
}

// setZeroHex normalises an empty or bare "0x" string to "0x0".
func setZeroHex(input string) string {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	return input
}

// decodeString decodes a hex string to *big.Int, returning zero on error.
func decodeString(val string) *big.Int {
	decoded, err := hexutil.DecodeBig(val)
	if err != nil {
		return new(big.Int).SetUint64(0)
	}
	return decoded
}

// toCallArg converts an ethereum.CallMsg to the map format expected by eth_estimateGas
// and eth_sendTransaction.
func toCallArg(msg *ethereum.CallMsg) any {
	arg := map[string]any{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

// convertBlockRawInfo converts a BlockRawInfo (hex strings) to a BlockInfo (*big.Int fields).
func convertBlockRawInfo(raw *BlockRawInfo) *BlockInfo {
	var baseFeePerGas *big.Int
	if raw.BaseFeePerGas != "" {
		baseFeePerGas = decodeString(raw.BaseFeePerGas)
	}

	return &BlockInfo{
		Number:           decodeString(raw.Number),
		Hash:             raw.Hash,
		ParentHash:       raw.ParentHash,
		Nonce:            decodeString(raw.Nonce),
		Sha3Uncles:       raw.Sha3Uncles,
		LogsBloom:        raw.LogsBloom,
		TransactionsRoot: raw.TransactionsRoot,
		StateRoot:        raw.StateRoot,
		Miner:            raw.Miner,
		Difficulty:       decodeString(raw.Difficulty),
		TotalDifficulty:  decodeString(raw.TotalDifficulty),
		ExtraData:        raw.ExtraData,
		Size:             decodeString(raw.Size),
		GasLimit:         decodeString(raw.GasLimit),
		GasUsed:          decodeString(raw.GasUsed),
		Timestamp:        decodeString(raw.Timestamp),
		Transactions:     raw.Transactions,
		Uncles:           raw.Uncles,
		BaseFeePerGas:    baseFeePerGas,
	}
}

func castToInt64(val any) (int64, error) {
	v, err := castToString(val)
	if err != nil {
		return 0, err
	}

	v2, err := hexutil.DecodeBig(v)
	if err != nil {
		return 0, err
	}
	return v2.Int64(), nil
}

func castToString(val any) (string, error) {
	v, ok := val.(string)
	if !ok {
		return "", errors.New("fail to cast to string")
	}
	return v, nil
}

func castToSliceString(val any) ([]string, error) {
	data, ok := val.([]any)
	if !ok {
		return nil, errors.New("fail to cast to []string")
	}
	if len(data) == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, len(data))
	for _, v := range data {
		sRet, err := castToString(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, sRet)
	}
	return ret, nil
}

// ValidateAddr validates that addr is a valid hex Ethereum address.
func ValidateAddr(addr string) error {
	if !common.IsHexAddress(addr) {
		return errors.New("address:" + addr + " is invalid")
	}
	return nil
}
