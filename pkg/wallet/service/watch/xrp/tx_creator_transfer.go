package xrp

import (
	"context"
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
// - sender pays fee,
// - any internal account should have only one address in Ethereum because no utxo
func (t *TxCreate) CreateTransferTx(
	sender, receiver domainAccount.AccountType, floatValue float64,
) (string, string, error) {
	targetAction := domainTx.ActionTypeTransfer

	// validation account
	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	// check sender's balance
	senderAddr, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(sender): %w", err)
	}
	senderBalance, err := t.xrp.GetBalance(context.TODO(), senderAddr.WalletAddress)
	if err != nil {
		return "", "", fmt.Errorf("fail to call xrp.GetAccountInfo(): %w", err)
	}
	if senderBalance <= 20 {
		return "", "", errors.New("sender balance is insufficient to send")
	}
	if floatValue != 0 && senderBalance <= floatValue {
		return "", "", errors.New("sender balance is insufficient to send")
	}

	logger.Debug("amount",
		"floatValue", floatValue,
		"senderBalance", senderBalance,
	)

	// get receiver address
	receiverAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(receiver): %w", err)
	}

	// call CreateRawTransaction
	instructions := &xrp.Instructions{
		MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
	}
	txJSON, rawTxString, err := t.xrp.CreateRawTransaction(
		context.TODO(), senderAddr.WalletAddress, receiverAddr.WalletAddress, floatValue, instructions)
	if err != nil {
		return "", "", fmt.Errorf(
			"fail to call xrp.CreateRawTransaction(), sender address: %s: %w",
			senderAddr.WalletAddress, err)
	}
	logger.Debug("txJSON", "txJSON", txJSON)
	grok.Value(txJSON)

	// generate UUID to trace transaction because unsignedTx is not unique
	uid, err := t.uuidHandler.GenerateV7()
	if err != nil {
		return "", "", fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	serializedTxs := []string{fmt.Sprintf("%s,%s", uid, rawTxString)}

	// create insert data for　eth_detail_tx
	txDetailItem := &models.XRPDetailTX{
		UUID:               uid.String(),
		CurrentTXType:      domainTx.TxTypeUnsigned.Int8(),
		SenderAccount:      sender.String(),
		SenderAddress:      senderAddr.WalletAddress,
		ReceiverAccount:    receiver.String(),
		ReceiverAddress:    receiverAddr.WalletAddress,
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
	txDetailItems := []*models.XRPDetailTX{txDetailItem}

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
