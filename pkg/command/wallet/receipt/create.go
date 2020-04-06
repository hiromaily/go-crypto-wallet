package receipt

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

const createName = "create"

//create subcommand
type CreateTxCommand struct {
	ui     cli.Ui
	wallet *wallet.Wallet
}

func (c *CreateTxCommand) Synopsis() string {
	return "detect receipt to our addressed and crate receipt unsigned transaction"
}

func (c *CreateTxCommand) Help() string {
	return `Usage: wallet receipt create [options...]
Options:
  -fee  adjustment fee
`
}

func (c *CreateTxCommand) Run(args []string) int {
	var (
		fee float64
	)
	flags := flag.NewFlagSet(createName, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-fee: %f", fee))

	// Detect receipt transaction from outside and create receipt unsigned transaction
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
