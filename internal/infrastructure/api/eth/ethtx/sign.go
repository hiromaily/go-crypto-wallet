package ethtx

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
)

// SignTxOffline signs a raw transaction using a private key directly, without requiring
// a keystore or Ethereum node connection. This enables true air-gapped offline signing.
//
// Uses types.LatestSignerForChainID for forward compatibility with future transaction
// types beyond EIP-1559 (London).
func SignTxOffline(
	rawTx *domainETH.RawTx,
	privKey *ecdsa.PrivateKey,
	chainID *big.Int,
) (*domainETH.RawTx, error) {
	if rawTx == nil {
		return nil, errors.New("rawTx must not be nil")
	}
	if privKey == nil {
		return nil, errors.New("private key must not be nil")
	}
	if chainID == nil || chainID.Sign() <= 0 {
		return nil, errors.New("chainID must be a positive integer")
	}

	tx, err := DecodeTx(rawTx.TxHex)
	if err != nil {
		return nil, fmt.Errorf("fail to decode raw transaction: %w", err)
	}

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, fmt.Errorf("fail to sign transaction: %w", err)
	}

	fromAddr, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, fmt.Errorf("fail to recover sender address: %w", err)
	}

	encodedTx, err := EncodeTx(signedTx)
	if err != nil {
		return nil, fmt.Errorf("fail to encode signed transaction: %w", err)
	}

	// To() is nil for contract creation transactions; preserve empty string in that case.
	toAddr := ""
	if signedTx.To() != nil {
		toAddr = signedTx.To().Hex()
	}

	return &domainETH.RawTx{
		UUID:  rawTx.UUID,
		From:  fromAddr.Hex(),
		To:    toAddr,
		Value: *signedTx.Value(),
		Nonce: signedTx.Nonce(),
		TxHex: *encodedTx,
		Hash:  signedTx.Hash().Hex(),
	}, nil
}
