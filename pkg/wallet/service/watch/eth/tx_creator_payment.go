package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
// Note
// - to avoid complex logic to create raw transaction
// - only one address of sender should afford to send coin to all payment request users.
func (t *TxCreate) CreatePaymentTx() (string, string, error) {
	sender := t.paymentSender
	receiver := domainAccount.AccountTypeAnonymous
	targetAction := domainTx.ActionTypePayment
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	// get payment data from payment_request
	userPayments, totalAmount, paymentRequestIds, err := t.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// get sender address
	senderAddr, err := t.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetAll(domainAccount.AccountTypeClient): %w", err)
	}
	err = t.validateAmount(senderAddr, totalAmount)
	if err != nil {
		return "", "", err
	}

	// create raw transaction each address
	serializedTxs, txDetailItems, err := t.createPaymentRawTransactions(sender, receiver, userPayments, senderAddr)
	if err != nil {
		return "", "", err
	}
	if len(txDetailItems) == 0 {
		return "", "", nil
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
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return "", generatedFileName, nil
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
		return nil, nil, nil, fmt.Errorf("fail to call repo.GetPaymentRequestAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		logger.Debug("no data in payment_request")
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
		var amt float64
		amt, err = strconv.ParseFloat(val.Amount.String(), 64)
		if err != nil {
			// fatal error because table includes invalid data
			logger.Error("payment_request table includes invalid amount field")
			return nil, nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if err = t.eth.ValidateAddr(userPayments[idx].receiverAddr); err != nil {
			// fatal error
			logger.Error("fail to call ValidationAddr",
				"address", userPayments[idx].receiverAddr,
				"error", err,
			)
			return nil, nil, nil, fmt.Errorf("fail to call eth.ValidationAddr(): %w", err)
		}

		// amount
		userPayments[idx].amount = t.eth.FloatToBigInt(userPayments[idx].floatAmount)
		totalAmount = new(big.Int).Add(totalAmount, userPayments[idx].amount)
	}

	return userPayments, totalAmount, paymentRequestIds, nil
}

func (t *TxCreate) validateAmount(senderAddr *models.Address, totalAmount *big.Int) error {
	// check sender's total balance
	senderBalance, err := t.eth.GetBalance(context.TODO(), senderAddr.WalletAddress, eth.QuantityTagPending)
	if err != nil {
		return fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}

	if senderBalance.Uint64() <= totalAmount.Uint64() {
		return errors.New("sender balance is insufficient to send")
	}
	return nil
}

func (t *TxCreate) createPaymentRawTransactions(
	sender, receiver domainAccount.AccountType, userPayments []UserPayment, senderAddr *models.Address,
) ([]string, []*models.EthDetailTX, error) {
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userPayments))
	additionalNonce := 0
	for _, userPayment := range userPayments {
		// call CreateRawTransaction
		rawTx, txDetailItem, err := t.eth.CreateRawTransaction(context.TODO(),
			senderAddr.WalletAddress, userPayment.receiverAddr, userPayment.amount.Uint64(), additionalNonce)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"fail to call addrRepo.CreateRawTransaction(), sender address: %s: %w",
				senderAddr.WalletAddress, err)
		}
		additionalNonce++

		rawTxHex := rawTx.TxHex
		logger.Debug("rawTxHex", "rawTxHex", rawTxHex)
		// TODO: `rawTxHex` should be used to trace progress to update database

		serializedTx, err := serial.EncodeToString(rawTx)
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
