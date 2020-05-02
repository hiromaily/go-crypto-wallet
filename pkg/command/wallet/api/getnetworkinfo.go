package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// GetNetworkInfoCommand getnetworkinfo subcommand
type GetNetworkInfoCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *GetNetworkInfoCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *GetNetworkInfoCommand) Help() string {
	return `Usage: wallet api getnetworkinfo`
}

// Run executes this subcommand
func (c *GetNetworkInfoCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// call getnetworkinfo
	infoData, err := c.wallet.GetBTC().GetNetworkInfo()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.GetNetworkInfo() %+v", err))
		return 1
	}
	grok.Value(infoData)

	return 0
}
