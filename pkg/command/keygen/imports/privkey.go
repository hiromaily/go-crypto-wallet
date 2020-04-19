package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// PrivKeyCommand privkey subcommand
type PrivKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *PrivKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *PrivKeyCommand) Help() string {
	return `Usage: keygen key import privkey [options...]
Options:
  -account  target account
`
}

// Run executes this subcommand
func (c *PrivKeyCommand) Run(args []string) int {
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
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization}) {
		c.ui.Error(fmt.Sprintf("account: %s is not allowed", account.AccountTypeAuthorization))
		return 1
	}

	//import generated private key to keygen wallet
	err := c.wallet.ImportPrivateKey(account.AccountType(acnt))
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportPrivateKey() %+v", err))
		return 1
	}
	c.ui.Output("Done!")

	return 0
}
