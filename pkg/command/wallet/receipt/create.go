package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

// TODO
//  - how to display help? upper layer's help displays by `wallet receipt create -h`
//  - as workaround, add undefined flag like `wallet receipt create -a`

//create subcommand
type CreateTxCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *CreateTxCommand) Synopsis() string {
	return c.synopsis
}

func (c *CreateTxCommand) Help() string {
	return `Usage: wallet receipt create [options...]
Options:
  -fee    adjustment fee
  -check  only check client addresses, not create unsigned transaction
`
}

func (c *CreateTxCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		fee     float64
		isCheck bool
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	flags.BoolVar(&isCheck, "check", false, "only check client addresses, not create unsigned transaction")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-fee: %f", fee))

	// Detect transaction for clients from blockchain network and create receipt unsigned transaction
	// It would be run manually on the daily basis because signature is manual task
	hex, fileName, err := c.wallet.DetectReceivedCoin(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call DetectReceivedCoin() %+v", err))
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
