package api

import (
	"flag"
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// LoggingCommand logging subcommand
type LoggingCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	btc      api.Bitcoiner
}

// Synopsis is explanation for this subcommand
func (c *LoggingCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *LoggingCommand) Help() string {
	return `Usage: wallet api logging`
}

// Run executes this subcommand
func (c *LoggingCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// logging
	logData, err := c.btc.Logging()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call BTC.Logging() %+v", err))
		return 1
	}
	grok.Value(logData)

	return 0
}
