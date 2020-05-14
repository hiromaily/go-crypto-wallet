package watchsrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// CreateDepositTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	targetAction := action.ActionTypeDeposit
	senderAccount := account.AccountTypeClient
	receiverAccount := account.AccountTypeDeposit

	//1. get addresses for client account
	addrs, err := t.addrRepo.GetAll(senderAccount)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	//addresses, err := t.eth.Accounts()

	// target addresses
	var userAmounts []eth.UserAmount

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
				userAmounts = append(userAmounts, eth.UserAmount{Address: addr.WalletAddress, Amount: balance.Uint64()})
			}
		}
	}

	if len(userAmounts) == 0 {
		t.logger.Info("no data")
		return "", "", nil
	}

	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(receiverAccount)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(account.AccountTypeDeposit)")
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userAmounts))
	for _, val := range userAmounts {
		// call CreateRawTransaction
		rawTx, txDetailItem, err := t.eth.CreateRawTransaction(val.Address, depositAddr.WalletAddress, val.Amount)
		if err != nil {
			return "", "", errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), sender address: %s", val.Address)
		}

		rawTxHex := rawTx.TxHex
		t.logger.Debug("rawTxHex", zap.String("rawTxHex", rawTxHex))
		//TODO: `rawTxHex` should be used to trace progress to update database

		serializedTx, err := serial.EncodeToString(rawTx)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call serial.EncodeToString(rawTx)")
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = senderAccount.String()
		txDetailItem.ReceiverAccount = receiverAccount.String()
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	return t.afterTxCreation(targetAction, senderAccount, serializedTxs, txDetailItems)
}
