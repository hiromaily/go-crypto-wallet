package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//logging subcommand
type LoggingCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   *service.Wallet
}

func (c *LoggingCommand) Synopsis() string {
	return c.synopsis
}

func (c *LoggingCommand) Help() string {
	return `Usage: wallet api logging`
}

func (c *LoggingCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// logging
	logData, err := c.wallet.BTC.Logging()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.Logging() %+v", err))
		return 1
	}
	grok.Value(logData)

	return 0
}
