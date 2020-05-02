package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// TODO
//  - how to display help? upper layer's help displays by `wallet deposit create -h`
//  - as workaround, add undefined flag like `wallet deposit create -a`

// DepositCommand deposit subcommand
type DepositCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *DepositCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *DepositCommand) Help() string {
	return `Usage: wallet create deposit [options...]
Options:
  -fee    adjustment fee
`
}

// Run executes this subcommand
func (c *DepositCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		fee float64
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Detect transaction for clients from blockchain network and create deposit unsigned transaction
	// It would be run manually on the daily basis because signature is manual task
	hex, fileName, err := c.wallet.CreateDepositTx(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call CreateDepositTx() %+v", err))
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
