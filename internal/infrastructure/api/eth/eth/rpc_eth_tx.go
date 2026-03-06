package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// Sign calculates an Ethereum specific signature with:
//
//	sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)))
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (e *Ethereum) Sign(ctx context.Context, hexAddr, message string) (string, error) {
	return e.pkgrpc.Sign(ctx, hexAddr, message)
}

// SendTransaction sends transaction and returns transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction
func (e *Ethereum) SendTransaction(ctx context.Context, msg *ethereum.CallMsg) (string, error) {
	return e.pkgrpc.SendTransaction(ctx, msg)
}

// SendRawTransaction creates new message call transaction or a contract creation for signed transactions
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func (e *Ethereum) SendRawTransaction(ctx context.Context, signedTx string) (string, error) {
	return e.pkgrpc.SendRawTransaction(ctx, signedTx)
}

// SendRawTransactionWithTypesTx call SendRawTransaction() by types.Transaction.
// Uses EIP-2718 binary format (MarshalBinary) to correctly encode both legacy
// and typed (EIP-1559) transactions for eth_sendRawTransaction.
func (e *Ethereum) SendRawTransactionWithTypesTx(ctx context.Context, tx *types.Transaction) (string, error) {
	return e.pkgrpc.SendRawTransactionWithTypesTx(ctx, tx)
}

// GetTransactionByHash returns the information about a transaction requested by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func (e *Ethereum) GetTransactionByHash(
	ctx context.Context, hashTx string,
) (*dtoeth.TransactionInfo, error) {
	result, err := e.pkgrpc.GetTransactionByHash(ctx, hashTx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &dtoeth.TransactionInfo{
		BlockHash:        result.BlockHash,
		BlockNumber:      result.BlockNumber,
		From:             result.From,
		Gas:              result.Gas,
		GasPrice:         result.GasPrice,
		Hash:             result.Hash,
		Input:            result.Input,
		Nonce:            result.Nonce,
		To:               result.To,
		TransactionIndex: result.TransactionIndex,
		Value:            result.Value,
		V:                result.V,
		R:                result.R,
		S:                result.S,
	}, nil
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
func (e *Ethereum) GetTransactionReceipt(
	ctx context.Context, hashTx string,
) (*dtoeth.RawTransactionReceipt, error) {
	result, err := e.pkgrpc.GetTransactionReceipt(ctx, hashTx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &dtoeth.RawTransactionReceipt{
		TransactionHash:   result.TransactionHash,
		TransactionIndex:  result.TransactionIndex,
		BlockHash:         result.BlockHash,
		BlockNumber:       result.BlockNumber,
		From:              result.From,
		To:                result.To,
		CumulativeGasUsed: result.CumulativeGasUsed,
		GasUsed:           result.GasUsed,
		ContractAddress:   result.ContractAddress,
		Logs:              result.Logs,
		LogsBloom:         result.LogsBloom,
		Status:            result.Status,
	}, nil
}

// GetTxReceipt retrieves a transaction receipt and converts it to the clean domain type.
// Returns (nil, nil) when the transaction has not yet been included in a block.
// Returns (nil, error) on node connectivity or parsing failures.
func (e *Ethereum) GetTxReceipt(ctx context.Context, txHash string) (*domainETH.TransactionReceipt, error) {
	result, err := e.pkgrpc.GetTxReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("fail to call GetTxReceipt(): %w", err)
	}
	return toDomainTransactionReceiptFromPkg(result), nil
}

// toDomainTransactionReceiptFromPkg converts pkg TransactionReceipt to domain.
func toDomainTransactionReceiptFromPkg(pkg *ethrpc.TransactionReceipt) *domainETH.TransactionReceipt {
	if pkg == nil {
		return nil
	}
	return &domainETH.TransactionReceipt{
		TransactionHash:   pkg.TransactionHash,
		TransactionIndex:  pkg.TransactionIndex,
		BlockHash:         pkg.BlockHash,
		BlockNumber:       pkg.BlockNumber,
		From:              pkg.From,
		To:                pkg.To,
		CumulativeGasUsed: pkg.CumulativeGasUsed,
		GasUsed:           pkg.GasUsed,
		ContractAddress:   pkg.ContractAddress,
		Status:            pkg.Status,
	}
}
