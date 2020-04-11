package monitoring

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//montoring subcommand
type MonitoringCommand struct {
	Name    string
	Version string
	UI      cli.Ui
	Wallet  wallets.Walleter
}

func (c *MonitoringCommand) Synopsis() string {
	return "montoring functionality"
}

var (
	senttxSynopsis  = "monitor sent transactions"
	balanceSynopsis = "monitor balance"
)

func (c *MonitoringCommand) Help() string {
	return fmt.Sprintf(`Usage: wallet receipt [Subcommands...]
Subcommands:
  senttx   %s
  balance  %s
`, senttxSynopsis, balanceSynopsis)
}

func (c *MonitoringCommand) Run(args []string) int {
	c.UI.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"senttx": func() (cli.Command, error) {
			return &SentTxCommand{
				name:     "senttx",
				synopsis: senttxSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.Wallet,
			}, nil
		},
		"balance": func() (cli.Command, error) {
			return &BalanceCommand{
				name:     "balance",
				synopsis: balanceSynopsis,
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
