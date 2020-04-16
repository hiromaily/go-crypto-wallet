package monitor

import (
	"flag"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//balance subcommand
type BalanceCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *BalanceCommand) Synopsis() string {
	return c.synopsis
}

func (c *BalanceCommand) Help() string {
	return `Usage: wallet monitor balance [options...]
Options:
  -account  target account
`
}

func (c *BalanceCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		acnt string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "account for monitoring")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//TODO monitor balance

	return 0
}
