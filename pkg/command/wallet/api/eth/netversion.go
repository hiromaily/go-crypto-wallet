package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// NetVersionCommand syncing subcommand
type NetVersionCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	eth      ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *NetVersionCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *NetVersionCommand) Help() string {
	return `Usage: wallet api netversion
`
}

// Run executes this subcommand
func (c *NetVersionCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	version, err := c.eth.NetVersion()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call eth.NetVersion() %+v", err))
		return 1
	}

	c.ui.Info(fmt.Sprintf("net version: %d", version))

	return 0
}
