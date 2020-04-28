package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ImportCommand import subcommand
type ImportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Walleter
}

// Synopsis is explanation for this subcommand
func (c *ImportCommand) Synopsis() string {
	return "importing functionality"
}

var (
	addressSynopsis = "import generatd addresses by keygen wallet"
)

// Help returns usage for this subcommand
func (c *ImportCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet import [Subcommands...]
Subcommands:
  address  %s
`, addressSynopsis)
}

// Run executes this subcommand
func (c *ImportCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"address": func() (cli.Command, error) {
			return &AddressCommand{
				name:     "address",
				synopsis: addressSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
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
