package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//validateaddress subcommand
type ValidateAddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *ValidateAddressCommand) Synopsis() string {
	return c.synopsis
}

func (c *ValidateAddressCommand) Help() string {
	return `Usage: wallet api validateaddress [options...]
Options:
  -address  address like '2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr'
`
}

func (c *ValidateAddressCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		address string
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&address, "address", "", "address")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validate address
	_, err := c.wallet.BTC.ValidateAddress(address)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.ValidateAddress() %+v", err))
		return 1
	}

	return 0
}
