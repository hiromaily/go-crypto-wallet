package db

import (
	"github.com/mitchellh/cli"

	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// CreateCommand create subcommand
type CreateCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	wallet   wallets.Watcher
}

// Synopsis is explanation for this subcommand
func (c *CreateCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *CreateCommand) Help() string {
	return `Usage: wallet db create [options...]
Options:
  -table  target table name
`
}

// Run executes this subcommand
func (c *CreateCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	//var (
	//	tableName string
	//)
	//flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	//flags.StringVar(&tableName, "table", "", "table name of database")
	//if err := flags.Parse(args); err != nil {
	//	return 1
	//}
	//
	//c.ui.Output(fmt.Sprintf("-table: %s", tableName))
	//
	////validator
	//if tableName == "" {
	//	tableName = "payment_request"
	//}
	//
	//switch tableName {
	//case "payment_request":
	//	// create payment_request table
	//	amtList := []float64{
	//		0.00001,
	//		0.00002,
	//		0.000025,
	//		0.000015,
	//		0.00003,
	//	}
	//
	//	// get client pubkeys
	//	pubkeyItems, err := c.wallet.GetDB().Addr().GetAll(account.AccountTypeClient)
	//	if err != nil {
	//		c.ui.Error(fmt.Sprintf("fail to call wallet.GetDB().Pubkey().GetAll() %+v", err))
	//		return 1
	//	}
	//	if len(pubkeyItems) < len(amtList)*2 {
	//		c.ui.Error(fmt.Sprintf("pubkey for client should be created at least %d", len(amtList)))
	//	}
	//	// start transaction
	//	dtx, err := c.wallet.GetDB().BeginTx()
	//	if err != nil {
	//		c.ui.Error(fmt.Sprintf("fail to start transaction %+v", err))
	//		return 1
	//	}
	//	defer func() {
	//		if err != nil {
	//			dtx.Rollback()
	//		} else {
	//			dtx.Commit()
	//		}
	//	}()
	//
	//	// delete payment request
	//	_, err = c.wallet.GetDB().PayReq().DeleteAll()
	//	if err != nil {
	//		c.ui.Error(fmt.Sprintf("fail to call wallet.GetDB().PayReq().DeleteAll() %+v", err))
	//		return 1
	//	}
	//	// insert payment_request
	//	payReqItems := make([]*models.PaymentRequest, 0, len(amtList))
	//	for _, amt := range amtList {
	//		payReqItems = append(payReqItems, &models.PaymentRequest{
	//			Coin:            c.wallet.GetBTC().CoinTypeCode().String(),
	//			PaymentID:       null.NewInt64(0, false),
	//			SenderAddress:   pubkeyItems[0].WalletAddress,
	//			SenderAccount:   pubkeyItems[0].Account,
	//			ReceiverAddress: pubkeyItems[len(amtList)].WalletAddress,
	//			Amount:          c.wallet.GetBTC().FloatToDecimal(amt),
	//			IsDone:          false,
	//			UpdatedAt:       null.TimeFrom(time.Now()),
	//		})
	//	}
	//	if err = c.wallet.GetDB().PayReq().InsertBulk(payReqItems); err != nil {
	//		c.ui.Error(fmt.Sprintf("fail to call wallet.GetDB().PayReq().InsertBulk() %+v", err))
	//		return 1
	//	}
	//
	//}

	return 0
}
