package watchsrv

import (
	"fmt"
	"strconv"

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

	// get payment data from payment_request
	userPayments, totalAmount, paymentRequestIds, err := t.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		t.logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// check sender's total balance
	//GetOneUnAllocated
	addrItem, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call addrRepo.GetAll(account.AccountTypeClient)")
	}
	senderBalance, err := t.xrp.GetBalance(addrItem.WalletAddress)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call xrp.GetAccountInfo()")
	}

	if senderBalance <= totalAmount {
		return "", "", errors.New("sender balance is insufficient to send")
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*models.XRPDetailTX, 0, len(userPayments))
	var sequence uint64
	for _, userPayment := range userPayments {
		// call CreateRawTransaction
		instructions := &pb.Instructions{
			MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		txJSON, rawTxString, err := t.xrp.CreateRawTransaction(addrItem.WalletAddress, userPayment.receiverAddr, userPayment.floatAmount, instructions)
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
			SenderAddress:      addrItem.WalletAddress,
			ReceiverAccount:    receiver.String(),
			ReceiverAddress:    userPayment.receiverAddr,
			Amount:             txJSON.Amount,
			XRPTXType:          txJSON.TransactionType,
			Fee:                txJSON.Fee,
			Flags:              txJSON.Flags,
			LastLedgerSequence: txJSON.LastLedgerSequence,
			Sequence:           txJSON.Sequence,
			//SigningPubkey:      txJSON.SigningPubKey,
			//TXNSignature:       txJSON.TxnSignature,
			//Hash:               txJSON.Hash,
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	return t.afterTxCreation(targetAction, sender, serializedTxs, txDetailItems, paymentRequestIds)
}

// UserPayment user's payment address and amount
type UserPayment struct {
	senderAddr   string  // sender address for just chacking
	receiverAddr string  // receiver address
	floatAmount  float64 // float amount (Ether)
}

// createUserPayment get payment data from payment_request table
func (t *TxCreate) createUserPayment() ([]UserPayment, float64, []int64, error) {
	// get payment_request
	paymentRequests, err := t.payReqRepo.GetAll()
	if err != nil {
		return nil, 0, nil, errors.Wrap(err, "fail to call repo.GetPaymentRequestAll()")
	}
	if len(paymentRequests) == 0 {
		t.logger.Debug("no data in payment_request")
		return nil, 0, nil, nil
	}

	userPayments := make([]UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))
	var totalAmount float64

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.SenderAddress
		userPayments[idx].receiverAddr = val.ReceiverAddress
		amt, err := strconv.ParseFloat(val.Amount.String(), 64)
		if err != nil {
			// fatal error because table includes invalid data
			t.logger.Error("payment_request table includes invalid amount field")
			return nil, 0, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if !xrp.ValidateAddress(userPayments[idx].receiverAddr) {
			// fatal error
			t.logger.Error("address is invalid",
				zap.String("address", userPayments[idx].receiverAddr),
				zap.Error(err),
			)
			return nil, 0, nil, errors.Wrapf(err, "address is invalid: %s", userPayments[idx].receiverAddr)
		}

		// total amount
		totalAmount += amt
	}

	return userPayments, totalAmount, paymentRequestIds, nil
}
