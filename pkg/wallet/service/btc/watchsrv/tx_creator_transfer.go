package watchsrv

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
)

// CreateTransferTx create unsigned tx for transfer coin among internal account except client, authorization
// FIXME: for now, receiver account covers fee, but is should be flexible
func (t *TxCreate) CreateTransferTx(sender, receiver account.AccountType, floatAmount, adjustmentFee float64) (string, string, error) {
	targetAction := action.ActionTypeTransfer

	// validation account
	if receiver == account.AccountTypeClient || receiver == account.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	//amount btcutil.Amount
	requiredAmount, err := t.btc.FloatToAmount(floatAmount)
	if err != nil {
		return "", "", err
	}

	// check balance for sender
	balance, err := t.btc.GetBalanceByAccount(sender)
	//balance, err := w.btc.GetReceivedByLabelAndMinConf(sender.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		//balance is short
		return "", "", errors.Errorf("account: %s balance is insufficient", sender)
	}

	// create transfer transaction
	return t.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}
