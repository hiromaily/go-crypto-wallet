package sign

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//sending subcommand
type SignatureCommand struct {
	Name   string
	UI     cli.Ui
	Wallet wallets.Keygener
}

func (c *SignatureCommand) Synopsis() string {
	return "sign on unsigned transaction (account would be found from file name)"
}

func (c *SignatureCommand) Help() string {
	return `Usage: wallet sending [options...]
Options:
  -file  unsigned transaction file path
`
}

func (c *SignatureCommand) Run(args []string) int {
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

	// sign on unsigned transactions, action(receipt/payment) could be found from file name
	hexTx, isSigned, generatedFileName, err := c.Wallet.SignTx(filePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call SignTx() %+v", err))
	}

	//TODO: output should be json if json option is true
	c.UI.Output(fmt.Sprintf("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName))

	return 0
}
