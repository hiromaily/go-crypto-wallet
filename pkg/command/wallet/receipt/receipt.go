package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//receipt subcommand
type ReceiptCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallet.Walleter
}

func (c *ReceiptCommand) Synopsis() string {
	return "receipt functionality"
}

var (
	createSynopsis = "create a receipt transaction file for client account"
	debugSynopsis  = "execute series of flows from creation of a receiving transaction to sending of a transaction"
)

func (c *ReceiptCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet receipt [Subcommands...]
Subcommands:
  create  %s
  debug   %s
`, createSynopsis, debugSynopsis)
}

func (c *ReceiptCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &CreateTxCommand{
				name:     "create",
				synopsis: createSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"debug": func() (cli.Command, error) {
			return &DebugSequenceCommand{
				name:     "debug",
				synopsis: debugSynopsis,
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
