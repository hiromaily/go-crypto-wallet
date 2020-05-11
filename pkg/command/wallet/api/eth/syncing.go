package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// SyncingCommand syncing subcommand
type SyncingCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	eth      ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *SyncingCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *SyncingCommand) Help() string {
	return `Usage: wallet api syncing
`
}

// Run executes this subcommand
func (c *SyncingCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	syncResult, isSyncing, err := c.eth.Syncing()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call eth.Syncing() %+v", err))
		return 1
	}

	c.ui.Info(fmt.Sprintf("is syncing? : %t, %v", isSyncing, syncResult))

	return 0
}
