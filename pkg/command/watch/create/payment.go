package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// PaymentCommand payment subcommand
type PaymentCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *PaymentCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *PaymentCommand) Help() string {
	return `Usage: wallet create payment [options...]
Options:
  -fee  adjustment fee
`
}

// Run executes this subcommand
func (c *PaymentCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		fee float64
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Create payment transaction
	hex, fileName, err := c.wallet.CreatePaymentTx(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call CreatePaymentTx() %+v", err))
		return 1
	}
	if (c.wallet.CoinTypeCode() != coin.ETH && c.wallet.CoinTypeCode() != coin.XRP) && hex == "" {
		c.ui.Info("No utxo")
		return 0
	}
	//TODO: output should be json if json option is true
	c.ui.Output(fmt.Sprintf("[hex]: %s\n[fileName]: %s", hex, fileName))

	return 0
}
