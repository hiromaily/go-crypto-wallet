package watchsrv

import (
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// CreateDepositTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateDepositTx() (string, string, error) {
	sender := account.AccountTypeClient
	receiver := t.depositReceiver
	targetAction := action.ActionTypeDeposit
	t.logger.Debug("account",
		zap.String("sender", sender.String()),
		zap.String("receiver", receiver.String()),
	)

	userAmounts, err := t.getUserAmounts(sender)
	if err != nil {
		return "", "", err
	}
	if len(userAmounts) == 0 {
		t.logger.Info("no data")
		return "", "", nil
	}

	serializedTxs, txDetailItems, err := t.createDepositRawTransactions(sender, receiver, userAmounts)
	if err != nil {
		return "", "", err
	}
	if len(txDetailItems) == 0 {
		return "", "", nil
	}

	txID, err := t.updateDB(targetAction, txDetailItems, nil)
	if err != nil {
		return "", "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = t.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return "", generatedFileName, nil
}

func (t *TxCreate) getUserAmounts(sender account.AccountType) ([]xrp.UserAmount, error) {
	// get addresses for client account
	addrs, err := t.addrRepo.GetAll(sender)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	// addresses, err := t.eth.Accounts()

	// target addresses
	var userAmounts []xrp.UserAmount
	// address list for client
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		clientBalance, err := t.xrp.GetBalance(addr.WalletAddress)
		if err != nil {
			t.logger.Warn("fail to call t.xrp.GetAccountInfo()",
				zap.String("address", addr.WalletAddress),
			)
		} else {
			t.logger.Debug("account_info", zap.String("address", addr.WalletAddress), zap.Float64("balance", clientBalance))
			if clientBalance != 0 {
				userAmounts = append(userAmounts, xrp.UserAmount{Address: addr.WalletAddress, Amount: clientBalance})
			}
		}
	}
	return userAmounts, nil
}

func (t *TxCreate) createDepositRawTransactions(sender, receiver account.AccountType, userAmounts []xrp.UserAmount) ([]string, []*models.XRPDetailTX, error) {
	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(account.AccountTypeDeposit)")
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.XRPDetailTX, 0, len(userAmounts))

	var sequence uint64
	for _, val := range userAmounts {
		// call CreateRawTransaction
		instructions := &xrp.Instructions{
			MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		txJSON, rawTxString, err := t.xrp.CreateRawTransaction(val.Address, depositAddr.WalletAddress, 0, instructions)
		if err != nil {
			t.logger.Warn("fail to call xrp.CreateRawTransaction()", zap.Error(err))
			continue
		}
		t.logger.Debug("txJSON", zap.Any("txJSON", txJSON))
		grok.Value(txJSON)

		// sequence for next rawTransaction
		sequence = txJSON.Sequence + 1

		// generate UUID to trace transaction because unsignedTx is not unique
		uid := uuid.NewV4().String()

		serializedTxs = append(serializedTxs, fmt.Sprintf("%s,%s", uid, rawTxString))

		// create insert data for　eth_detail_tx
		txDetailItem := &models.XRPDetailTX{
			UUID:               uid,
			CurrentTXType:      tx.TxTypeUnsigned.Int8(),
			SenderAccount:      sender.String(),
			SenderAddress:      val.Address,
			ReceiverAccount:    receiver.String(),
			ReceiverAddress:    depositAddr.WalletAddress,
			Amount:             txJSON.Amount,
			XRPTXType:          txJSON.TransactionType,
			Fee:                txJSON.Fee,
			Flags:              txJSON.Flags,
			LastLedgerSequence: txJSON.LastLedgerSequence,
			Sequence:           txJSON.Sequence,
			// SigningPubkey:      txJSON.SigningPubKey,
			// TXNSignature:       txJSON.TxnSignature,
			// Hash:               txJSON.Hash,
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	return serializedTxs, txDetailItems, nil
}
