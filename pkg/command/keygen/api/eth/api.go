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
	return "Bitcoin API functionality"
}

var (
	importrawkeySynopsis = "import raw key"
)

// Help returns usage for this subcommand
func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  importrawkey    %s
`, importrawkeySynopsis)
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
		"importrawkey": func() (cli.Command, error) {
			return &ImportRawKeyCommand{
				name:     "importrawkey",
				synopsis: importrawkeySynopsis,
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
