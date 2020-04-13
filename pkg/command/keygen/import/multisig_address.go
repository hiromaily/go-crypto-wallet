package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

func (c *MultisigCommand) Help() string {
	return `Usage: keygen key import multisig [options...]
Options:
  -file     signed transaction file path
  -account  target account
`
}

func (c *MultisigCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		filePath string
		acnt     string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&filePath, "file", "", "import file path for signed transactions")
	flags.StringVar(&acnt, "account", "", "target account")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if filePath == "" {
		c.ui.Error("file path option [-file] is required")
		return 1
	}
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		c.ui.Error(fmt.Sprintf("account: %s/%s is not allowd", account.AccountTypeAuthorization, account.AccountTypeClient))
		return 1
	}

	// import multisig addresses exported by signature wallet from csv file to database
	err := c.wallet.ImportMultisigAddrForColdWallet1(filePath, account.AccountType(acnt))
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportMultisigAddrForColdWallet1() %+v", err))
		return 1
	}
	c.ui.Output("Done!")

	return 0
}
