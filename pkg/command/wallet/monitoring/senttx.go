package monitoring

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//senttx subcommand
type SentTxCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Walleter
}

func (c *SentTxCommand) Synopsis() string {
	return c.synopsis
}

func (c *SentTxCommand) Help() string {
	return `Usage: wallet monitoring sendtx [options...]
Options:
  -account  target account
`
}

func (c *SentTxCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

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
	err := c.wallet.UpdateStatus()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call UpdateStatus() %+v", err))
		return 1
	}

	return 0
}
