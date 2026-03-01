package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// Sign calculates an Ethereum specific signature with:
//
//	sign(keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)))
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (e *Ethereum) Sign(ctx context.Context, hexAddr, message string) (string, error) {
	return ethrpc.Sign(ctx, e.rpcClient, hexAddr, message)
}

// SendTransaction sends transaction and returns transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction
func (e *Ethereum) SendTransaction(ctx context.Context, msg *ethereum.CallMsg) (string, error) {
	return ethrpc.SendTransaction(ctx, e.rpcClient, msg)
}

// SendRawTransaction creates new message call transaction or a contract creation for signed transactions
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction
func (e *Ethereum) SendRawTransaction(ctx context.Context, signedTx string) (string, error) {
	return ethrpc.SendRawTransaction(ctx, e.rpcClient, signedTx)
}

// SendRawTransactionWithTypesTx call SendRawTransaction() by types.Transaction.
// Uses EIP-2718 binary format (MarshalBinary) to correctly encode both legacy
// and typed (EIP-1559) transactions for eth_sendRawTransaction.
func (e *Ethereum) SendRawTransactionWithTypesTx(ctx context.Context, tx *types.Transaction) (string, error) {
	return ethrpc.SendRawTransactionWithTypesTx(ctx, e.rpcClient, tx)
}

// GetTransactionByHash returns the information about a transaction requested by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash
func (e *Ethereum) GetTransactionByHash(
	ctx context.Context, hashTx string,
) (*ethrpc.ResponseGetTransaction, error) {
	return ethrpc.GetTransactionByHash(ctx, e.rpcClient, hashTx)
}

// GetTransactionReceipt returns the receipt of a transaction by transaction hash
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt
func (e *Ethereum) GetTransactionReceipt(
	ctx context.Context, hashTx string,
) (*ethrpc.ResponseGetTransactionReceipt, error) {
	return ethrpc.GetTransactionReceipt(ctx, e.rpcClient, hashTx)
}

// GetTxReceipt retrieves a transaction receipt and converts it to the clean domain type.
// Returns (nil, nil) when the transaction has not yet been included in a block.
// Returns (nil, error) on node connectivity or parsing failures.
func (e *Ethereum) GetTxReceipt(ctx context.Context, txHash string) (*domainETH.TransactionReceipt, error) {
	result, err := ethrpc.GetTxReceipt(ctx, e.rpcClient, txHash)
	if err != nil {
		return nil, fmt.Errorf("fail to call GetTxReceipt(): %w", err)
	}
	return ToDomainTransactionReceiptFromPkg(result), nil
}
