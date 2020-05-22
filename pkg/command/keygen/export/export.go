package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// ExportCommand export subcommand
type ExportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *ExportCommand) Synopsis() string {
	return "export resources"
}

var (
	addressSynopsis = "export generated PublicKey as csv file"
)

// Help returns usage for this subcommand
func (c *ExportCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen export [Subcommands...]
Subcommands:
  address   %s
`, addressSynopsis)
}

// Run executes this subcommand
func (c *ExportCommand) Run(args []string) int {
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
