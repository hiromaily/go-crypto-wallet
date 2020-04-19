package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ListUnspentCommand listunspent subcommand
type ListUnspentCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *ListUnspentCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ListUnspentCommand) Help() string {
	return `Usage: wallet api listunspent`
}

// Run executes this subcommand
func (c *ListUnspentCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// call listunspent
	unspentList, err := c.wallet.GetBTC().Client().ListUnspentMin(c.wallet.GetBTC().ConfirmationBlock()) //6
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.ListUnspentMin() %+v", err))
		return 1
	}
	grok.Value(unspentList)

	return 0
}
