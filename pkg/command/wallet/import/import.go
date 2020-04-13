package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//import subcommand
type ImportCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Walleter
}

func (c *ImportCommand) Synopsis() string {
	return "importing functionality"
}

var (
	keySynopsis = "import generatd addresses by keygen wallet"
)

func (c *ImportCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet import [Subcommands...]
Subcommands:
  key  %s
`, keySynopsis)
}

func (c *ImportCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"key": func() (cli.Command, error) {
			return &KeyCommand{
				name:     "key",
				synopsis: keySynopsis,
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
