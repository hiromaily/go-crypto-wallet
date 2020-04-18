package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// TODO
//  - how to display help? upper layer's help displays by `wallet receipt create -h`
//  - as workaround, add undefined flag like `wallet receipt create -a`

// ReceiptCommand receipt subcommand
type ReceiptCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis
func (c *ReceiptCommand) Synopsis() string {
	return c.synopsis
}

// Help
func (c *ReceiptCommand) Help() string {
	return `Usage: wallet create receipt [options...]
Options:
  -fee    adjustment fee
`
}

// Run
func (c *ReceiptCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		fee float64
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Detect transaction for clients from blockchain network and create receipt unsigned transaction
	// It would be run manually on the daily basis because signature is manual task
	hex, fileName, err := c.wallet.CreateReceiptTx(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call CreateReceiptTx() %+v", err))
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
