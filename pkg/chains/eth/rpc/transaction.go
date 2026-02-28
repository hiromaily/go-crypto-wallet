package rpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// errReceiptNotFound is returned by GetTransactionReceipt when the node responds
// with null (transaction not yet included in a block).
var errReceiptNotFound = errors.New("response is empty")

// Sign calculates an Ethereum-specific signature.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func Sign(ctx context.Context, caller RPCCaller, hexAddr, message string) (string, error) {
	var signature string
	err := caller.CallContext(ctx, &signature, "eth_sign", hexAddr, message)
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(eth_sign): %w", err)
	}

	return signature, nil
}

// SendTransaction sends a transaction and returns the transaction hash.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction
func SendTransaction(ctx context.Context, caller RPCCaller, msg *ethereum.CallMsg) (string, error) {
	var txHash string
	err := caller.CallContext(ctx, &txHash, "eth_sendTransaction", toCallArg(msg))
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(eth_sendTransaction): %w", err)
	}

	return txHash, nil
}

// SendRawTransaction sends a signed raw transaction.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func SendRawTransaction(ctx context.Context, caller RPCCaller, signedTx string) (string, error) {
	var txHash string
	err := caller.CallContext(ctx, &txHash, "eth_sendRawTransaction", signedTx)
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(eth_sendTransaction): %w", err)
	}

	return txHash, nil
}

// SendRawTransactionWithTypesTx encodes tx using EIP-2718 binary format and calls
// SendRawTransaction. Correctly handles both legacy and typed (EIP-1559) transactions.
func SendRawTransactionWithTypesTx(ctx context.Context, caller RPCCaller, tx *types.Transaction) (string, error) {
	encodedTx, err := tx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("fail to call tx.MarshalBinary(): %w", err)
	}
	return SendRawTransaction(ctx, caller, hexutil.Encode(encodedTx))
}

// GetTransactionByHash returns transaction information by hash.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func GetTransactionByHash(
	ctx context.Context, caller RPCCaller, hashTx string,
) (*ResponseGetTransaction, error) {
	var resMap map[string]string
	err := caller.CallContext(ctx, &resMap, "eth_getTransactionByHash", hashTx)
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_getTransactionByHash): %w", err)
	}
	if len(resMap) == 0 {
		return nil, errors.New("response of eth_getTransactionByHash is empty")
	}

	blockNumber, err := hexutil.DecodeBig(setZeroHex(resMap["blockNumber"]))
	if err != nil {
		return nil, errors.New("response[blockNumber] is invalid")
	}
	gas, err := hexutil.DecodeBig(setZeroHex(resMap["gas"]))
	if err != nil {
		return nil, errors.New("response[gas] is invalid")
	}
	gasPrice, err := hexutil.DecodeBig(setZeroHex(resMap["gasPrice"]))
	if err != nil {
		return nil, errors.New("response[gasPrice] is invalid")
	}
	nonce, err := hexutil.DecodeBig(setZeroHex(resMap["nonce"]))
	if err != nil {
		return nil, errors.New("response[nonce] is invalid")
	}
	transactionIndex, err := hexutil.DecodeBig(setZeroHex(resMap["transactionIndex"]))
	if err != nil {
		return nil, errors.New("response[transactionIndex] is invalid")
	}
	value, err := hexutil.DecodeBig(setZeroHex(resMap["value"]))
	if err != nil {
		return nil, errors.New("response[value] is invalid")
	}
	v, err := hexutil.DecodeBig(setZeroHex(resMap["v"]))
	if err != nil {
		return nil, errors.New("response[v] is invalid")
	}

	return &ResponseGetTransaction{
		BlockHash:        resMap["blockHash"],
		BlockNumber:      blockNumber.Int64(),
		From:             resMap["from"],
		Gas:              gas.Int64(),
		GasPrice:         gasPrice.Int64(),
		Hash:             resMap["hash"],
		Input:            resMap["input"],
		Nonce:            nonce.Int64(),
		To:               resMap["to"],
		TransactionIndex: transactionIndex.Int64(),
		Value:            value.Int64(),
		V:                v.Int64(),
		R:                resMap["r"],
		S:                resMap["s"],
	}, nil
}

