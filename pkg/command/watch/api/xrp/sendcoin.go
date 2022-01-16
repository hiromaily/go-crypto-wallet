package xrp

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/mitchellh/cli"
	"google.golang.org/grpc/status"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	pb "github.com/hiromaily/ripple-lib-proto/v2/pb/go/rippleapi"
)

// SendCoinCommand syncing subcommand
type SendCoinCommand struct {
	name     string
	synopsis string
	ui       cli.Ui
	xrp      xrpgrp.Rippler
	txData   *config.RippleTxData
}

// Synopsis is explanation for this subcommand
func (c *SendCoinCommand) Synopsis() string {
	return c.synopsis
}

// Help returns usage for this subcommand
func (c *SendCoinCommand) Help() string {
	return `Usage: wallet api sendcoin [options...]
Options:
  -address  receiver address
  -amount   amount
`
}

// Run executes this subcommand
func (c *SendCoinCommand) Run(args []string) int {
	c.ui.Info(c.Synopsis())

	var (
		receiverAddr string
		amount       float64
	)

	flags := flag.NewFlagSet(c.name, flag.ContinueOnError)
	flags.StringVar(&receiverAddr, "address", "", "receiver address")
	flags.Float64Var(&amount, "amount", 0, "amount")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// validator
	if receiverAddr == "" {
		c.ui.Error("address option [-address] is invalid")
		return 1
	}

	// send coin
	// PrepareTransaction
	instructions := &pb.Instructions{
		MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
	}
	c.ui.Info(fmt.Sprintf("sender: %s, receiver: %s, amount: %v", c.txData.Account, receiverAddr, amount))
	txJSON, _, err := c.xrp.CreateRawTransaction(c.txData.Account, receiverAddr, amount, instructions)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call xrp.CreateRawTransaction() %v", err))
		return 1
	}
	grok.Value(txJSON)

	// SingTransaction
	txID, txBlob, err := c.xrp.SignTransaction(txJSON, c.txData.Secret)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call xrp.SignTransaction() %v", err))
		return 1
	}

	// SendTransaction
	sentTx, earlistLedgerVersion, err := c.xrp.SubmitTransaction(txBlob)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call xrp.SubmitTransaction() %v", err))
		return 1
	}
	if strings.Contains(sentTx.ResultCode, "UNFUNDED_PAYMENT") {
		c.ui.Error(fmt.Sprintf("fail to call SubmitTransaction. resultCode: %s, resultMessage: %s", sentTx.ResultCode, sentTx.ResultMessage))
		return 1
	}

	// validate transaction
	_, err = c.xrp.WaitValidation(sentTx.TxJSON.LastLedgerSequence)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call xrp.WaitValidation() %v", err))
		return 1
	}

	// get transaction info
	txInfo, err := c.xrp.GetTransaction(txID, earlistLedgerVersion)
	if err != nil {
		c.ui.Error(fmt.Sprintf("fail to call xrp.GetTransaction() %v", err))
		return 1
	}
	c.ui.Info(fmt.Sprintf("transaction Info: %v", txInfo))

	// get receiver info
	accountInfo, err := c.xrp.GetAccountInfo(receiverAddr)
	if err != nil {
		errStatus, _ := status.FromError(err)
		c.ui.Error(fmt.Sprintf("fail to call xrp.GetAccountInfo() code: %d, message: %s", errStatus.Code(), errStatus.Message()))
		return 1
	}
	c.ui.Info(fmt.Sprintf("receiver account Info: %v", accountInfo))

	return 0
}
