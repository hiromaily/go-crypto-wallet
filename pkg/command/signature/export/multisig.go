package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// MultisigCommand multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *MultisigCommand) Help() string {
	return `Usage: sign export multisig [options...]
Options:
  -account  target account
`
}

// Run executes this subcommand
func (c *MultisigCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		acnt string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "target account")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		c.ui.Error(fmt.Sprintf("account: %s/%s is not allowed", account.AccountTypeAuthorization, account.AccountTypeClient))
		return 1
	}

	// export created multisig address as csv file
	fileName, err := c.wallet.ExportAddedPubkeyHistory(account.AccountType(acnt))
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ExportAddedPubkeyHistory() %+v", err))
		return 1
	}
	c.ui.Output(fmt.Sprintf("[fileName]: %s", fileName))

	return 0
}
