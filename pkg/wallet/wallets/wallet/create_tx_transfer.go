package wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/pkg/errors"
)

// CreateTransferTx transfer coin among internal account except client, authorization
// - if amount=0, all coin is sent
// TODO: temporary fixed account is used `receipt to payment`
// TODO: after this func, what if `listtransactions` api is called to see result
func (w *Wallet) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	targetAction := action.ActionTypeTransfer

	// validation
	if receiver == account.AccountTypeClient || receiver == account.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	//amount btcutil.Amount
	amount, err := w.btc.FloatBitToAmount(floatAmount)
	if err != nil {
		return "", "", err
	}

	// check balance for sender
	balance, err := w.btc.GetReceivedByLabelAndMinConf(sender.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= amount {
		//balance is short
		return "", "", errors.Errorf("account: %s balance is insufficient", sender)
	}

	// create transfer transaction
	return w.createTx(sender, receiver, targetAction, amount, adjustmentFee)
}
