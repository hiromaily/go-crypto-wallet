package xrp

import (
	"context"
	"fmt"

	"github.com/bookerzzz/grok"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// CreateDepositTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateDepositTx() (string, string, error) {
	sender := domainAccount.AccountTypeClient
	receiver := t.depositReceiver
	targetAction := domainTx.ActionTypeDeposit
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	userAmounts, err := t.getUserAmounts(sender)
	if err != nil {
		return "", "", err
	}
	if len(userAmounts) == 0 {
		logger.Info("no data")
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
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return "", generatedFileName, nil
}

func (t *TxCreate) getUserAmounts(sender domainAccount.AccountType) ([]xrp.UserAmount, error) {
	// get addresses for client account
	addrs, err := t.addrRepo.GetAll(sender)
	if err != nil {
		return nil, fmt.Errorf("fail to call addrRepo.GetAll(domainAccount.AccountTypeClient): %w", err)
	}
	// addresses, err := t.eth.Accounts()

	// target addresses
	var userAmounts []xrp.UserAmount
	// address list for client
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		var clientBalance float64
		clientBalance, err = t.xrp.GetBalance(context.TODO(), addr.WalletAddress)
		if err != nil {
			logger.Warn("fail to call t.xrp.GetAccountInfo()",
				"address", addr.WalletAddress,
			)
		} else {
			logger.Debug("account_info",
				"address", addr.WalletAddress, "balance", clientBalance)
			if clientBalance != 0 {
				userAmounts = append(userAmounts, xrp.UserAmount{Address: addr.WalletAddress, Amount: clientBalance})
			}
		}
	}
	return userAmounts, nil
}

func (t *TxCreate) createDepositRawTransactions(
	sender, receiver domainAccount.AccountType, userAmounts []xrp.UserAmount,
) ([]string, []*models.XRPDetailTX, error) {
	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"fail to call addrRepo.GetOneUnAllocated(domainAccount.AccountTypeDeposit): %w", err,
		)
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
		var txJSON *xrp.TxInput
		var rawTxString string
		txJSON, rawTxString, err = t.xrp.CreateRawTransaction(
			context.TODO(), val.Address, depositAddr.WalletAddress, 0, instructions)
		if err != nil {
			logger.Warn("fail to call xrp.CreateRawTransaction()", "error", err)
			continue
		}
		logger.Debug("txJSON", "txJSON", txJSON)
		grok.Value(txJSON)

		// sequence for next rawTransaction
		sequence = txJSON.Sequence + 1

		// generate UUID to trace transaction because unsignedTx is not unique
		uid, err := t.uuidHandler.GenerateV7()
		if err != nil {
			return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
		}

		serializedTxs = append(serializedTxs, fmt.Sprintf("%s,%s", uid, rawTxString))

		// create insert data for　eth_detail_tx
		txDetailItem := &models.XRPDetailTX{
			UUID:               uid.String(),
			CurrentTXType:      domainTx.TxTypeUnsigned.Int8(),
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
