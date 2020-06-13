package watchsrv

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account(payment) covers fee, but is should be flexible
// Note
// - to avoid complex logic to create raw transaction
// - only one address of sender should afford to send coin to all payment request users.
func (t *TxCreate) CreatePaymentTx() (string, string, error) {
	return "", "", nil
}
