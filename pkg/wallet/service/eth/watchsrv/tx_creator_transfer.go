package watchsrv

import "github.com/hiromaily/go-bitcoin/pkg/account"

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
// TODO: implement
func (t *TxCreate) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	return "", "", nil
}
