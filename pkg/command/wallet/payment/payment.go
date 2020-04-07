package payment

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	paymentName = "payment"
)

//payment subcommand
type PaymentCommand struct {
	name    string
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
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(paymentName, flag.ContinueOnError)
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
	cl := command.CreateSubCommand(paymentName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", paymentName, err))
	}
	return code
}
