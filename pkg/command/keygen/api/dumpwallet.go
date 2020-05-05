package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// DumpWalletCommand balance subcommand
type DumpWalletCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *DumpWalletCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *DumpWalletCommand) Help() string {
	return `Usage: keygen api dumpwallet [options...]
Options:
  -file  filename
`
}

// Run executes this subcommand
func (c *DumpWalletCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var fileName string

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&fileName, "file", "", "file name")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if fileName == "" {
		c.ui.Error("finename option [-file] is required")
		return 1
	}

	err := c.btc.DumpWallet(fileName)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call btc.DumpWallet() %+v", err))
		return 1
	}

	c.ui.Info("wallet file is dumped!")

	return 0
}
