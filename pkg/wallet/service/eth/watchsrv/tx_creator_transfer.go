package watchsrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
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

	// check sernder's balance
	senderAddr, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(sender)")
	}
	senderBalnce, err := t.eth.GetBalance(senderAddr.WalletAddress, eth.QuantityTagLatest)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call eth.GetBalance(sender)")
	}

	if senderBalnce.Uint64() == 0 {
		return "", "", errors.New("sender has no balance")
	}

	requiredValue := t.eth.FromFloatEther(floatValue)
	if floatValue != 0 && (senderBalnce.Uint64() <= requiredValue.Uint64()) {
		return "", "", errors.New("sender balance is insufficient to send")
	}
	t.logger.Debug("amount",
		zap.Float64("floatValue(Ether)", floatValue),
		zap.Uint64("requiredValue(Ether)", requiredValue.Uint64()),
		zap.Uint64("senderBalnce", senderBalnce.Uint64()),
	)

	// get receiver address
	receiverAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(receiver)")
	}

	//txDetailItems := make([]*models.EthDetailTX, 0, 1)

	// call CreateRawTransaction
	rawTx, txDetailItem, err := t.eth.CreateRawTransaction(senderAddr.WalletAddress, receiverAddr.WalletAddress, requiredValue.Uint64(), 0)
	if err != nil {
		return "", "", errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), sender address: %s", senderAddr.WalletAddress)
	}

	rawTxHex := rawTx.TxHex
	t.logger.Debug("rawTxHex", zap.String("rawTxHex", rawTxHex))

	serializedTx, err := serial.EncodeToString(rawTx)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString(rawTx)")
	}
	serializedTxs := []string{serializedTx}

	// create insert data for　eth_detail_tx
	txDetailItem.SenderAccount = sender.String()
	txDetailItem.ReceiverAccount = receiver.String()
	txDetailItems := []*models.EthDetailTX{txDetailItem}

	return t.afterTxCreation(targetAction, sender, serializedTxs, txDetailItems, nil)
}
