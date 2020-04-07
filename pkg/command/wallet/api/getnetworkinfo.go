package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//getnetworkinfo subcommand
type GetnetworkInfoCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *GetnetworkInfoCommand) Synopsis() string {
	return c.synopsis
}

func (c *GetnetworkInfoCommand) Help() string {
	return `Usage: wallet api getnetworkinfo`
}

func (c *GetnetworkInfoCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// call getnetworkinfo
	infoData, err := c.wallet.BTC.GetNetworkInfo()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.GetNetworkInfo() %+v", err))
		return 1
	}
	grok.Value(infoData)

	return 0
}
