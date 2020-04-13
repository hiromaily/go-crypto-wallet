package send

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

//send subcommand
type SendCommand struct {
	Name   string
	UI     cli.Ui
	Wallet wallets.Walleter
}

func (c *SendCommand) Synopsis() string {
	return "send signed transaction to bitcoin blockchain network"
}

func (c *SendCommand) Help() string {
	return `Usage: wallet send [options...]
Options:
  -file  signed transaction file path
`
}

func (c *SendCommand) Run(args []string) int {
	c.UI.Info(c.Synopsis())

	var (
		filePath string
	)
	flags := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	flags.StringVar(&filePath, "file", "", "import file path for signed transactions")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if filePath == "" {
		c.UI.Error("file path option [-file] is required")
		return 1
	}

	// send signed transactions
	txID, err := c.Wallet.SendFromFile(filePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call SendFromFile() %+v", err))
		return 1
	}

	//TODO: output should be json if json option is true
	c.UI.Output(fmt.Sprintf("[Done]送信までDONE!! txID: %s", txID))

	return 0
}
