package api

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// ImportWalletCommand importwallet subcommand
type ImportWalletCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *ImportWalletCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ImportWalletCommand) Help() string {
	return `Usage: keygen api dumpwallet [options...]
Options:
  -file  filename
`
}

// Run executes this subcommand
func (c *ImportWalletCommand) Run(args []string) int {
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

	err := c.btc.ImportWallet(fileName)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call btc.ImportWallet() %+v", err))
		return 1
	}

	c.ui.Info("wallet file is imported!")

	return 0
}
