package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

//TODO: this code is almost same to keygen wallet

//seed subcommand
type SeedCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Signer
}

func (c *SeedCommand) Synopsis() string {
	return c.synopsis
}

func (c *SeedCommand) Help() string {
	return `Usage: sign key create seed
`
}

func (c *SeedCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// create seed
	bSeed, err := c.wallet.GenerateSeed()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call GenerateSeed() %+v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("seed: %s", key.SeedToString(bSeed)))

	return 0
}
