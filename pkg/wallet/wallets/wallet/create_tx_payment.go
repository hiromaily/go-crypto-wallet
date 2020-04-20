package wallet

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
)

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
func (w *Wallet) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypePayment
	receiver := account.AccountType("")
	targetAction := action.ActionTypePayment

	// get payment data from payment_request
	userPayments, paymentRequestIds, err := w.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		w.logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// calculate total amount to send from payment_request
	var requiredAmount btcutil.Amount
	for _, val := range userPayments {
		requiredAmount += val.validAmount
	}

	// get balance for payment account
	balance, err := w.btc.GetBalanceByAccount(account.AccountTypePayment)
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		//balance is short
		return "", "", errors.New("balance for payment account is insufficient")
	}
	w.logger.Debug("payment balane and userTotal",
		zap.Any("balance", balance),
		zap.Any("userTotal", requiredAmount))

	// create transfer transaction
	return w.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, paymentRequestIds, userPayments)
}

// userPayments is given for receiverAddr
func (w *Wallet) createPaymentOutputs(userPayments []UserPayment, changeAddr string, changeAmount btcutil.Amount) map[btcutil.Address]btcutil.Amount {
	var (
		txOutputs = map[btcutil.Address]btcutil.Amount{}
		//if key of map is btcutil.Address which is interface type, uniqueness can't be found from map key
		// so this key is string
		tmpOutputs = map[string]btcutil.Amount{}
	)

	// create txOutput from userPayment
	for _, userPayment := range userPayments {
		if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
			// if there are multiple receiver same address in user_payment,
			//  sum these
			tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
		} else {
			// new receiver address
			tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
		}
	}

	// add txOutput as change address and amount for change
	//TODO:
	// - what if user register for address which is same to payment address?
	//   Though it's impossible in real but systematically possible
	// - BIP44, hdwallet has `ChangeType`. ideally this address should be used
	if _, ok := tmpOutputs[changeAddr]; ok {
		// in real, this is impossible
		tmpOutputs[changeAddr] += changeAmount
	} else {
		// of course, change address is new in txOutputs
		tmpOutputs[changeAddr] = changeAmount
	}

	// create txOutputs from tmpOutputs switching string address type to btcutil.Address
	for strAddr, amount := range tmpOutputs {
		addr, err := w.btc.DecodeAddress(strAddr)
		if err != nil {
			// this case is impossible because addresses are checked in advance
			w.logger.Error("fail to call DecodeAddress",
				zap.String("address", strAddr))
			continue
		}
		txOutputs[addr] = amount
	}

	return txOutputs
}
