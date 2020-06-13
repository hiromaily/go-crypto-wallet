package eth

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
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
	return "Ethereum API functionality"
}

var (
	clientVersionSynopsis  = "network version"
	nodeinfoSynopsis       = "node info"
	syncingSynopsis        = "sync info"
	networkVersionSynopsis = "network version"
)

// Help returns usage for this subcommand
func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  clientversion    %s
  nodeinfo         %s
  syncing          %s
  netversion       %s
`, clientVersionSynopsis, nodeinfoSynopsis, syncingSynopsis, networkVersionSynopsis)
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
		"clientversion": func() (cli.Command, error) {
			return &ClientVersionCommand{
				name:     "clientversion",
				synopsis: clientVersionSynopsis,
				ui:       command.ClolorUI(),
				eth:      c.ETH,
			}, nil
		},
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
		"netversion": func() (cli.Command, error) {
			return &NetVersionCommand{
				name:     "netversion",
				synopsis: networkVersionSynopsis,
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
