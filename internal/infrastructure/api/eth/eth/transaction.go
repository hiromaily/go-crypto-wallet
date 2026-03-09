package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// when creating multiple transaction from same address, nonce should be incremented
func (e *Ethereum) getNonce(ctx context.Context, fromAddr string, additionalNonce int) (uint64, error) {
	// by calling GetTransactionCount()
	nonce, err := e.GetTransactionCount(ctx, fromAddr, domainETH.QuantityTagPending)
	if err != nil {
		return 0, fmt.Errorf("fail to call eth.GetTransactionCount(): %w", err)
	}
	if additionalNonce != 0 {
		nonce = nonce.Add(nonce, new(big.Int).SetUint64(uint64(additionalNonce)))
	}
	logger.Debug("nonce",
		"GetTransactionCount(fromAddr, QuantityTagPending)", nonce.Uint64(),
	)

	return nonce.Uint64(), nil
}

// SignOnRawTransaction signs on raw transaction
// - https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go
// - Note: this requires private key on this machine, if node is working remotely, it would not work.
func (e *Ethereum) SignOnRawTransaction(rawTx *domainETH.RawTx, passphrase string) (*domainETH.RawTx, error) {
	txHex := rawTx.TxHex
	fromAddr := rawTx.From
	tx, err := pkgeth.DecodeTx(txHex)
	if err != nil {
		return nil, fmt.Errorf("fail to call decodeTx(txHex): %w", err)
	}

	// get private key
	key, err := e.GetPrivKey(fromAddr, passphrase)
	if err != nil {
		return nil, fmt.Errorf("fail to call e.GetPrivKey(): %w", err)
	}

	// chain id
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
	chainID := big.NewInt(int64(e.netID))
	if chainID.Uint64() == 0 {
		return nil, fmt.Errorf("chainID can't get from netID:  %d", e.netID)
	}

	logger.Debug("call types.SignTx",
		"tx", tx,
		"chainID", chainID.Uint64(),
	)
	signer := types.LatestSignerForChainID(chainID)

	// sign
	signedTX, err := types.SignTx(tx, signer, key.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("fail to call types.SignTx(): %w", err)
	}

	// TODO: baseFee *big.Int param is added in AsMessage method and maybe useful
	fromSignedAddr, err := types.Sender(signer, signedTX)
	if err != nil {
		return nil, fmt.Errorf("fail to cll signedTX.AsMessage(): %w", err)
	}

	encodedTx, err := pkgeth.EncodeTx(signedTX)
	if err != nil {
		return nil, fmt.Errorf("fail to call encodeTx(): %w", err)
	}

	return &domainETH.RawTx{
		UUID:  rawTx.UUID,
		From:  fromSignedAddr.Hex(),
		To:    signedTX.To().Hex(),
		Value: *signedTX.Value(),
		Nonce: signedTX.Nonce(),
		TxHex: *encodedTx,
		Hash:  signedTX.Hash().Hex(),
	}, nil
}

// SendRawTx sends raw transaction
func (e *Ethereum) SendRawTx(ctx context.Context, tx *types.Transaction) error {
	return e.ethClient.SendTransaction(ctx, tx)
}

// SendSignedRawTransaction sends signed raw transaction
// - SendRawTransaction in rpc_eth_tx.go
// - SendRawTx in client.go
func (e *Ethereum) SendSignedRawTransaction(ctx context.Context, signedTxHex string) (string, error) {
	decodedTx, err := pkgeth.DecodeTx(signedTxHex)
	if err != nil {
		return "", fmt.Errorf("fail to call decodeTx(signedTxHex): %w", err)
	}

	txHash, err := e.SendRawTransactionWithTypesTx(ctx, decodedTx)
	if err != nil {
		return "", fmt.Errorf("fail to call SendRawTransactionWithTypesTx(): %w", err)
	}

	return txHash, err
}

// GetConfirmation returns confirmation number
func (e *Ethereum) GetConfirmation(ctx context.Context, hashTx string) (uint64, error) {
	txInfo, err := e.GetTransactionByHash(ctx, hashTx)
	if err != nil {
		return 0, err
	}
	if txInfo.BlockNumber == 0 {
		return 0, errors.New("block number can't retrieved")
	}
	currentBlockNum, err := e.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	confirmation := currentBlockNum.Int64() - txInfo.BlockNumber

	return uint64(confirmation), nil
}

