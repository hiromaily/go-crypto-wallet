package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//MultisigCommand  multisig subcommand
type MultisigCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *MultisigCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *MultisigCommand) Help() string {
	return `Usage: keygen key create multisig
`
}

// Run executes this subcommand
func (c *MultisigCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// create multisig address for debug use
	resAddr, err := c.wallet.GetBTC().AddMultisigAddress(
		2,
		[]string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"},
		"multi01",
		address.AddrTypeP2shSegwit)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call AddMultisigAddress() %+v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript))

	return 0
}
