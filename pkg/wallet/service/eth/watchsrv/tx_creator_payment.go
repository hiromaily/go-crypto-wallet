package watchsrv

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
// TODO: implement
func (t *TxCreate) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	return "", "", nil
}
