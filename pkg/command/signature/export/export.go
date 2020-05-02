package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ExportCommand export subcommand
type ExportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *ExportCommand) Synopsis() string {
	return "export resources"
}

var (
	fullpubkeySynopsis = "export full pubkey"
)

// Help returns usage for this subcommand
func (c *ExportCommand) Help() string {
	return fmt.Sprintf(`Usage: sign export [Subcommands...]
Subcommands:
  fullpubkey   %s
`, fullpubkeySynopsis)
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
		"fullpubkey": func() (cli.Command, error) {
			return &FullPubkeyCommand{
				name:     "fullpubkey",
				synopsis: fullpubkeySynopsis,
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
