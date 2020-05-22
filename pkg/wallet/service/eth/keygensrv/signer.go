package keygensrv

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

// Sign type
type Sign struct {
	eth        ethgrp.Ethereumer
	logger     *zap.Logger
	txFileRepo tx.FileRepositorier
	wtype      wallet.WalletType
}

// NewSign returns sign object
func NewSign(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType) *Sign {

	return &Sign{
		eth:        eth,
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

	var senderAccount account.AccountType

	// get hex tx from file
	data, err := s.txFileRepo.ReadFileSlice(filePath)
	if err != nil {
		return "", false, "", errors.Wrap(err, "fail to call txFileRepo.ReadFileSlice()")
	}
	if len(data) > 1 {
		senderAccount = account.AccountType(data[0])
	} else {
		return "", false, "", errors.New("file is invalid")
	}
	serializedTxs := data[1:]

	txHexs := make([]string, 0, len(serializedTxs))
	for _, serializedTx := range serializedTxs {
		var rawTx eth.RawTx
		if err = serial.DecodeFromString(serializedTx, &rawTx); err != nil {
			return "", false, "", errors.Wrap(err, "fail to call serial.DecodeFromString()")
		}
		// sign
		signedRawTx, err := s.eth.SignOnRawTransaction(&rawTx, eth.Password, senderAccount)
		if err != nil {
			return "", false, "", errors.Wrap(err, "fail to call eth.SignOnRawTransaction()")
		}
		txHexs = append(txHexs, fmt.Sprintf("%s,%s", rawTx.UUID, signedRawTx.TxHex))
	}

	// write file
	path := s.txFileRepo.CreateFilePath(actionType, tx.TxTypeSigned, txID, signedCount)
	generatedFileName, err := s.txFileRepo.WriteFileSlice(path, txHexs)
	if err != nil {
		return "", false, "", errors.Wrap(err, "fail to call txFileRepo.WriteFileSlice()")
	}

	//return hexTx, isSigned, generatedFileName, nil
	return "", true, generatedFileName, nil
}
