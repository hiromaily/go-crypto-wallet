package watchsrv

import (
	"context"
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
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
	senderBalnce, err := t.eth.GetBalance(context.TODO(), senderAddr.WalletAddress, eth.QuantityTagLatest)
	if err != nil {
		return "", "", fmt.Errorf("fail to call eth.GetBalance(sender): %w", err)
	}

	if senderBalnce.Uint64() == 0 {
		return "", "", errors.New("sender has no balance")
	}

	requiredValue := t.eth.FloatToBigInt(floatValue)
	if floatValue != 0 && (senderBalnce.Uint64() <= requiredValue.Uint64()) {
		return "", "", errors.New("sender balance is insufficient to send")
	}
	logger.Debug("amount",
		"floatValue(Ether)", floatValue,
		"requiredValue(Ether)", requiredValue.Uint64(),
		"senderBalance", senderBalnce.Uint64(),
	)

	// get receiver address
	receiverAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(receiver): %w", err)
	}

	// call CreateRawTransaction
	rawTx, txDetailItem, err := t.eth.CreateRawTransaction(context.TODO(),
		senderAddr.WalletAddress, receiverAddr.WalletAddress, requiredValue.Uint64(), 0)
	if err != nil {
		return "", "", fmt.Errorf(
			"fail to call eth.CreateRawTransaction(), sender address: %s: %w",
			senderAddr.WalletAddress, err)
	}

	rawTxHex := rawTx.TxHex
	logger.Debug("rawTxHex", "rawTxHex", rawTxHex)

	serializedTx, err := serial.EncodeToString(rawTx)
	if err != nil {
		return "", "", fmt.Errorf("fail to call serial.EncodeToString(rawTx): %w", err)
	}
	serializedTxs := []string{serializedTx}

	// create insert data for　eth_detail_tx
	txDetailItem.SenderAccount = sender.String()
	txDetailItem.ReceiverAccount = receiver.String()
	txDetailItems := []*models.EthDetailTX{txDetailItem}

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
