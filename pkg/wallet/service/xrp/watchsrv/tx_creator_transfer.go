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
	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
// - sender pays fee,
// - any internal account should have only one address in Ethereum because no utxo
func (t *TxCreate) CreateTransferTx(sender, receiver account.AccountType, floatValue float64) (string, string, error) {
	targetAction := action.ActionTypeTransfer

	// validation account
	if receiver == account.AccountTypeClient || receiver == account.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	// check sender's balance
	senderAddr, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(sender)")
	}
	senderBalance, err := t.xrp.GetBalance(senderAddr.WalletAddress)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call xrp.GetAccountInfo()")
	}
	if senderBalance <= 20 {
		return "", "", errors.New("sender balance is insufficient to send")
	}
	if floatValue != 0 && senderBalance <= floatValue {
		return "", "", errors.New("sender balance is insufficient to send")
	}

	t.logger.Debug("amount",
		zap.Float64("floatValue", floatValue),
		zap.Float64("senderBalance", senderBalance),
	)

	// get receiver address
	receiverAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(receiver)")
	}

	// call CreateRawTransaction
	instructions := &pb.Instructions{
		MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
	}
	txJSON, rawTxString, err := t.xrp.CreateRawTransaction(senderAddr.WalletAddress, receiverAddr.WalletAddress, floatValue, instructions)
	if err != nil {
		return "", "", errors.Wrapf(err, "fail to call xrp.CreateRawTransaction(), sender address: %s", senderAddr.WalletAddress)
	}
	t.logger.Debug("txJSON", zap.Any("txJSON", txJSON))
	grok.Value(txJSON)

	// generate UUID to trace transaction because unsignedTx is not unique
	uid := uuid.NewV4().String()

	serializedTxs := []string{fmt.Sprintf("%s,%s", uid, rawTxString)}

	// create insert data for　eth_detail_tx
	txDetailItem := &models.XRPDetailTX{
		UUID:               uid,
		CurrentTXType:      tx.TxTypeUnsigned.Int8(),
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
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return "", generatedFileName, nil
}
