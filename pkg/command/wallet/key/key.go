package key

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//key subcommand
type KeyCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Walleter
}

func (c *KeyCommand) Synopsis() string {
	return "key importing functionality"
}

var (
	importSynopsis = "import generatd addresses by keygen wallet"
)

func (c *KeyCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet key [Subcommands...]
Subcommands:
  import  %s
`, importSynopsis)
}

func (c *KeyCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"import": func() (cli.Command, error) {
			return &ImportCommand{
				name:     "import",
				synopsis: importSynopsis,
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
