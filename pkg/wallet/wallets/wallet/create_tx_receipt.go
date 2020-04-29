package wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
)

// CreateReceiptTx create unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but is should be flexible
func (w *Wallet) CreateReceiptTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypeClient
	receiver := account.AccountTypeDeposit
	targetAction := action.ActionTypeDeposit
	requiredAmount, err := w.btc.FloatToAmount(0)
	if err != nil {
		return "", "", err
	}

	// create deposit transaction
	return w.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}
