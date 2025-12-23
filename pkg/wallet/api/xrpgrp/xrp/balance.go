package xrp

import "context"

// GetBalance returns amount of address
func (r *Ripple) GetBalance(ctx context.Context, addr string) (float64, error) {
	accountInfo, err := r.GetAccountInfo(ctx, addr)
	if err != nil {
		return 0, err
	}
	return ToFloat64(accountInfo.XrpBalance), nil
}

// GetTotalBalance returns total amount in address list
func (r *Ripple) GetTotalBalance(ctx context.Context, addrs []string) float64 {
	var total float64
	for _, addr := range addrs {
		amt, err := r.GetBalance(ctx, addr)
		if err == nil {
			total += amt
		}
	}
	return total
}
