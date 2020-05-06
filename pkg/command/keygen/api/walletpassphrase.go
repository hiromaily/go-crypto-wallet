package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
)

// WalletPassphraseCommand walletpassphrase subcommand
type WalletPassphraseCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      btcgrp.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *WalletPassphraseCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *WalletPassphraseCommand) Help() string {
	return `Usage: keygen api walletpassphrase [options...]
Options:
  -passphrase  passphrase
`
}

// Run executes this subcommand
func (c *WalletPassphraseCommand) Run(args []string) int {
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

	err := c.btc.WalletPassphrase(passphrase, 10)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call btc.WalletPassphrase() %+v", err))
		return 1
	}

	c.ui.Info("wallet encryption is unlocked for 10s!")

	return 0
}
