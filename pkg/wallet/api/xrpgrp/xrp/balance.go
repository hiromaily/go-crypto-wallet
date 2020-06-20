package xrp

// GetBalance returns amount of address
func (r *Ripple) GetBalance(addr string) (float64, error) {
	accountInfo, err := r.GetAccountInfo(addr)
	if err != nil {
		return 0, err
	}
	return ToFloat64(accountInfo.XrpBalance), nil
}

// GetTotalBalance returns total amount in address list
func (r *Ripple) GetTotalBalance(addrs []string) float64 {
	var total float64
	for _, addr := range addrs {
		amt, err := r.GetBalance(addr)
		if err == nil {
			total += amt
		}
	}
	return total
}
