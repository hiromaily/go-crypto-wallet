package eth

import (
	"context"
	"fmt"
	"math/big"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/ethtx"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
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
	logger.Debug("update result",
		"txID", txID,
		"error", err,
	)
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

func (t *TxCreate) getUserAmounts(sender domainAccount.AccountType) ([]eth.UserAmount, error) {
	// get addresses for client account
	addrs, err := t.addrRepo.GetAll(sender)
	if err != nil {
		return nil, fmt.Errorf("fail to call addrRepo.GetAll(domainAccount.AccountTypeClient): %w", err)
	}

	// target addresses
	var userAmounts []eth.UserAmount

	// address list for client
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		var balance *big.Int
		balance, err = t.eth.GetBalance(context.TODO(), addr.WalletAddress, eth.QuantityTagLatest)
		if err != nil {
			logger.Warn("fail to call .GetBalance()",
				"address", addr.WalletAddress,
				"error", err,
			)
		} else if balance.Uint64() != 0 {
			userAmounts = append(userAmounts, eth.UserAmount{Address: addr.WalletAddress, Amount: balance.Uint64()})
		}
	}

	return userAmounts, nil
}

func (t *TxCreate) createDepositRawTransactions(
	sender, receiver domainAccount.AccountType, userAmounts []eth.UserAmount,
) ([]string, []*models.EthDetailTX, error) {
	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"fail to call addrRepo.GetOneUnAllocated(domainAccount.AccountTypeDeposit): %w", err,
		)
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userAmounts))
	// additionalNonce := 0
	for _, val := range userAmounts {
		// call CreateRawTransaction
		var rawTx *ethtx.RawTx
		var txDetailItem *models.EthDetailTX
		rawTx, txDetailItem, err = t.eth.CreateRawTransaction(
			context.TODO(), val.Address, depositAddr.WalletAddress, 0, 0)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"fail to call addrRepo.CreateRawTransaction(), sender address: %s: %w",
				val.Address, err)
		}
		// additionalNonce++

		rawTxHex := rawTx.TxHex
		logger.Debug("rawTxHex", "rawTxHex", rawTxHex)

		var serializedTx string
		serializedTx, err = serial.EncodeToString(rawTx)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to call serial.EncodeToString(rawTx): %w", err)
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = sender.String()
		txDetailItem.ReceiverAccount = receiver.String()
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems, nil
}
