package monitoring

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	montoringName = "montoring"
)

//montoring subcommand
type MontoringCommand struct {
	name    string
	version string
	ui      cli.Ui
	wallet  *service.Wallet
}

func (c *MontoringCommand) Synopsis() string {
	return "montoring functionality"
}

var (
	senttxSynopsis  = "monitor sent transactions"
	balanceSynopsis = "monitor balance"
)

func (c *MontoringCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet receipt [Subcommands...]
Subcommands:
  senttx   %s
  balance  %s
`, senttxSynopsis, balanceSynopsis)
}

func (c *MontoringCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(montoringName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"senttx": func() (cli.Command, error) {
			return &SentTxCommand{
				name:   "senttx",
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
		"balance": func() (cli.Command, error) {
			return &BalanceCommand{
				name:   "balance",
				ui:     command.ClolorUI(),
				wallet: c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(montoringName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of monitoring: %v", err))
	}
	return code
}
