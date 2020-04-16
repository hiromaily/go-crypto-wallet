package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//listunspent subcommand
type ListUnspentCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *ListUnspentCommand) Synopsis() string {
	return c.synopsis
}

func (c *ListUnspentCommand) Help() string {
	return `Usage: wallet api listunspent`
}

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
