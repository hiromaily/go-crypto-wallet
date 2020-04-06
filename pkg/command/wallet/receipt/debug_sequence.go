package receipt

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

const debugName = "debug"

//debug subcommand
type DebugSequenceCommand struct {
	ui     cli.Ui
	wallet *service.Wallet
}

func (c *DebugSequenceCommand) Synopsis() string {
	return "detect receipt to our addressed and crate receipt unsigned transaction"
}

func (c *DebugSequenceCommand) Help() string {
	return `Usage: wallet receipt create [options...]
Options:
  -fee  adjustment fee
`
}

func (c *DebugSequenceCommand) Run(args []string) int {
	var (
		fee float64
	)
	flags := flag.NewFlagSet(debugName, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	c.ui.Output(fmt.Sprintf("-fee: %f", fee))

	// make sure sequence from detecting receipt transactions, sign on unsigned transaction, send signed transaction
	// 1.Detect receipt transaction from outside and create receipt unsigned transaction
	c.ui.Info("[1]Run: Detect receipt transaction")
	hex, fileName, err := c.wallet.DetectReceivedCoin(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call DetectReceivedCoin() %+v", err))
		return 1
	}
	if hex == "" {
		c.ui.Info("No utxo")
		return 0
	}
	c.ui.Output(fmt.Sprintf("[hex]: %s\n[fileName]: %s", hex, fileName))

	// 2. sign on unsigned transaction. actually it is available for sign wallet
	c.ui.Info("\n[2]Run: sign")
	hexTx, isSigned, generatedFileName, err := c.wallet.SignatureFromFile(fileName)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call SignatureFromFile() %+v", err))
		return 1
	}
	c.ui.Output(fmt.Sprintf("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName))

	// 3. send signed transaction to blockchain network
	c.ui.Info("\n[3]Run: send signed transaction")
	txID, err := c.wallet.SendFromFile(generatedFileName)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call SendFromFile() %+v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("[Done] sent transaction ID:%s", txID))

	return 0
}
