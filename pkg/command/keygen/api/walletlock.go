package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
)

// WalletLockCommand walletlock subcommand
type WalletLockCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      btcgrp.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *WalletLockCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *WalletLockCommand) Help() string {
	return `Usage: keygen api walletlock`
}

// Run executes this subcommand
func (c *WalletLockCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	err := c.btc.WalletLock()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call WalletLock() %+v", err))
		return 1
	}

	c.ui.Info("wallet is locked!")

	return 0
}
