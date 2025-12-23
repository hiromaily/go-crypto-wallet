package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runBalance(btc btcgrp.Bitcoiner, acnt string) error {
	// validator
	if acnt != "" && !account.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}

	var (
		balance btcutil.Amount
		err     error
	)
	if acnt == "" {
		balance, err = btc.GetBalance()
		if err != nil {
			return fmt.Errorf("fail to call btc.GetBalance() %w", err)
		}
	} else {
		// get received by account
		balance, err = btc.GetBalanceByAccount(account.AccountType(acnt), btc.ConfirmationBlock())
		if err != nil {
			return fmt.Errorf("fail to call btc.GetBalanceByAccount() %w", err)
		}
	}

	// FIXME: even spent tx looks to be left, GetReceivedByLabelAndMinConf may be wrong to get balance
	// balance, err := wallet.GetBTC().GetReceivedByLabelAndMinConf(acnt, wallet.GetBTC().ConfirmationBlock())
	// if err != nil {
	//	return fmt.Errorf("fail to call BTC.GetReceivedByAccountAndMinConf() %w", err)
	//}

	fmt.Printf("balance: %v\n", balance)

	return nil
}
