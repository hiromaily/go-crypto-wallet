package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
)

// WalletPassphraseChangeCommand walletpassphrasechange subcommand
type WalletPassphraseChangeCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      btcgrp.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *WalletPassphraseChangeCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *WalletPassphraseChangeCommand) Help() string {
	return `Usage: keygen api walletpassphrasechange [options...]
Options:
  -old  old passphrase
  -new  new passphrase
`
}

// Run executes this subcommand
func (c *WalletPassphraseChangeCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var old, new string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&old, "old", "", "old passphrase")
	flags.StringVar(&new, "new", "", "new passphrase")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if old == "" {
		c.ui.Error("old passphrase option [-old] is required")
		return 1
	}
	if new == "" {
		c.ui.Error("new passphrase option [-new] is required")
		return 1
	}

	err := c.btc.WalletPassphraseChange(old, new)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call btc.WalletPassphraseChange() %+v", err))
		return 1
	}

	c.ui.Info("wallet passphrase was changed!")

	return 0
}
