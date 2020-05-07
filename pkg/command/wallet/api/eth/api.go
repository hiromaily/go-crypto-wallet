package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// APICommand api subcommand
type APICommand struct {
	Name    string
	Version string
	UI      cli.Ui
	ETH     ethgrp.Ethereumer
}

// Synopsis is explanation for this subcommand
func (c *APICommand) Synopsis() string {
	return "Bitcoin API functionality"
}

var (
	nodeinfoSynopsis = "node info"
	syncingSynopsis  = "sync info"
)

// Help returns usage for this subcommand
func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  nodeinfo         %s
  syncing          %s
`, nodeinfoSynopsis, syncingSynopsis)
}

// Run executes this subcommand
func (c *APICommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"nodeinfo": func() (cli.Command, error) {
			return &NodeInfoCommand{
				name:     "nodeinfo",
				synopsis: nodeinfoSynopsis,
				ui:       command.ClolorUI(),
				eth:      c.ETH,
			}, nil
		},
		"syncing": func() (cli.Command, error) {
			return &SyncingCommand{
				name:     "syncing",
				synopsis: syncingSynopsis,
				ui:       command.ClolorUI(),
				eth:      c.ETH,
			}, nil
		},
	}
	cl := command.CreateSubCommand(c.Name, c.Version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", c.Name, err))
	}
	return code
}
