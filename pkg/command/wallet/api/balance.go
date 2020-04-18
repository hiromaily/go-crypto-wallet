package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//balance subcommand
type BalanceCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *BalanceCommand) Synopsis() string {
	return c.synopsis
}

func (c *BalanceCommand) Help() string {
	return `Usage: wallet api balance [options...]
Options:
  -account  account
`
}

func (c *BalanceCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		acnt string
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "account")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}

	//TODO: get all balance without account
	var (
		balance float64
		err     error
	)
	if acnt == "" {
		balance, err = c.wallet.GetBTC().GetBalance()
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.GetBalance() %+v", err))
			return 1
		}
	} else {
		//get received by account
		balance, err = c.wallet.GetBTC().GetBalanceByAccount(account.AccountType(acnt))
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.GetBalanceByAccount() %+v", err))
			return 1
		}
	}

	// FIXME: even spent tx looks to be left, GetReceivedByLabelAndMinConf may be wrong to get balance
	//balance, err := c.wallet.GetBTC().GetReceivedByLabelAndMinConf(acnt, c.wallet.GetBTC().ConfirmationBlock())
	//if err != nil {
	//	c.ui.Error(fmt.Sprintf("fail to call BTC.GetReceivedByAccountAndMinConf() %+v", err))
	//	return 1
	//}

	c.ui.Info(fmt.Sprintf("balance: %v", balance))

	return 0
}
