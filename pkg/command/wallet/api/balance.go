package api

import (
	"flag"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// BalanceCommand balance subcommand
type BalanceCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *BalanceCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *BalanceCommand) Help() string {
	return `Usage: wallet api balance [options...]
Options:
  -account  account
`
}

// Run executes this subcommand
func (c *BalanceCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var acnt string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "account")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if acnt != "" && !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}

	var (
		balance btcutil.Amount
		err     error
	)
	if acnt == "" {
		balance, err = c.btc.GetBalance()
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.GetBalance() %+v", err))
			return 1
		}
	} else {
		//get received by account
		balance, err = c.btc.GetBalanceByAccount(account.AccountType(acnt))
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
