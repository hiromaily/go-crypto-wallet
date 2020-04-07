package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	receiptName = "receipt"
)

//receipt subcommand
type ReceiptCommand struct {
	name    string
	version string
	ui      cli.Ui
	wallet  *service.Wallet
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
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(receiptName, flag.ContinueOnError)
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
				wallet:   c.wallet,
			}, nil
		},
		"debug": func() (cli.Command, error) {
			return &DebugSequenceCommand{
				name:     "debug",
				synopsis: debugSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(receiptName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", receiptName, err))
	}
	return code
}
