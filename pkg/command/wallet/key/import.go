package key

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const importName = "import"

//import subcommand
type ImportCommand struct {
	ui     cli.Ui
	wallet *service.Wallet
}

func (c *ImportCommand) Synopsis() string {
	return "import generatd addresses by keygen wallet"
}

func (c *ImportCommand) Help() string {
	return `Usage: wallet key import [options...]
Options:
  -file  import file path for generated addresses
`
}

func (c *ImportCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		filePath string
		acnt     string
		isRescan bool
	)
	flags := flag.NewFlagSet(importName, flag.ContinueOnError)
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
		c.ui.Error(fmt.Sprintf("account: %s is not allowd", account.AccountTypeAuthorization))
		return 1
	}

	//import public key
	err := c.wallet.ImportPublicKeyForWatchWallet(filePath, account.AccountType(acnt), isRescan)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportPublicKeyForWatchWallet() %+v", err))
		return 1
	}
	c.ui.Info("Done!")

	return 0
}
