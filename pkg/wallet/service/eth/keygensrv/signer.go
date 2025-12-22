package keygensrv

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/ethtx"
)

// Sign type
type Sign struct {
	eth        ethgrp.Ethereumer
	logger     logger.Logger
	txFileRepo tx.FileRepositorier
	wtype      wallet.WalletType
}

// NewSign returns sign object
func NewSign(
	ethAPI ethgrp.Ethereumer,
	logger logger.Logger,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType,
) *Sign {
	return &Sign{
		eth:        ethAPI,
		logger:     logger,
		txFileRepo: txFileRepo,
		wtype:      wtype,
	}
}

// SignTx sign on tx in csv file
// - multisig equivalent functionality is not implemented yet in ETH
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	// get tx_deposit_id from tx file name
	actionType, _, txID, signedCount, err := s.txFileRepo.ValidateFilePath(filePath, tx.TxTypeUnsigned)
	if err != nil {
		return "", false, "", err
	}

	// get hex tx from file
	data, err := s.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.ReadFileSlice(): %w", err)
	}
	if len(data) <= 1 {
		return "", false, "", errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		var rawTx ethtx.RawTx
		if err = serial.DecodeFromString(serializedTx, &rawTx); err != nil {
			return "", false, "", fmt.Errorf("fail to call serial.DecodeFromString(): %w", err)
		}
		// sign
		var signedRawTx *ethtx.RawTx
		signedRawTx, err = s.eth.SignOnRawTransaction(&rawTx, eth.Password)
		if err != nil {
			return "", false, "", fmt.Errorf("fail to call eth.SignOnRawTransaction(): %w", err)
		}
		txHexs = append(txHexs, fmt.Sprintf("%s,%s", rawTx.UUID, signedRawTx.TxHex))
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, tx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := s.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return "", false, "", fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	// return hexTx, isSigned, generatedFileName, nil
	return "", true, generatedFileName, nil
}
