package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddressCommand address subcommand
type AddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *AddressCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *AddressCommand) Help() string {
	return `Usage: keygen key export address [options...]
Options:
  -account  target account
`
}

// Run executes this subcommand
func (c *AddressCommand) Run(args []string) int {
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

	// export generated PublicKey as csv file
	fileName, err := c.wallet.ExportAddress(account.AccountType(acnt))
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ExportAddress() %+v", err))
		return 1
	}
	c.ui.Output(fmt.Sprintf("[fileName]: %s", fileName))

	return 0
}
