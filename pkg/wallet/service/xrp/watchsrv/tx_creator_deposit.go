package watchsrv

import (
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
	receiver := account.AccountTypeDeposit
	targetAction := action.ActionTypeDeposit

	//1. get addresses for client account
	addrs, err := t.addrRepo.GetAll(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	//addresses, err := t.eth.Accounts()

	// target addresses
	var userAmounts []xrp.UserAmount

	// address list for client
	for _, addr := range addrs {
		//TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		accountInfo, err := t.xrp.GetAccountInfo(addr.WalletAddress)
		if err != nil {
			t.logger.Warn("fail to call t.xrp.GetAccountInfo()",
				zap.String("address", addr.WalletAddress),
				zap.Error(err),
			)
		} else {
			amt := xrp.ToFloat64(accountInfo.XrpBalance)
			if amt != 0 {
				userAmounts = append(userAmounts, xrp.UserAmount{Address: addr.WalletAddress, Amount: amt})
			}
		}
	}

	if len(userAmounts) == 0 {
		t.logger.Info("no data")
		return "", "", nil
	}

	// get address for deposit account
	depositAddr, err := t.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetOneUnAllocated(account.AccountTypeDeposit)")
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.XRPDetailTX, 0, len(userAmounts))
	for _, val := range userAmounts {
		// call CreateRawTransaction
		txJSON, rawTxString, err := t.xrp.CreateRawTransaction(val.Address, depositAddr.WalletAddress, val.Amount)
		if err != nil {
			t.logger.Warn("fail to call xrp.CreateRawTransaction()", zap.Error(err))
			//return "", "", errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), sender address: %s", val.Address)
			continue
		}

		t.logger.Debug("txJSON", zap.Any("txJSON", txJSON))

		//serializedTx, err := serial.EncodeToString(txJSON)
		//if err != nil {
		//	return "", "", errors.Wrap(err, "fail to call serial.EncodeToString(txJSON)")
		//}
		serializedTxs = append(serializedTxs, rawTxString)

		// generate UUID to trace transaction because unsignedTx is not unique
		uid := uuid.NewV4().String()

		// create insert data for　eth_detail_tx
		txDetailItem := &models.XRPDetailTX{
			UUID:               uid,
			CurrentTXType:      tx.TxTypeUnsigned.Int8(),
			SenderAccount:      sender.String(),
			SenderAddress:      val.Address,
			ReceiverAccount:    receiver.String(),
			ReceiverAddress:    depositAddr.WalletAddress,
			Amount:             txJSON.Amount, //TODO:compare to fmt.Sprint(val.Amount)
			XRPTXType:          txJSON.TransactionType,
			Fee:                txJSON.Fee,
			Flags:              txJSON.Flags,
			LastLedgerSequence: txJSON.LastLedgerSequence,
			Sequence:           txJSON.Sequence,
			SigningPubkey:      txJSON.SigningPubKey,
			TXNSignature:       txJSON.TxnSignature,
			Hash:               txJSON.Hash,
		}

		//ID                    int64     `boil:"id" json:"id" toml:"id" yaml:"id"`
		//TXID                  int64     `boil:"tx_id" json:"tx_id" toml:"tx_id" yaml:"tx_id"`

		//EarliestLedgerVersion uint64    `boil:"earliest_ledger_version" json:"earliest_ledger_version" toml:"earliest_ledger_version" yaml:"earliest_ledger_version"`
		//SignedTXID            string    `boil:"signed_tx_id" json:"signed_tx_id" toml:"signed_tx_id" yaml:"signed_tx_id"`
		//SignedTXBlob          string    `boil:"signed_tx_blob" json:"signed_tx_blob" toml:"signed_tx_blob" yaml:"signed_tx_blob"`
		//SentTXBlob            string    `boil:"sent_tx_blob" json:"sent_tx_blob" toml:"sent_tx_blob" yaml:"sent_tx_blob"`
		//UnsignedUpdatedAt     null.Time `boil:"unsigned_updated_at" json:"unsigned_updated_at,omitempty" toml:"unsigned_updated_at" yaml:"unsigned_updated_at,omitempty"`
		//SentUpdatedAt         null.Time `boil:"sent_updated_at" json:"sent_updated_at,omitempty" toml:"sent_updated_at" yaml:"sent_updated_at,omitempty"`

		txDetailItems = append(txDetailItems, txDetailItem)
	}
	if len(txDetailItems) == 0 {
		//TODO: what should be returned?
		return "", "", nil
	}

	return t.afterTxCreation(targetAction, sender, serializedTxs, txDetailItems, nil)
}
