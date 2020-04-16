package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
)

//seed subcommand
type SeedCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

func (c *SeedCommand) Synopsis() string {
	return c.synopsis
}

func (c *SeedCommand) Help() string {
	return `Usage: keygen key create seed  [options...]
Options:
  -seed  given seed is used to store in database instead of generating new seed (development use)
`
}

func (c *SeedCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		seed  string
		bSeed []byte
		err   error
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&seed, "seed", "", "given seed is used to store in database instead of generating new seed (development use)")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if seed != "" {
		// store seed into database, not generate seed
		bSeed, err = c.wallet.StoreSeed(seed)
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call StoreSeed() %+v", err))
			return 1
		}
	} else {
		// create seed
		bSeed, err = c.wallet.GenerateSeed()
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call GenerateSeed() %+v", err))
			return 1
		}
	}
	c.ui.Info(fmt.Sprintf("seed: %s", key.SeedToString(bSeed)))

	return 0
}
