package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//find subcommand
type FindCommand struct {
	ui     cli.Ui
	wallet *service.Wallet
}

func (c *FindCommand) Synopsis() string {
	return fmt.Sprintf("%s", findSynopsis)
}

func (c *FindCommand) Help() string {
	return `Usage: wallet receipt find [options...]
Options:
  -fee  adjustment fee
`
}

func (c *FindCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(findName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//TODO: Detect receipt transaction from outside, not create transaction
	return 0
}
