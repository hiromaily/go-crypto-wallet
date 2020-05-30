package xrp

// PrepareTransaction calls PrepareTransaction API
func (r *Ripple) PrepareTransaction(senderAccount, receiverAccount string, amount float64) error {
	return r.api.PrepareTransaction(senderAccount, receiverAccount, amount)
}
