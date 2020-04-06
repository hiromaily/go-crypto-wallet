package receipt

import (
	"flag"
	"log"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const receiptName = "receipt"

//receipt subcommand
type ReceiptCommand struct {
	version string
	ui      cli.Ui
	wallet  *service.Wallet
}

func (c *ReceiptCommand) Synopsis() string {
	return "receipt functionality"
}

func (c *ReceiptCommand) Help() string {
	return `Usage: wallet receipt [Subcommands...]
Subcommands:
  create  create receipt transaction file
  find    find receipts transaction from bitcoin blockchain network
  debug   sequences from creating receipts transaction, sign, send signed transaction
`
}

func (c *ReceiptCommand) Run(args []string) int {
	flags := flag.NewFlagSet(receiptName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &CreateTxCommand{
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
		"find": func() (cli.Command, error) {
			return &FindCommand{
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
		"debug": func() (cli.Command, error) {
			return &DebugSequenceCommand{
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(receiptName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() subcommand of receipt: %v\n", err)
	}
	return code
}
