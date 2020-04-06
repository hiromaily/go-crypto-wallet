package payment

import (
	"flag"
	"fmt"
	"log"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	paymentName = "paymentName"
	createName  = "create"
	debugName   = "debug"
)

//payment subcommand
type PaymentCommand struct {
	version string
	ui      cli.Ui
	wallet  *service.Wallet
}

func (c *PaymentCommand) Synopsis() string {
	return "payment functionality"
}

var (
	createSynopsis = "create a payment transaction file"
	debugSynopsis  = "execute series of flows from creation of a payment transaction to sending of a transaction"
)

func (c *PaymentCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet payment [Subcommands...]
Subcommands:
  create  %s
  debug   %s
`, createSynopsis, debugSynopsis)
}

func (c *PaymentCommand) Run(args []string) int {
	flags := flag.NewFlagSet(paymentName, flag.ContinueOnError)
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
		"debug": func() (cli.Command, error) {
			return &DebugSequenceCommand{
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(paymentName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() subcommand of payment: %v\n", err)
	}
	return code
}
