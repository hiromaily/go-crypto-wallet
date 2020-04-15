package create

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallets"
)

// TODO
//  - how to display help? upper layer's help displays by `wallet receipt create -h`
//  - as workaround, add undefined flag like `wallet receipt create -a`

//receipt subcommand
type ReceiptCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Walleter
}

func (c *ReceiptCommand) Synopsis() string {
	return c.synopsis
}

func (c *ReceiptCommand) Help() string {
	return `Usage: wallet create receipt [options...]
Options:
  -fee    adjustment fee
  -debug  execute series of flows from creation of a receiving transaction to sending of a transaction
`
}

func (c *ReceiptCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		fee float64
		//isCheck bool
		isDebug bool
	)
	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.Float64Var(&fee, "fee", 0, "adjustment fee")
	//flags.BoolVar(&isCheck, "check", false, "only check client addresses, not create unsigned transaction")
	flags.BoolVar(&isDebug, "debug", false, "execute series of flows from creation of a receiving transaction to sending of a transaction")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if isDebug {
		return c.runDebug(fee)
	}

	// Detect transaction for clients from blockchain network and create receipt unsigned transaction
	// It would be run manually on the daily basis because signature is manual task
	hex, fileName, err := c.wallet.DetectReceivedCoin(fee)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call DetectReceivedCoin() %+v", err))
		return 1
	}
	if hex == "" {
		c.ui.Info("No utxo")
		return 0
	}
	//TODO: output should be json if json option is true
	c.ui.Output(fmt.Sprintf("[hex]: %s\n[fileName]: %s", hex, fileName))

	return 0
}

func (c *ReceiptCommand) runDebug(fee float64) int {
	c.ui.Output("debug mode")

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

	//FIXME: no SignTx in walleter interface
	// 2. sign on unsigned transaction. actually it is available for sign wallet
	//c.ui.Info("\n[2]Run: sign")
	//hexTx, isSigned, generatedFileName, err := c.wallet.SignTx(fileName)
	//if err != nil {
	//	c.ui.Error(fmt.Sprintf("fail to call SignTx() %+v", err))
	//	return 1
	//}
	//c.ui.Output(fmt.Sprintf("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName))

	// 3. send signed transaction to blockchain network
	//c.ui.Info("\n[3]Run: send signed transaction")
	//txID, err := c.wallet.SendFromFile(generatedFileName)
	//if err != nil {
	//	c.ui.Error(fmt.Sprintf("fail to call SendFromFile() %+v", err))
	//	return 1
	//}
	//c.ui.Info(fmt.Sprintf("[Done] sent transaction ID:%s", txID))

	return 0
}
