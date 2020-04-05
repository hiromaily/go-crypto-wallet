package command

import (
	"os"

	"github.com/mitchellh/cli"
)

func CreateSubCommand(appName, version string, args []string, commands map[string]cli.CommandFactory) *cli.CLI {
	c := cli.NewCLI(appName, version)
	c.Args = args

	//register subcommand
	c.Commands = commands
	return c
}

func ClolorUI() *cli.ColoredUi {
	return &cli.ColoredUi{
		InfoColor:  cli.UiColorBlue,
		ErrorColor: cli.UiColorRed,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
			Reader:      os.Stdin,
		},
	}
}
