package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/keystatus"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//address subcommand
type AddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

func (c *AddressCommand) Synopsis() string {
	return c.synopsis
}

func (c *AddressCommand) Help() string {
	return `Usage: keygen key export address [options...]
Options:
  -account  target account
`
}

func (c *AddressCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

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
		c.ui.Error(fmt.Sprintf("account: %s is not allowd", account.AccountTypeAuthorization))
		return 1
	}

	// export generated PublicKey as csv file to use at watch only wallet
	fileName, err := c.wallet.ExportAccountKey(account.AccountType(acnt), keystatus.KeyStatusImportprivkey)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ExportAccountKey() %+v", err))
		return 1
	}
	c.ui.Output(fmt.Sprintf("[fileName]: %s", fileName))

	return 0
}
