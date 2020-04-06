package receipt

import (
	"flag"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

const findName = "find"

//find subcommand
type FindCommand struct {
	ui     cli.Ui
	wallet *wallet.Wallet
}

func (c *FindCommand) Synopsis() string {
	return "detect receipt to our addressed and crate receipt unsigned transaction"
}

func (c *FindCommand) Help() string {
	return `Usage: wallet receipt create [options...]
Options:
  -fee  adjustment fee
`
}

func (c *FindCommand) Run(args []string) int {
	flags := flag.NewFlagSet(findName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//TODO: Detect receipt transaction from outside, not create transaction
	return 0
}
