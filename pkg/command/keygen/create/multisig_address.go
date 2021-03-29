package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// MultisigCommand  multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *MultisigCommand) Help() string {
	return `Usage: keygen create multisig [options...]
Options:
  -account  target account
`
}

// Run executes this subcommand
func (c *MultisigCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var acnt string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "target account")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validator
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}

	// create multisig address
	err := c.wallet.CreateMultisigAddress(account.AccountType(acnt))
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call CreateMultisigAddress() %+v", err))
		return 1
	}

	return 0
}
