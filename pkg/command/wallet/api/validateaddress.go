package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ValidateAddressCommand validateaddress subcommand
type ValidateAddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

// Synopsis
func (c *ValidateAddressCommand) Synopsis() string {
	return c.synopsis
}

// Help
func (c *ValidateAddressCommand) Help() string {
	return `Usage: wallet api validateaddress [options...]
Options:
  -address  address like '2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr'
`
}

// Run
func (c *ValidateAddressCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		address string
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&address, "address", "", "address")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validate args
	if address == "" {
		c.ui.Error("address option [-address] is required")
		return 1
	}

	// validate address
	_, err := c.wallet.GetBTC().ValidateAddress(address)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.ValidateAddress() %+v", err))
		return 1
	}

	return 0
}
