package watchsrv

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
)

// CreateDepositTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateDepositTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypeClient
	receiver := account.AccountTypeDeposit
	targetAction := action.ActionTypeDeposit
	requiredAmount, err := t.btc.FloatToAmount(0)
	if err != nil {
		return "", "", err
	}

	// create deposit transaction
	return t.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}