// GetTransactionReceipt returns the receipt of a transaction by hash.
// Returns errReceiptNotFound when the transaction is not yet included in a block.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
//
//nolint:gocyclo
func GetTransactionReceipt(
	ctx context.Context, caller RPCCaller, hashTx string,
) (*ResponseGetTransactionReceipt, error) {
	ch := make(chan error, 1)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		cancel()
	}()

	var resMap map[string]any
	go func() {
		err := caller.CallContext(timeoutCtx, &resMap, "eth_getTransactionReceipt", hashTx)
		if err != nil {
			ch <- fmt.Errorf("fail to call rpc.CallContext(eth_getTransactionReceipt): %w", err)
		}
		ch <- nil
	}()

	select {
	case <-timeoutCtx.Done():
		err := timeoutCtx.Err()
		switch err {
		case context.Canceled:
			logger.Debug("context.Canceled for calling eth_getTransactionReceipt")
		case context.DeadlineExceeded:
			logger.Debug("context.DeadlineExceeded for calling eth_getTransactionReceipt")
		case nil:
			// no error
		default:
			logger.Debug(err.Error())
			return nil, err
		}
	case retErr := <-ch:
		if retErr != nil {
			return nil, retErr
		}
	}

	if len(resMap) == 0 {
		return nil, errReceiptNotFound
	}

	transactionHash, err := castToString(resMap["transactionHash"])
	if err != nil {
		return nil, errors.New("response[transactionHash] is invalid")
	}
	transactionIndex, err := castToInt64(resMap["transactionIndex"])
	if err != nil {
		return nil, errors.New("response[transactionIndex] is invalid")
	}
	blockHash, err := castToString(resMap["blockHash"])
	if err != nil {
		return nil, errors.New("response[blockHash] is invalid")
	}
	blockNumber, err := castToInt64(resMap["blockNumber"])
	if err != nil {
		return nil, errors.New("response[blockNumber] is invalid")
	}
	from, err := castToString(resMap["from"])
	if err != nil {
		return nil, errors.New("response[from] is invalid")
	}
	to, err := castToString(resMap["to"])
	if err != nil {
		return nil, errors.New("response[to] is invalid")
	}
	cumulativeGasUsed, err := castToInt64(resMap["cumulativeGasUsed"])
	if err != nil {
		return nil, errors.New("response[cumulativeGasUsed] is invalid")
	}
	gasUsed, err := castToInt64(resMap["gasUsed"])
	if err != nil {
		return nil, errors.New("response[gasUsed] is invalid")
	}
	var contractAddress string
	if resMap["contractAddress"] != nil {
		contractAddress, err = castToString(resMap["contractAddress"])
		if err != nil {
			return nil, errors.New("response[contractAddress] is invalid")
		}
	}
	logs, err := castToSliceString(resMap["logs"])
	if err != nil {
		return nil, errors.New("response[logs] is invalid")
	}
	logsBloom, err := castToString(resMap["logsBloom"])
	if err != nil {
		return nil, errors.New("response[logsBloom] is invalid")
	}
	status, err := castToInt64(resMap["status"])
	if err != nil {
		return nil, errors.New("response[status] is invalid")
	}

	return &ResponseGetTransactionReceipt{
		TransactionHash:   transactionHash,
		TransactionIndex:  transactionIndex,
		BlockHash:         blockHash,
		BlockNumber:       blockNumber,
		From:              from,
		To:                to,
		CumulativeGasUsed: cumulativeGasUsed,
		GasUsed:           gasUsed,
		ContractAddress:   contractAddress,
		Logs:              logs,
		LogsBloom:         logsBloom,
		Status:            status,
	}, nil
}

// GetTxReceipt retrieves a transaction receipt and returns the decoded TransactionReceipt.
// Returns (nil, nil) when the transaction has not yet been included in a block.
func GetTxReceipt(ctx context.Context, caller RPCCaller, txHash string) (*TransactionReceipt, error) {
	resp, err := GetTransactionReceipt(ctx, caller, txHash)
	if err != nil {
		if errors.Is(err, errReceiptNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("fail to call GetTransactionReceipt(): %w", err)
	}
	return &TransactionReceipt{
		TransactionHash:   resp.TransactionHash,
		TransactionIndex:  uint64(resp.TransactionIndex),
		BlockHash:         resp.BlockHash,
		BlockNumber:       uint64(resp.BlockNumber),
		From:              resp.From,
		To:                resp.To,
		CumulativeGasUsed: uint64(resp.CumulativeGasUsed),
		GasUsed:           uint64(resp.GasUsed),
		ContractAddress:   resp.ContractAddress,
		Status:            uint64(resp.Status),
	}, nil
}
