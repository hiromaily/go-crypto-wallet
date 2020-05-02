package imports

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// PrivKeyCommand privkey subcommand
type PrivKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *PrivKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *PrivKeyCommand) Help() string {
	return `Usage: sign import privkey
`
}

// Run executes this subcommand
func (c *PrivKeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// import generated private key for Authorization account to database
	err := c.wallet.PrivKey().Import()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call ImportPrivateKey() %+v", err))
	}
	c.ui.Output("Done!")

	return 0
}
