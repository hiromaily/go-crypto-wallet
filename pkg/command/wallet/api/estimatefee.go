package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//estimatefee subcommand
type EstimateFeeCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *EstimateFeeCommand) Synopsis() string {
	return c.synopsis
}

func (c *EstimateFeeCommand) Help() string {
	return `Usage: wallet api estimatefee`
}

func (c *EstimateFeeCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// estimate fee
	feePerKb, err := c.wallet.BTC.EstimateSmartFee()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.EstimateSmartFee() %+v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("EstimateSmartFee: %f", feePerKb))

	return 0
}
