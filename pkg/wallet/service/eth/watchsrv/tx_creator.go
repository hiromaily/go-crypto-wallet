package watchsrv

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/repository/watchrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// TxCreate type
type TxCreate struct {
	eth        ethgrp.Ethereumer
	logger     *zap.Logger
	dbConn     *sql.DB
	addrRepo   watchrepo.AddressRepositorier
	txRepo     watchrepo.TxRepositorier
	payReqRepo watchrepo.PaymentRequestRepositorier
	txFileRepo tx.FileRepositorier
	wtype      wallet.WalletType
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType) *TxCreate {

	return &TxCreate{
		eth:        eth,
		logger:     logger,
		dbConn:     dbConn,
		addrRepo:   addrRepo,
		txRepo:     txRepo,
		payReqRepo: payReqRepo,
		txFileRepo: txFileRepo,
		wtype:      wtype,
	}
}

// generateHexFile generate file for hex and encoded previous addresses
func (t *TxCreate) generateHexFile(actionType action.ActionType, bTxs [][]byte) (string, error) {
	var (
		generatedFileName string
		//err               error
	)

	//savedata := hex
	//if encodedAddrsPrevs != "" {
	//	savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
	//}
	//
	//// create file
	//path := t.txFileRepo.CreateFilePath(actionType, tx.TxTypeUnsigned, id, 0)
	//generatedFileName, err = t.txFileRepo.WriteFile(path, savedata)
	//if err != nil {
	//	return "", errors.Wrap(err, "fail to call txFileRepo.WriteFile()")
	//}

	return generatedFileName, nil
}
