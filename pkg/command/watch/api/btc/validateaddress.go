package btc

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// ValidateAddressCommand validateaddress subcommand
type ValidateAddressCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      btcgrp.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *ValidateAddressCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ValidateAddressCommand) Help() string {
	return `Usage: wallet api validateaddress [options...]
Options:
  -address  address like '2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr'
`
}

// Run executes this subcommand
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
	_, err := c.btc.ValidateAddress(address)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.ValidateAddress() %+v", err))
		return 1
	}

	return 0
}
