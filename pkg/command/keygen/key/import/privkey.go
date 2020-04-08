package _import

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

//privkey subcommand
type PrivKeyCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallet.Keygener
}

func (c *PrivKeyCommand) Synopsis() string {
	return c.synopsis
}

func (c *PrivKeyCommand) Help() string {
	return `Usage: keygen key import privkey [options...]
Options:
  -table  target table name
`
}

func (c *PrivKeyCommand) Run(args []string) int {
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

	////validator
	//if tableName == "" {
	//	tableName = "payment_request"
	//	//c.ui.Error("table name option [-table] is required")
	//	//return 1
	//}
	//
	////create payment_request table
	//err := testdata.CreateInitialTestData(c.wallet.GetDB(), c.wallet.GetBTC())
	//if err != nil {
	//	c.ui.Error(fmt.Sprintf("fail to call testdata.CreateInitialTestData() %+v", err))
	//	return 1
	//}
	//c.ui.Info("Done!")

	return 0
}
