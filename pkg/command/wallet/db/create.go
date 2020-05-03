package db

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// CreateCommand create subcommand
type CreateCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *CreateCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *CreateCommand) Help() string {
	return `Usage: wallet db create [options...]
Options:
  -table  target table name
`
}

// Run executes this subcommand
func (c *CreateCommand) Run(args []string) int {
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
	}
	switch tableName {
	case "payment_request":
		// create payment_request table
		if err := c.wallet.CreatePaymentRequest(); err != nil {
			c.ui.Error(fmt.Sprintf("fail to call CreatePaymentRequest() %+v", err))
		}
		return 1
	}

	return 0
}
