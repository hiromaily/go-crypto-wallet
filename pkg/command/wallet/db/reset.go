package db

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// ResetCommand reset subcommand
type ResetCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *ResetCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *ResetCommand) Help() string {
	return `Usage: wallet db reset [options...]
Options:
  -table  target table name
`
}

// Run executes this subcommand
func (c *ResetCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		tableName string
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&tableName, "table", "", "table name of database")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-table: %s", tableName))

	//validator
	if tableName == "" {
		tableName = "payment_request"
		//c.ui.Error("table name option [-table] is required")
		//return 1
	}

	//reset payment_request table
	//FIXME: ResetTestData is under refactoring
	c.ui.Info("FIXME: ResetTestData is under refactoring")

	return 0
}
