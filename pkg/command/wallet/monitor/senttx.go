package monitor

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// SentTxCommand senttx subcommand
type SentTxCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *SentTxCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *SentTxCommand) Help() string {
	return `Usage: wallet monitor sendtx [options...]
Options:
  -account  target account
`
}

// Run executes this subcommand
func (c *SentTxCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		acnt string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&acnt, "account", "", "account for monitoring")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// monitor sent transactions
	//TODO: add account parameter
	err := c.wallet.UpdateTxStatus()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call UpdateTxStatus() %+v", err))
		return 1
	}

	return 0
}
