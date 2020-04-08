package db

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//reset subcommand
type ResetCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *ResetCommand) Synopsis() string {
	return c.synopsis
}

func (c *ResetCommand) Help() string {
	return `Usage: wallet db reset [options...]
Options:
  -table  target table name
`
}

func (c *ResetCommand) Run(args []string) int {
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

	//reset payment_request table
	_, err := c.wallet.DB.ResetAnyFlagOnPaymentRequestForTestOnly(nil, true)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call db.ResetAnyFlagOnPaymentRequestForTestOnly() %+v", err))
		return 1
	}
	c.ui.Info("Done!")

	return 0
}
