package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

// CreateSubCommand creates sub command calling cli.NewCLI
func CreateSubCommand(appName, version string, args []string, commands map[string]cli.CommandFactory) *cli.CLI {
	c := cli.NewCLI(appName, version)
	c.Args = args

	//register subcommand
	c.Commands = commands
	return c
}

// ClolorUI returns cli.ColoredUi
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

// HelpFunc modify default help explanation
//  if new args are added on the top, modify this func
func HelpFunc(appName string) func(c map[string]cli.CommandFactory) string {
	return func(c map[string]cli.CommandFactory) string {
		// Replace basic help header by new one
		// because it doesn't show optional flags.
		header := fmt.Sprintf(
			"Usage: %s [-version] [-help] [-conf] <command> [<args>]",
			appName)
		s := cli.BasicHelpFunc(appName)(c)
		i := strings.Index(s, "\n")
		s = strings.Replace(s, s[:i], header, 1)
		return s
	}
}

func SearchArg(key string) bool {
	for _, v := range os.Args[1:] {
		if v == key {
			return true
		}
	}
	return false
}

func SearchArgs(keys []string) bool {
	for _, v := range os.Args[1:] {
		for _, key := range keys {
			if v == key {
				return true
			}
		}
	}
	return false
}