// SupportsEIP1559 detects whether the connected Ethereum node supports EIP-1559 transactions.
//
// EIP-1559 was introduced in the London hard fork and includes a base fee mechanism.
// This method checks for EIP-1559 support by:
//  1. Returning true immediately if the client is Anvil (always supports EIP-1559)
//  2. Checking if the latest block contains a baseFeePerGas field (indicates post-London)
//
// Returns true if EIP-1559 is supported, false otherwise.
func (e *Ethereum) SupportsEIP1559(ctx context.Context) bool {
	// Anvil always supports EIP-1559
	if e.clientType == pkgeth.ClientVersionAnvil {
		return true
	}

	// Check if latest block has baseFeePerGas (indicates EIP-1559 support)
	currentBlockNum, err := e.BlockNumber(ctx)
	if err != nil {
		logger.Warn("failed to get block number for EIP-1559 detection", "error", err)
		return false
	}

	blockInfo, err := e.GetBlockByNumber(ctx, currentBlockNum.Uint64())
	if err != nil {
		logger.Warn("failed to get block info for EIP-1559 detection", "error", err)
		return false
	}

	return blockInfo.BaseFeePerGas != nil
}

// CreateRawTransactionEIP1559 creates an EIP-1559 transaction with dynamic fee pricing.
//
// EIP-1559 transactions use a two-tier gas pricing model:
//   - maxPriorityFeePerGas: Tip paid directly to miners
//   - maxFeePerGas: Maximum total fee willing to pay (base fee + priority fee)
//
// This method calculates fees as:
//   - maxPriorityFeePerGas = 2 Gwei (reasonable default for most networks)
//   - maxFeePerGas = (baseFeePerGas * 2) + maxPriorityFeePerGas
//
// The doubling of baseFeePerGas provides buffer for base fee increases between
// transaction submission and inclusion in a block.
//
// Note: EIP-1559 is only supported on networks that have completed the London hard fork.
// Use SupportsEIP1559() to check compatibility before calling this method.
func (e *Ethereum) CreateRawTransactionEIP1559(
	ctx context.Context, fromAddr, toAddr string, amount uint64, additionalNonce int,
) (*domainETH.RawTx, *apieth.TxCreateParams, error) {
	// validation check
	if pkgeth.ValidateAddr(fromAddr) != nil || pkgeth.ValidateAddr(toAddr) != nil {
		return nil, nil, errors.New("address validation error")
	}
	logger.Debug("eth.CreateRawTransactionEIP1559()",
		"fromAddr", fromAddr,
		"toAddr", toAddr,
		"amount", amount,
	)

	// Check EIP-1559 support
	if !e.SupportsEIP1559(ctx) {
		return nil, nil, errors.New("EIP-1559 transactions not supported by this network")
	}

	// Get balance
	balance, err := e.GetBalance(ctx, fromAddr, domainETH.QuantityTagPending)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}
	logger.Info("balance", "balance", balance.Int64())
	if balance.Uint64() == 0 {
		return nil, nil, errors.New("balance is needed to send eth")
	}

	// nonce
	nonce, err := e.getNonce(ctx, fromAddr, additionalNonce)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.getNonce(): %w", err)
	}

	// Get base fee from latest block
	currentBlockNum, err := e.BlockNumber(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.BlockNumber(): %w", err)
	}

	blockInfo, err := e.GetBlockByNumber(ctx, currentBlockNum.Uint64())
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.GetBlockByNumber(): %w", err)
	}

	if blockInfo.BaseFeePerGas == nil {
		return nil, nil, errors.New("baseFeePerGas not found in block (EIP-1559 not activated)")
	}

	// Calculate EIP-1559 fees
	// maxPriorityFeePerGas: suggested by the node with config fallback
	maxPriorityFeePerGas, err := e.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call eth.SuggestGasTipCap(): %w", err)
	}

	// maxFeePerGas: (baseFee * 2) + maxPriorityFee
	// The doubling provides buffer for baseFee increases between tx creation and inclusion
	maxFeePerGas := new(big.Int).Mul(blockInfo.BaseFeePerGas, big.NewInt(2))
	maxFeePerGas = maxFeePerGas.Add(maxFeePerGas, maxPriorityFeePerGas)

	logger.Info("EIP-1559 fees",
		"baseFeePerGas", blockInfo.BaseFeePerGas.Uint64(),
		"maxPriorityFeePerGas", maxPriorityFeePerGas.Uint64(),
		"maxFeePerGas", maxFeePerGas.Uint64(),
	)

	// Estimate gas
	tmpToAddrForEstimate := common.HexToAddress(toAddr)
	msg := &ethereum.CallMsg{
		From:  common.HexToAddress(fromAddr),
		To:    &tmpToAddrForEstimate,
		Gas:   0,
		Value: nil,
		Data:  nil,
	}
	estimatedGas, err := e.EstimateGas(ctx, msg)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call EstimateGas(): %w", err)
	}

	// Calculate transaction fee
	txFee := new(big.Int).Mul(maxFeePerGas, estimatedGas)

	// Calculate value to send
	newValue := new(big.Int)
	if amount == 0 {
		// receiver pays fee (deposit, transfer(pays all) action)
		newValue = newValue.Sub(balance, txFee)
	} else {
		// sender pays fee (payment, transfer(pays partially))
		newValue = new(big.Int).SetUint64(amount)
		if balance.Cmp(new(big.Int).Add(newValue, txFee)) == -1 {
			return nil, nil, fmt.Errorf(
				"balance `%d` is insufficient to send `%d` + fee `%d`",
				balance.Uint64(), newValue.Uint64(), txFee.Uint64())
		}
	}

	logger.Debug("EIP-1559 tx parameter",
		"estimatedGas", estimatedGas.Uint64(),
		"txFee", txFee.Uint64(),
		"newValue", newValue.Uint64())

	// create EIP-1559 transaction
	tmpToAddr := common.HexToAddress(toAddr)
	chainID := big.NewInt(int64(e.netID))

	// Note: Using estimatedGas for accuracy. For simple ETH transfers, this should equal GasLimit (21000).
	// Unlike legacy transactions which use the constant GasLimit, EIP-1559 uses the estimated value
	// to ensure sufficient gas for any transaction type.
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       estimatedGas.Uint64(),
		To:        &tmpToAddr,
		Value:     newValue,
		Data:      nil,
	})

	// Use the signer hash (EIP-1559 pre-image) not tx.Hash() for MPC-TSS signing.
	// tx.Hash() encodes the unsigned tx with zero r/s/v fields, producing a different hash
	// than what types.Sender recovers from. signer.Hash(tx) is the canonical pre-image.
	signer := types.LatestSignerForChainID(chainID)
	txHash := signer.Hash(tx).Hex()
	rawTxHex, err := pkgeth.EncodeTx(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call encodeTx(): %w", err)
	}

	// generate UUID to trace transaction
	uid, err := e.uuidHandler.GenerateV7()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	// create domain RawTx
	domainRawTx := &domainETH.RawTx{
		UUID:  uid.String(),
		From:  fromAddr,
		To:    toAddr,
		Value: *newValue,
		Nonce: nonce,
		TxHex: *rawTxHex,
		Hash:  txHash,
	}

	// create TxCreateParams DTO for use case layer
	txParams := &apieth.TxCreateParams{
		UUID:                 uid.String(),
		FromAddress:          fromAddr,
		ToAddress:            toAddr,
		Amount:               newValue.Uint64(),
		Fee:                  txFee.Uint64(),
		GasLimit:             uint32(estimatedGas.Uint64()),
		Nonce:                nonce,
		EthTxType:            2, // EIP-1559 transaction
		ChainID:              chainID.Uint64(),
		MaxFeePerGas:         maxFeePerGas.Uint64(),
		MaxPriorityFeePerGas: maxPriorityFeePerGas.Uint64(),
	}

	return domainRawTx, txParams, nil
}

// SignTxWithPrivateKey signs a raw transaction using a private key directly,
// without requiring a keystore or node connection (fully offline operation).
// Uses LatestSignerForChainID for forward compatibility with future transaction types.
func (*Ethereum) SignTxWithPrivateKey(
	rawTx *domainETH.RawTx, privKey *ecdsa.PrivateKey, chainID *big.Int,
) (*domainETH.RawTx, error) {
	return pkgeth.SignTxOffline(rawTx, privKey, chainID)
}
