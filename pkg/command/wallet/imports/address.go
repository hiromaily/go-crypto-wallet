package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// AddressCommand address subcommand
type AddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *AddressCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *AddressCommand) Help() string {
	return `Usage: wallet import address [options...]
Options:
  -file  import file path for generated addresses
`
}

// Run executes this subcommand
func (c *AddressCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		filePath string
		acnt     string
		isRescan bool
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&filePath, "file", "", "import file path for generated addresses")
	flags.StringVar(&acnt, "account", "", "user account")
	flags.BoolVar(&isRescan, "rescan", false, "run rescan when importing addresses or not")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-file: %s", filePath))

	//validator
	if filePath == "" {
		c.ui.Error("file path option [-file] is required")
		return 1
	}
	if !account.ValidateAccountType(acnt) {
		c.ui.Error("account option [-account] is invalid")
		return 1
	}
	if account.AccountTypeAuthorization.Is(acnt) {
		c.ui.Error(fmt.Sprintf("account: %s is not allowed", account.AccountTypeAuthorization))
		return 1
	}

	//import public address
	err := c.wallet.ImportAddress(filePath, account.AccountType(acnt), isRescan)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportAddress() %+v", err))
		return 1
	}
	c.ui.Info("Done!")

	return 0
}
