package export

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// FullPubkeyCommand multisig subcommand
type FullPubkeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *FullPubkeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *FullPubkeyCommand) Help() string {
	return `Usage: sign export fullpubkey`
}

// Run executes this subcommand
func (c *FullPubkeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// export full pubkey as csv file
	fileName, err := c.wallet.FullPubkeyExport().ExportFullPubkey()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ExportAddedPubkeyHistory() %+v", err))
		return 1
	}
	c.ui.Output(fmt.Sprintf("[fileName]: %s", fileName))

	return 0
}
