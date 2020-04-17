package wallet

import (
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// UserPayment user's payment address and amount
type UserPayment struct {
	senderAddr   string          // sender address for just chacking
	receiverAddr string          // receiver address
	validRecAddr btcutil.Address // decoded receiver address
	amount       float64         // amount
	validAmount  btcutil.Amount  // decoded amount
}

// createUserPayment get payment data from payment_request table
func (w *Wallet) createUserPayment() ([]UserPayment, []int64, error) {
	// get payment_request
	paymentRequests, err := w.storager.GetPaymentRequestAll()
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call storager.GetPaymentRequestAll()")
	}
	if len(paymentRequests) == 0 {
		w.logger.Debug("no data in payment_request")
		return nil, nil, nil
	}

	userPayments := make([]UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.AddressFrom
		userPayments[idx].receiverAddr = val.AddressTo
		amt, err := strconv.ParseFloat(val.Amount, 64)
		if err != nil {
			// fatal error because table includes invalid data
			w.logger.Error("payment_request table includes invalid amount field")
			return nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].amount = amt

		// decode address
		// TODO: may be it's not necessary
		userPayments[idx].validRecAddr, err = w.btc.DecodeAddress(userPayments[idx].receiverAddr)
		if err != nil {
			// fatal error
			w.logger.Error("unexpected error occurred converting receiverAddr from string type  to address type")
			return nil, nil, errors.New("unexpected error occurred converting receiverAddr from string type to address type")
		}

		// amount
		userPayments[idx].validAmount, err = w.btc.FloatBitToAmount(userPayments[idx].amount)
		if err != nil {
			// fatal error
			w.logger.Error("unexpected error occurred converting amount from float64 type to Amount type")
			return nil, nil, errors.New("unexpected error occurred converting amount from float64 type to Amount type")
		}
	}

	return userPayments, paymentRequestIds, nil
}
