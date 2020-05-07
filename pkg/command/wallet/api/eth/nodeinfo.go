package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// NodeInfoCommand nodeinfo subcommand
type NodeInfoCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	eth      ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *NodeInfoCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *NodeInfoCommand) Help() string {
	return `Usage: wallet api nodeinfo
`
}

// Run executes this subcommand
func (c *NodeInfoCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	peerInfo, err := c.eth.NodeInfo()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call eth.NodeInfo() %+v", err))
		return 1
	}

	c.ui.Info(fmt.Sprintf("nodeinfo: %v", peerInfo))

	return 0
}
