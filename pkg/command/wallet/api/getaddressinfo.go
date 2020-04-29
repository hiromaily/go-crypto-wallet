package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// GetAddressInfoCommand getaddressinfo subcommand
type GetAddressInfoCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *GetAddressInfoCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *GetAddressInfoCommand) Help() string {
	return `Usage: wallet api getaddressinfo [options...]
Options:
	-address  address
`
}

// Run executes this subcommand
func (c *GetAddressInfoCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var addr string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&addr, "address", "", "address")
	if err := flags.Parse(args); err != nil {
		return 1
	}
	//validator
	if addr == "" {
		c.ui.Error("address option [-address] is required")
		return 1
	}

	// call getaddressinfo
	addrData, err := c.wallet.GetBTC().GetAddressInfo(addr)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.GetAddressInfo() %+v", err))
		return 1
	}
	grok.Value(addrData)

	return 0
}
