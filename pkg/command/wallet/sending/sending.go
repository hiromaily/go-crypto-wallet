package sending

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const (
	sendingName = "sending"
)

//sending subcommand
type SendingCommand struct {
	name   string
	ui     cli.Ui
	wallet *service.Wallet
}

func (c *SendingCommand) Synopsis() string {
	return "send signed transaction to bitcoin blockchain network"
}

func (c *SendingCommand) Help() string {
	return `Usage: wallet sending [options...]
Options:
  -file  signed transaction file path
`
}

//WIP
func (c *SendingCommand) Run(args []string) int {
	c.ui.Output(c.Synopsis())

	var (
		filePath string
	)
	flags := flag.NewFlagSet(sendingName, flag.ContinueOnError)
	flags.StringVar(&filePath, "file", "", "import file path for signed transactions")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	//validator
	if filePath == "" {
		c.ui.Error("file path option [-file] is required")
		return 1
	}

	// send signed transactions
	txID, err := c.wallet.SendFromFile(filePath)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call SendFromFile() %+v", err))
	}

	//TODO: output should be json if json option is true
	c.ui.Output(fmt.Sprintf("[Done]送信までDONE!! txID: %s", txID))

	return 0
}
