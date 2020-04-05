package wallet

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

const importName = "import"

//import subcommand
type ImportCommand struct {
	ui     cli.Ui
	wallet *wallet.Wallet
}

func (c *ImportCommand) Synopsis() string {
	return "key importing functionality"
}

func (c *ImportCommand) Help() string {
	return `Usage: wallet key import [options...]
Options:
  -file  import file path for generated addresses
`
}

func (c *ImportCommand) Run(args []string) int {
	var (
		filePath string
	)
	flags := flag.NewFlagSet(importName, flag.ContinueOnError)
	flags.StringVar(&filePath, "file", "", "import file path for generated addresses")
	if err := flags.Parse(args); err != nil {
		return 1
	}
	c.ui.Output(fmt.Sprintf("-file: %s", filePath))

	//validator
	if filePath == "" {
		c.ui.Error("file path option [-file] is required")
	}

	//logic

	return 0
}
