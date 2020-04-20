package api

import (
	"flag"
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/account"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ListUnspentCommand listunspent subcommand
type ListUnspentCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *ListUnspentCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ListUnspentCommand) Help() string {
	return `Usage: wallet api listunspent [options...]
Options:
  -account  account
`
}

// Run executes this subcommand
func (c *ListUnspentCommand) Run(args []string) int {
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

	if acnt != "" {
		// call listunspent
		unspentList, unspentAddrs, err := c.wallet.GetBTC().ListUnspentByAccount(account.AccountType(acnt))
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.ListUnspentByAccount() %+v", err))
			return 1
		}
		grok.Value(unspentList)
		for _, addr := range unspentAddrs {
			grok.Value(addr)
		}
	} else {
		// call listunspent
		// ListUnspentMin doesn't have proper response, label can't be retrieved
		unspentList, err := c.wallet.GetBTC().ListUnspent()
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.ListUnspent() %+v", err))
			return 1
		}
		grok.Value(unspentList)
	}

	return 0
}
