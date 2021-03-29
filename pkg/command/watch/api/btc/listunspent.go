package btc

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// ListUnspentCommand listunspent subcommand
type ListUnspentCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      btcgrp.Bitcoiner
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
  -num      confirmation number
`
}

// Run executes this subcommand
func (c *ListUnspentCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		acnt            string
		argsNum         int64
		confirmationNum uint64
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "account")
	flags.Int64Var(&argsNum, "num", -1, "confirmation number")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validator
	if acnt != "" && !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}
	if argsNum == -1 {
		confirmationNum = c.btc.ConfirmationBlock()
	} else {
		confirmationNum = uint64(argsNum)
	}

	if acnt != "" {
		// call listunspent
		unspentList, err := c.btc.ListUnspentByAccount(account.AccountType(acnt), confirmationNum)
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.ListUnspentByAccount() %+v", err))
			return 1
		}
		grok.Value(unspentList)

		unspentAddrs := c.btc.GetUnspentListAddrs(unspentList, account.AccountType(acnt))
		for _, addr := range unspentAddrs {
			grok.Value(addr)
		}
	} else {
		// call listunspent
		// ListUnspentMin doesn't have proper response, label can't be retrieved

		unspentList, err := c.btc.ListUnspent(confirmationNum)
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call btc.ListUnspent() %+v", err))
			return 1
		}
		grok.Value(unspentList)
	}

	return 0
}
