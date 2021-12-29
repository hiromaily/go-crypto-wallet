package watchsrv

import (
	"math/big"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
// Note
// - to avoid complex logic to create raw transaction
// - only one address of sender should afford to send coin to all payment request users.
func (t *TxCreate) CreatePaymentTx() (string, string, error) {
	sender := t.paymentSender
	receiver := account.AccountTypeAnonymous
	targetAction := action.ActionTypePayment
	t.logger.Debug("account",
		zap.String("sender", sender.String()),
		zap.String("receiver", receiver.String()),
	)

	userPayments, paymentRequestIds, addrItem, err := t.getUserPayments(sender)
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		return "", "", nil
	}

	serializedTxs, txDetailItems, err := t.createPaymentRawTransactions(sender, receiver, userPayments, addrItem)
	if err != nil {
		return "", "", err
	}

	txID, err := t.updateDB(targetAction, txDetailItems, paymentRequestIds)
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

func (t *TxCreate) getUserPayments(sender account.AccountType) ([]UserPayment, []int64, *models.Address, error) {
	// get payment data from payment_request
	userPayments, totalAmount, paymentRequestIds, err := t.createUserPayment()
	if err != nil {
		return nil, nil, nil, err
	}
	if len(userPayments) == 0 {
		t.logger.Debug("no data in userPayments")
		// no data
		return nil, nil, nil, nil
	}

	// check sernder's total balance
	// GetOneUnAllocated
	addrItem, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	senderBalance, err := t.eth.GetBalance(addrItem.WalletAddress, eth.QuantityTagPending)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "fail to call eth.GetBalance()")
	}

	if senderBalance.Uint64() <= totalAmount.Uint64() {
		return nil, nil, nil, errors.New("sender balance is insufficient to send")
	}
	return userPayments, paymentRequestIds, addrItem, nil
}

// UserPayment user's payment address and amount
type UserPayment struct {
	senderAddr   string   // sender address for just chacking
	receiverAddr string   // receiver address
	floatAmount  float64  // float amount (Ether)
	amount       *big.Int // amount (Wei)
}

// createUserPayment get payment data from payment_request table
func (t *TxCreate) createUserPayment() ([]UserPayment, *big.Int, []int64, error) {
	// get payment_request
	paymentRequests, err := t.payReqRepo.GetAll()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "fail to call repo.GetPaymentRequestAll()")
	}
	if len(paymentRequests) == 0 {
		t.logger.Debug("no data in payment_request")
		return nil, nil, nil, nil
	}

	userPayments := make([]UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))
	totalAmount := new(big.Int)

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.SenderAddress
		userPayments[idx].receiverAddr = val.ReceiverAddress
		amt, err := strconv.ParseFloat(val.Amount.String(), 64)
		if err != nil {
			// fatal error because table includes invalid data
			t.logger.Error("payment_request table includes invalid amount field")
			return nil, nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if err = t.eth.ValidateAddr(userPayments[idx].receiverAddr); err != nil {
			// fatal error
			t.logger.Error("fail to call ValidationAddr",
				zap.String("address", userPayments[idx].receiverAddr),
				zap.Error(err),
			)
			return nil, nil, nil, errors.Wrap(err, "fail to call eth.ValidationAddr()")
		}

		// amount
		userPayments[idx].amount = t.eth.FloatToBigInt(userPayments[idx].floatAmount)
		totalAmount = new(big.Int).Add(totalAmount, userPayments[idx].amount)
	}

	return userPayments, totalAmount, paymentRequestIds, nil
}

func (t *TxCreate) createPaymentRawTransactions(sender, receiver account.AccountType, userPayments []UserPayment, addrItem *models.Address) ([]string, []*models.EthDetailTX, error) {
	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userPayments))
	additionalNonce := 0
	for _, userPayment := range userPayments {
		// call CreateRawTransaction
		rawTx, txDetailItem, err := t.eth.CreateRawTransaction(addrItem.WalletAddress, userPayment.receiverAddr, userPayment.amount.Uint64(), additionalNonce)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "fail to call addrRepo.CreateRawTransaction(), sender address: %s", addrItem.WalletAddress)
		}
		additionalNonce++

		rawTxHex := rawTx.TxHex
		t.logger.Debug("rawTxHex", zap.String("rawTxHex", rawTxHex))
		// TODO: `rawTxHex` should be used to trace progress to update database

		serializedTx, err := serial.EncodeToString(rawTx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "fail to call serial.EncodeToString(rawTx)")
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = sender.String()
		txDetailItem.ReceiverAccount = receiver.String()
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems, nil
}
