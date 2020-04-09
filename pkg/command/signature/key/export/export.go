package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

// export subcommand
type ExportCommand struct {
	Name        string
	Version     string
	SynopsisExp string
	UI          cli.Ui
	Wallet      wallet.Signer
}

func (c *ExportCommand) Synopsis() string {
	return c.SynopsisExp
}

var (
	multisigSynopsis = ""
)

func (c *ExportCommand) Help() string {
	return fmt.Sprintf(`Usage: sign export [Subcommands...]
Subcommands:
  multisig     %s
`, multisigSynopsis)
}

func (c *ExportCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"multisig": func() (cli.Command, error) {
			return &MultisigCommand{
				name:     "multisig",
				synopsis: multisigSynopsis,
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
