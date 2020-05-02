package create

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// HDKeyCommand hd key subcommand
type HDKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Signer
}

// Synopsis is explanation for this subcommand
func (c *HDKeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *HDKeyCommand) Help() string {
	return `Usage: sign create hdkey
`
}

// Run executes this subcommand
func (c *HDKeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// create seed
	bSeed, err := c.wallet.Seed().Generate()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call GenerateSeed() %+v", err))
		return 1
	}

	// create key for hd wallet for Authorization account
	keys, err := c.wallet.HDWallet().Generate(c.wallet.GetAuthType().AccountType(), bSeed, 1)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call GenerateAccountKey() %+v", err))
		return 1
	}
	grok.Value(keys)

	return 0
}
