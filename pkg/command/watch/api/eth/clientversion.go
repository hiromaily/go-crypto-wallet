package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// ClientVersionCommand syncing subcommand
type ClientVersionCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	eth      ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *ClientVersionCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ClientVersionCommand) Help() string {
	return `Usage: wallet api clientversion
`
}

// Run executes this subcommand
func (c *ClientVersionCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	version, err := c.eth.ClientVersion()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call eth.ClientVersion() %+v", err))
		return 1
	}

	c.ui.Info(fmt.Sprintf("client version: %s", version))

	return 0
}
