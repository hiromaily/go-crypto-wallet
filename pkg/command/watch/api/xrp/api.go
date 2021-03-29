package xrp

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

// APICommand api subcommand
type APICommand struct {
	Name    string
	Version string
	UI      cli.Ui
	XRP     xrpgrp.Rippler
	TxData  *config.RippleTxData
}

// Synopsis is explanation for this subcommand
func (c *APICommand) Synopsis() string {
	return "Ripple API functionality"
}

var sendCoinSynopsis = "send coin from faucet coin"

// Help returns usage for this subcommand
func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet api [Subcommands...]
Subcommands:
  sendcoin    %s
`, sendCoinSynopsis)
}

// Run executes this subcommand
func (c *APICommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"sendcoin": func() (cli.Command, error) {
			return &SendCoinCommand{
				name:     "clientversion",
				synopsis: sendCoinSynopsis,
				ui:       command.ClolorUI(),
				xrp:      c.XRP,
				txData:   c.TxData,
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
