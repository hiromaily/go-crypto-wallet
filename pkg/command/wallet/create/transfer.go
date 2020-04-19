package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// TransferCommand transfer subcommand
type TransferCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *TransferCommand) Synopsis() string {
	return "create unsigned transaction for transfer among accounts"
}

// Help returns usage for this subcommand
func (c *TransferCommand) Help() string {
	return `Usage: wallet create transfer [options...]
Options:
  -account1  sender account
  -account2  receiver account
  -amount    amount to send coin. if amount=0, all coin is sent
  -fee       adjustment fee
`
}

// Run executes this subcommand
func (c *TransferCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		account1 string
		account2 string
		amount   float64
		fee      float64
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&account1, "account1", "", "sender account")
	flags.StringVar(&account2, "account2", "", "receiver account")
	flags.Float64Var(&amount, "amount", 0, "amount to send coin")
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
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
	// This logic should be implemented in wallet.CreateTransferTx()
	//if !account.NotAllow(account1, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
	//	c.ui.Error(fmt.Sprintf("account1: %s/%s is not allowed", account.AccountTypeAuthorization, account.AccountTypeClient))
	//	return 1
	//}
	//if !account.NotAllow(account2, []account.AccountType{account.AccountTypeAuthorization, account.AccountTypeClient}) {
	//	c.ui.Error(fmt.Sprintf("account2: %s/%s is not allowed", account.AccountTypeAuthorization, account.AccountTypeClient))
	//	return 1
	//}
	//if amount == 0{
	//	c.ui.Error("amount option [-amount] is invalid")
	//}

	hex, fileName, err := c.wallet.CreateTransferTx(
		account.AccountType(account1),
		account.AccountType(account2),
		amount,
		fee)

	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call CreateTransferTx() %+v", err))
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
