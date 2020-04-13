package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//create subcommand
type CreateCommand struct {
	Name        string
	Version     string
	SynopsisExp string
	UI          cli.Ui
	Wallet      wallets.Signer
}

func (c *CreateCommand) Synopsis() string {
	return c.SynopsisExp
}

var (
	hdkeySynopsis = "create key for hd wallet for Authorization account"
	seedSynopsis  = "create seed"
)

func (c *CreateCommand) Help() string {
	return fmt.Sprintf(`Usage: sign create [Subcommands...]
Subcommands:
  hdkey     %s
  seed      %s
`, hdkeySynopsis, seedSynopsis)
}

func (c *CreateCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"hdkey": func() (cli.Command, error) {
			return &HDKeyCommand{
				name:     "hdkey",
				synopsis: hdkeySynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"seed": func() (cli.Command, error) {
			return &SeedCommand{
				name:     "seed",
				synopsis: seedSynopsis,
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
