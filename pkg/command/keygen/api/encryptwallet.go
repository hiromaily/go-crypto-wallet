package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// EncryptWalletCommand encryptwallet subcommand
type EncryptWalletCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *EncryptWalletCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *EncryptWalletCommand) Help() string {
	return `Usage: keygen api encryptwallet [options...]
Options:
  -passphrase  passphrase
`
}

// Run executes this subcommand
func (c *EncryptWalletCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var passphrase string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&passphrase, "passphrase", "", "passphrase")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if passphrase == "" {
		c.ui.Error("passphrase option [-passphrase] is required")
		return 1
	}

	err := c.btc.EncryptWallet(passphrase)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call btc.EncryptWallet() %+v", err))
		return 1
	}

	c.ui.Info("wallet is encrypted!")

	return 0
}
