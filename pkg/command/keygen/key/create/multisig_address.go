package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

// multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Keygener
}

func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

func (c *MultisigCommand) Help() string {
	return `Usage: keygen key create multisig
`
}

func (c *MultisigCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// create multisig address for debug use
	resAddr, err := c.wallet.GetBTC().AddMultisigAddress(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01", enum.AddressTypeP2shSegwit)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call AddMultisigAddress() %+v", err))
	}
	c.ui.Info(fmt.Sprintf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript))

	return 0
}