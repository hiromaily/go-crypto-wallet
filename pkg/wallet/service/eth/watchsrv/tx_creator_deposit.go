package watchsrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// UserAmount user address and amount
type UserAmount struct {
	address string
	amount  uint64
}

// CreateDepositTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	targetAction := action.ActionTypeDeposit

	//1. get addresses for client account
	addrs, err := t.addrRepo.GetAll(account.AccountTypeClient)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	//addresses, err := t.eth.Accounts()

	// target addresses
	var userAmounts []UserAmount

	// address list for client
	for _, addr := range addrs {
		//TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		balance, err := t.eth.GetBalance(addr.WalletAddress, eth.QuantityTagLatest)
		if err != nil {
			t.logger.Warn("fail to call t.eth.GetBalance()",
				zap.String("address", addr.WalletAddress),
				zap.Error(err),
			)
		} else {
			if balance.Uint64() != 0 {
				userAmounts = append(userAmounts, UserAmount{address: addr.WalletAddress, amount: balance.Uint64()})
			}
		}
	}

	if len(userAmounts) == 0 {
		t.logger.Info("no data")
		return "", "", nil
	}

	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(account.AccountTypeDeposit)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(account.AccountTypeDeposit)")
	}

	// create raw transaction each address
	bTxs := make([][]byte, 0, len(userAmounts))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userAmounts))
	for _, val := range userAmounts {
		// call CreateRawTransaction
		rawTxHex, bTx, txDetailItem, err := t.eth.CreateRawTransaction(val.address, depositAddr.WalletAddress, val.amount)
		if err != nil {
			return "", "", errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), address: %s", val.address)
		}
		t.logger.Debug("rawTxHex", zap.String("rawTxHex", rawTxHex))
		bTxs = append(bTxs, bTx)

		txDetailItem.SenderAccount = account.AccountTypeClient.String()
		txDetailItem.ReceiverAccount = account.AccountTypeDeposit.String()

		// create insert data for　eth_detail_tx
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	// start transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return "", "", errors.Wrap(err, "fail to start transaction")
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// Insert eth_tx
	txID, err := t.txRepo.InsertUnsignedTx(action.ActionTypeDeposit)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call txRepo.InsertUnsignedTx()")
	}
	// TODO: Insert to eth_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = t.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return "", "", errors.Wrap(err, "fail to call txDetailRepo.InsertBulk()")
	}

	// save transaction result to file
	var generatedFileName string
	if len(bTxs) != 0 {
		//TODO: implement generateHexFile()
		generatedFileName, err = t.generateHexFile(targetAction, bTxs)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return "", generatedFileName, nil
}
