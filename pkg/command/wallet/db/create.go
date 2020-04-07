package db

import (
	"flag"
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/testdata"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//create subcommand
type CreateCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *CreateCommand) Synopsis() string {
	return c.synopsis
}

func (c *CreateCommand) Help() string {
	return `Usage: wallet db create [options...]
Options:
  -table  target table name
`
}

func (c *CreateCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

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

	//create payment_request table
	err := testdata.CreateInitialTestData(c.wallet.DB, c.wallet.BTC)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call testdata.CreateInitialTestData() %+v", err))
		return 1
	}
	c.ui.Info("Done!")

	return 0
}
