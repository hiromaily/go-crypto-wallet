package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//unlocktx subcommand
type UnLockTxCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *UnLockTxCommand) Synopsis() string {
	return c.synopsis
}

func (c *UnLockTxCommand) Help() string {
	return `Usage: wallet api unlocktx`
}

func (c *UnLockTxCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// unlock locked transaction for unspent transaction
	err := c.wallet.GetBTC().UnlockAllUnspentTransaction()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.UnlockAllUnspentTransaction() %+v", err))
		return 1
	}

	return 0
}
