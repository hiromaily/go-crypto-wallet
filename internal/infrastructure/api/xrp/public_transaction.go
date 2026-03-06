package xrp

import (
	"context"
	"fmt"

	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

// https://xrpl.org/transaction-methods.html

// Sign calls sign method
func (w *WSClient) Sign(
	ctx context.Context, txJSON *xrprpcpublic.PublicTxType, secret string, offline bool,
) (*xrprpcpublic.ResponseSign, error) {
	res, err := w.publicRPC.Sign(ctx, txJSON, secret, offline)
	if err != nil {
		return nil, fmt.Errorf("fail to call publicRPC.Sign: %w", err)
	}
	return res, nil
}
