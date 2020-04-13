package add

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Signer
}

func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

func (c *MultisigCommand) Help() string {
	return `Usage: sign key add multisig [options...]
Options:
  -account  target account
`
}

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
		c.ui.Error(fmt.Sprintf("account: %s/%s is not allowd", account.AccountTypeAuthorization, account.AccountTypeClient))
		return 1
	}

	// call `addmultisigaddress` which adds a P2SH multisig address to the wallet
	//  address would be acount and authorization
	err := c.wallet.AddMultisigAddressByAuthorization(account.AccountType(acnt), address.AddressTypeP2shSegwit)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call AddMultisigAddressByAuthorization() %+v", err))
		return 1
	}
	c.ui.Output("Done!")

	return 0
}
