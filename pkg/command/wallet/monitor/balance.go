package monitor

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// BalanceCommand balance subcommand
type BalanceCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *BalanceCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *BalanceCommand) Help() string {
	return `Usage: wallet monitor balance
`
}

// Run executes this subcommand
func (c *BalanceCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if err := c.wallet.MonitorBalance(); err != nil {
		c.ui.Error(fmt.Sprintf("fail to call MonitorBalance() %+v", err))
		return 1
	}

	return 0
}
