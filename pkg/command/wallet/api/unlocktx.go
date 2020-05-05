package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// UnLockTxCommand unlocktx subcommand
type UnLockTxCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *UnLockTxCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *UnLockTxCommand) Help() string {
	return `Usage: wallet api unlocktx`
}

// Run executes this subcommand
func (c *UnLockTxCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// unlock locked transaction for unspent transaction
	err := c.btc.UnlockUnspent()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.UnlockUnspent() %+v", err))
		return 1
	}

	return 0
}
