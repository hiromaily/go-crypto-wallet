package watchsrv

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
)

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
func (t *TxCreate) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	sender := t.paymentSender
	receiver := account.AccountTypeAnonymous
	targetAction := action.ActionTypePayment
	t.logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	// get payment data from payment_request
	userPayments, paymentRequestIds, err := t.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		t.logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// calculate total amount to send from payment_request
	var requiredAmount btcutil.Amount
	for _, val := range userPayments {
		requiredAmount += val.validAmount
	}

	// get balance for payment account
	balance, err := t.btc.GetBalanceByAccount(account.AccountTypePayment, t.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		// balance is short
		t.logger.Info("balance for payment account is insufficient",
			"payment_balance", balance.ToBTC(),
			"required_amount", requiredAmount.ToBTC(),
		)
		return "", "", nil
	}
	t.logger.Debug("payment balane and userTotal",
		"balance", balance,
		"userTotal", requiredAmount)

	// create transfer transaction
	return t.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, paymentRequestIds, userPayments)
}

// userPayments is given for receiverAddr
func (t *TxCreate) createPaymentTxOutputs(
	userPayments []UserPayment, changeAddr string, changeAmount btcutil.Amount,
) map[btcutil.Address]btcutil.Amount {
	var (
		txOutputs = map[btcutil.Address]btcutil.Amount{}
		// if key of map is btcutil.Address which is interface type, uniqueness can't be found from map key
		// so this key is string
		tmpOutputs = map[string]btcutil.Amount{}
	)

	// create txOutput from userPayment
	for _, userPayment := range userPayments {
		tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
	}

	// add txOutput as change address and amount for change
	// TODO:
	// - what if user register for address which is same to payment address?
	//   Though it's impossible in real but systematically possible
	// - BIP44, hdwallet has `ChangeType`. ideally this address should be used
	tmpOutputs[changeAddr] += changeAmount

	// create txOutputs from tmpOutputs switching string address type to btcutil.Address
	for strAddr, amount := range tmpOutputs {
		addr, err := t.btc.DecodeAddress(strAddr)
		if err != nil {
			// this case is impossible because addresses are checked in advance
			t.logger.Error("fail to call DecodeAddress",
				"address", strAddr)
			continue
		}
		txOutputs[addr] = amount
	}

	return txOutputs
}

// UserPayment user's payment address and amount
type UserPayment struct {
	senderAddr   string          // sender address for just chacking
	receiverAddr string          // receiver address
	validRecAddr btcutil.Address // decoded receiver address
	amount       float64         // amount
	validAmount  btcutil.Amount  // decoded amount
}

// createUserPayment get payment data from payment_request table
func (t *TxCreate) createUserPayment() ([]UserPayment, []int64, error) {
	// get payment_request
	paymentRequests, err := t.payReqRepo.GetAll()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call repo.GetPaymentRequestAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		t.logger.Debug("no data in payment_request")
		return nil, nil, nil
	}

	userPayments := make([]UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.SenderAddress
		userPayments[idx].receiverAddr = val.ReceiverAddress
		amt, parseErr := strconv.ParseFloat(val.Amount.String(), 64)
		if parseErr != nil {
			// fatal error because table includes invalid data
			t.logger.Error("payment_request table includes invalid amount field")
			return nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].amount = amt

		// decode address
		// TODO: may be it's not necessary
		userPayments[idx].validRecAddr, err = t.btc.DecodeAddress(userPayments[idx].receiverAddr)
		if err != nil {
			// fatal error
			t.logger.Error("unexpected error occurred converting receiverAddr from string type  to address type")
			return nil, nil, errors.New("unexpected error occurred converting receiverAddr from string type to address type")
		}

		// amount
		userPayments[idx].validAmount, err = t.btc.FloatToAmount(userPayments[idx].amount)
		if err != nil {
			// fatal error
			t.logger.Error("unexpected error occurred converting amount from float64 type to Amount type")
			return nil, nil, errors.New("unexpected error occurred converting amount from float64 type to Amount type")
		}
	}

	return userPayments, paymentRequestIds, nil
}
