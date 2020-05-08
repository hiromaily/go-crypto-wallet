package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/btcwallet"
)

// KeyCommand key subcommand
type KeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Keygener
}

// Synopsis is explanation for this subcommand
func (c *KeyCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *KeyCommand) Help() string {
	return `Usage: keygen create key
`
}

// Run executes this subcommand
func (c *KeyCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// create one key for debug use
	if v, ok := c.wallet.(*btcwallet.BTCKeygen); ok {
		wif, pubAddress, err := key.GenerateWIF(v.BTC.GetChainConf())
		if err != nil {
			c.ui.Error(fmt.Sprintf("fail to call key.GenerateKey() %+v", err))
			return 1
		}
		c.ui.Info(fmt.Sprintf("[WIF] %s - [Pub Address] %s", wif.String(), pubAddress))
	} else {
		c.ui.Info("for only BTC")
	}
	return 0
}
