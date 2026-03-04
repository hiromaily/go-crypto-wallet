package xrp

import (
	"context"

	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
)

// GetBalance returns amount of address
func (w *WSClient) GetBalance(ctx context.Context, addr string) (float64, error) {
	accountInfo, err := w.GetAccountInfo(ctx, addr)
	if err != nil {
		return 0, err
	}
	return xrpkg.ToFloat64(accountInfo.XrpBalance), nil
}

// GetTotalBalance returns total amount in address list
func (w *WSClient) GetTotalBalance(ctx context.Context, addrs []string) float64 {
	var total float64
	for _, addr := range addrs {
		amt, err := w.GetBalance(ctx, addr)
		if err == nil {
			total += amt
		}
	}
	return total
}
