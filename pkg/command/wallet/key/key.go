package key

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const keyName = "key"

//key subcommand
type KeyCommand struct {
	name    string
	version string
	ui      cli.Ui
	wallet  *service.Wallet
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
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(keyName, flag.ContinueOnError)
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
				wallet:   c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(keyName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", keyName, err))
	}
	return code
}
