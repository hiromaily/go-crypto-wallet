package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	transferName = "transfer"
)

//transfer subcommand
type TransferCommand struct {
	ui     cli.Ui
	wallet *service.Wallet
}

func (c *TransferCommand) Synopsis() string {
	return "create unsigned transaction for transfer among accounts"
}

func (c *TransferCommand) Help() string {
	return `Usage: wallet transfer [options...]
Options:
  -account1  account for transfer from
  -account2  account for transfer to
`
}

//WIP
func (c *TransferCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		account1 string
		account2 string
	)
	flags := flag.NewFlagSet(transferName, flag.ContinueOnError)
	flags.StringVar(&account1, "account1", "", "account for transfer from")
	flags.StringVar(&account2, "account2", "", "account for transfer to")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if !account.ValidateAccountType(account1) {
		c.ui.Error("account option [-account1] is invalid")
		return 1
	}
	if !account.ValidateAccountType(account2) {
		c.ui.Error("account option [-account2] is invalid")
		return 1
	}
	if !account.NotAllow(account1, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		c.ui.Error(fmt.Sprintf("account1: %s/%s is not allowd", account.AccountTypeAuthorization, account.AccountTypeClient))
		return 1
	}
	if !account.NotAllow(account2, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
		c.ui.Error(fmt.Sprintf("account2: %s/%s is not allowd", account.AccountTypeAuthorization, account.AccountTypeClient))
		return 1
	}

	//TODO: amount should be set
	hex, fileName, err := c.wallet.SendToAccount(account.AccountType(account1), account.AccountType(account2), 0)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call SendToAccount() %+v", err))
		return 1
	}
	if hex == "" {
		c.ui.Info("No utxo")
		return 0
	}
	//TODO: output should be json if json option is true
	c.ui.Output(fmt.Sprintf("[hex]: %s\n[fileName]: %s", hex, fileName))

	return 0
}
