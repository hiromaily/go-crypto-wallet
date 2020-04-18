package sign

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

//sign subcommand
type SignCommand struct {
	Name   string
	UI     cli.Ui
	Wallet wallets.Signer
}

func (c *SignCommand) Synopsis() string {
	return "sign on signed transaction for multsig address (account would be found from file name)"
}

func (c *SignCommand) Help() string {
	return `Usage: sign sign [options...]
Options:
  -file  signed transaction file path for multisig address
`
}

func (c *SignCommand) Run(args []string) int {
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

	// sign on signed transactions for multisig, action(receipt/payment) could be found from file name
	hexTx, isSigned, generatedFileName, err := c.Wallet.SignTx(filePath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("fail to call SignTx() %+v", err))
	}

	//TODO: output should be json if json option is true
	c.UI.Output(fmt.Sprintf("[hex]: %s\n[isCompleted]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName))

	return 0
}
