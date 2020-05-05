package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// EstimateFeeCommand estimatefee subcommand
type EstimateFeeCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *EstimateFeeCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *EstimateFeeCommand) Help() string {
	return `Usage: wallet api estimatefee`
}

// Run executes this subcommand
func (c *EstimateFeeCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// estimate fee
	feePerKb, err := c.btc.EstimateSmartFee()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.EstimateSmartFee() %+v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("EstimateSmartFee: %f", feePerKb))

	return 0
}
