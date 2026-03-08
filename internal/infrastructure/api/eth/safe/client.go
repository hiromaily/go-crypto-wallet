// Package safe provides infrastructure for interacting with Safe (Gnosis Safe v1.4.1)
// smart contract accounts via the generated ABI bindings.
package safe

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	safeabi "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/contract/safe"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// Compile-time check that SafeClient implements SafeClientDeps.
var _ apieth.SafeClientDeps = (*SafeClient)(nil)

// NewSafeClientParams holds DI-provided signing parameters for the SafeClient.
// From is the address whose ETH pays gas for execTransaction submissions.
// SignerFn signs the outer Ethereum transaction; when nil ExecuteSafeTransaction
// returns an error.
type NewSafeClientParams struct {
	From     common.Address
	SignerFn bind.SignerFn
}

// SafeClient implements SafeClientDeps — the four Safe port interfaces —
// using the abigen-generated Safe v1.4.1 ABI bindings and a go-ethereum JSON-RPC client.
type SafeClient struct {
	client   *ethclient.Client
	chainID  *big.Int
	from     common.Address
	signerFn bind.SignerFn
}

// NewSafeClient constructs a SafeClient.
// chainID is fetched from the node once and cached.
func NewSafeClient(ctx context.Context, client *ethclient.Client, params NewSafeClientParams) (*SafeClient, error) {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("safe: failed to fetch chain ID: %w", err)
	}
	return &SafeClient{
		client:   client,
		chainID:  chainID,
		from:     params.From,
		signerFn: params.SignerFn,
	}, nil
}

// GetSafeTxHash calls getTransactionHash on the Safe contract and returns the
// result as a 0x-prefixed hex string.
// All gas parameters are zero for a simple ETH transfer — the Safe spec requires
// them to be present but they do not affect the hash when zero.
func (s *SafeClient) GetSafeTxHash(
	ctx context.Context,
	safeAddr, to string,
	value *big.Int,
	data []byte,
	operation uint8,
	nonce *big.Int,
) (string, error) {
	safe, err := safeabi.NewSafe(common.HexToAddress(safeAddr), s.client)
	if err != nil {
		return "", fmt.Errorf("safe: failed to bind Safe at %s: %w", safeAddr, err)
	}

	opts := &bind.CallOpts{Context: ctx}
	toAddr := common.HexToAddress(to)
	zero := new(big.Int)

	hash, err := safe.GetTransactionHash(
		opts,
		toAddr,
		value,
		data,
		operation,
		zero,             // safeTxGas
		zero,             // baseGas
		zero,             // gasPrice
		common.Address{}, // gasToken (zero address = ETH)
		common.Address{}, // refundReceiver (zero address)
		nonce,
	)
	if err != nil {
		return "", fmt.Errorf("safe: GetTransactionHash failed: %w", err)
	}

	return "0x" + hex.EncodeToString(hash[:]), nil
}

// GetSafeNonce calls Nonce() on the Safe contract binding and returns the
// current nonce as a *big.Int.
func (s *SafeClient) GetSafeNonce(ctx context.Context, safeAddr string) (*big.Int, error) {
	safe, err := safeabi.NewSafe(common.HexToAddress(safeAddr), s.client)
	if err != nil {
		return nil, fmt.Errorf("safe: failed to bind Safe at %s: %w", safeAddr, err)
	}

	opts := &bind.CallOpts{Context: ctx}
	nonce, err := safe.Nonce(opts)
	if err != nil {
		return nil, fmt.Errorf("safe: Nonce() failed: %w", err)
	}
	return nonce, nil
}

