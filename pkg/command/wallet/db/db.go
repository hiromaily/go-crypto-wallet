package db

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	dbName = "db"
)

//db subcommand
type APICommand struct {
	name    string
	version string
	ui      cli.Ui
	wallet  *service.Wallet
}

func (c *APICommand) Synopsis() string {
	return "Database functionality"
}

var (
	createSynopsis = "create table"
	resetSynopsis  = "reset table"
)

func (c *APICommand) Help() string {
	return fmt.Sprintf(`Usage: wallet db [Subcommands...]
Subcommands:
  create  %s
  reset  %s
`, createSynopsis, resetSynopsis)
}

func (c *APICommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	flags := flag.NewFlagSet(dbName, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//farther subcommand import
	cmds := map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &CreateCommand{
				name:     "create",
				synopsis: createSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
		"reset": func() (cli.Command, error) {
			return &ResetCommand{
				name:     "reset",
				synopsis: resetSynopsis,
				ui:       command.ClolorUI(),
				wallet:   c.wallet,
			}, nil
		},
	}
	cl := command.CreateSubCommand(dbName, c.version, args, cmds)

	code, err := cl.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call Run() subcommand of %s: %v", dbName, err))
	}
	return code
}
