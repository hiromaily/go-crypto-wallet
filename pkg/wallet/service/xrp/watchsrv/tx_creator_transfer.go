package watchsrv

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
)

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
// - sender pays fee,
// - any internal account should have only one address in Ethereum because no utxo
func (t *TxCreate) CreateTransferTx(sender, receiver account.AccountType, floatValue float64) (string, string, error) {
	return "", "", nil
}