// ExecuteSafeTransaction builds an EIP-1559 transaction calling execTransaction
// on the Safe contract, estimates gas, submits it, and returns the transaction hash.
func (s *SafeClient) ExecuteSafeTransaction(ctx context.Context, params apieth.SafeExecParams) (string, error) {
	if s.signerFn == nil {
		return "", errors.New("safe: no signer configured for ExecuteSafeTransaction")
	}

	safe, err := safeabi.NewSafe(common.HexToAddress(params.SafeAddress), s.client)
	if err != nil {
		return "", fmt.Errorf("safe: failed to bind Safe at %s: %w", params.SafeAddress, err)
	}

	// Build EIP-1559 gas fields.
	tip, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("safe: SuggestGasTipCap failed: %w", err)
	}

	// maxFeePerGas = baseFee * 2 + tip (provides headroom for base-fee increases).
	header, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("safe: HeaderByNumber failed: %w", err)
	}
	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee = new(big.Int) // pre-London chain: fall back to 0
	}
	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tip)

	// Estimate gas for the execTransaction call.
	toAddr := common.HexToAddress(params.SafeAddress)
	abiPacked, err := packExecTransaction(params)
	if err != nil {
		return "", fmt.Errorf("safe: failed to pack execTransaction call: %w", err)
	}
	estimated, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From: s.from,
		To:   &toAddr,
		Data: abiPacked,
	})
	if err != nil {
		return "", fmt.Errorf("safe: EstimateGas failed: %w", err)
	}
	// Add 20% buffer so the tx doesn't run out of gas due to estimation variance.
	gasLimit := estimated * 12 / 10

	logger.Info("safe: ExecuteSafeTransaction",
		"safeAddress", params.SafeAddress,
		"to", params.To,
		"value", params.Value,
		"estimatedGas", estimated,
		"gasLimit", gasLimit,
	)

	// Build transact opts for EIP-1559.
	txOpts := &bind.TransactOpts{
		From:      s.from,
		Signer:    s.signerFn,
		GasFeeCap: gasFeeCap,
		GasTipCap: tip,
		GasLimit:  gasLimit,
		Context:   ctx,
	}

	tx, err := safe.ExecTransaction(
		txOpts,
		common.HexToAddress(params.To),
		params.Value,
		params.Data,
		params.Operation,
		params.SafeTxGas,
		params.BaseGas,
		params.GasPrice,
		common.HexToAddress(params.GasToken),
		common.HexToAddress(params.RefundReceiver),
		params.Signatures,
	)
	if err != nil {
		return "", fmt.Errorf("safe: ExecTransaction failed: %w", err)
	}

	return tx.Hash().Hex(), nil
}

// GetSafeInfo retrieves on-chain state of a Safe contract.
func (s *SafeClient) GetSafeInfo(ctx context.Context, safeAddr string) (*apieth.SafeInfo, error) {
	addr := common.HexToAddress(safeAddr)
	safe, err := safeabi.NewSafe(addr, s.client)
	if err != nil {
		return nil, fmt.Errorf("safe: failed to bind Safe at %s: %w", safeAddr, err)
	}

	opts := &bind.CallOpts{Context: ctx}

	owners, err := safe.GetOwners(opts)
	if err != nil {
		return nil, fmt.Errorf("safe: GetOwners() failed: %w", err)
	}

	threshold, err := safe.GetThreshold(opts)
	if err != nil {
		return nil, fmt.Errorf("safe: GetThreshold() failed: %w", err)
	}

	nonce, err := safe.Nonce(opts)
	if err != nil {
		return nil, fmt.Errorf("safe: Nonce() failed: %w", err)
	}

	balance, err := s.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("safe: BalanceAt() failed: %w", err)
	}

	ownerStrs := make([]string, len(owners))
	for i, o := range owners {
		ownerStrs[i] = o.Hex()
	}

	return &apieth.SafeInfo{
		Owners:    ownerStrs,
		Threshold: threshold.Uint64(),
		Nonce:     nonce.Uint64(),
		Balance:   balance,
	}, nil
}

// packExecTransaction ABI-encodes the execTransaction call data for gas estimation.
func packExecTransaction(params apieth.SafeExecParams) ([]byte, error) {
	safeABI, err := safeabi.SafeMetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to get Safe ABI: %w", err)
	}

	toAddr := common.HexToAddress(params.To)
	gasTokenAddr := common.HexToAddress(params.GasToken)
	refundReceiverAddr := common.HexToAddress(params.RefundReceiver)

	packed, err := safeABI.Pack(
		"execTransaction",
		toAddr,
		params.Value,
		params.Data,
		params.Operation,
		params.SafeTxGas,
		params.BaseGas,
		params.GasPrice,
		gasTokenAddr,
		refundReceiverAddr,
		params.Signatures,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execTransaction: %w", err)
	}
	return packed, nil
}
