package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//create subcommand
type CreateCommand struct {
	Name        string
	Version     string
	SynopsisExp string
	UI          cli.Ui
	Wallet      wallet.Keygener
}

func (c *CreateCommand) Synopsis() string {
	return c.SynopsisExp
}

var (
	keySynopsis  = "create key"
	seedSynopsis = "create seed"
)

func (c *CreateCommand) Help() string {
	return fmt.Sprintf(`Usage: keygen create [Subcommands...]
Subcommands:
  key   %s
  seed  %s
`, keySynopsis, seedSynopsis)
}

func (c *CreateCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &KeyCommand{
				name:     "create",
				synopsis: keySynopsis,
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
