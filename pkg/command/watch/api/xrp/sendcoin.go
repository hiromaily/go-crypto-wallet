package xrp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bookerzzz/grok"
	"google.golang.org/grpc/status"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
)

func runSendCoin(xrpAPI ripple.Rippler, txData *config.RippleTxData, receiverAddr string, amount float64) error {
	// validator
	if receiverAddr == "" {
		return errors.New("address option [-address] is invalid")
	}

	// send coin
	// PrepareTransaction
	instructions := &xrp.Instructions{
		MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
	}
	fmt.Printf("sender: %s, receiver: %s, amount: %v\n", txData.Account, receiverAddr, amount)
	txJSON, _, err := xrpAPI.CreateRawTransaction(context.TODO(), txData.Account, receiverAddr, amount, instructions)
	if err != nil {
		return fmt.Errorf("fail to call xrp.CreateRawTransaction() %w", err)
	}
	grok.Value(txJSON)

	// SingTransaction
	txID, txBlob, err := xrpAPI.SignTransaction(context.TODO(), txJSON, txData.Secret)
	if err != nil {
		return fmt.Errorf("fail to call xrp.SignTransaction() %w", err)
	}

	// SendTransaction
	sentTx, earlistLedgerVersion, err := xrpAPI.SubmitTransaction(context.TODO(), txBlob)
	if err != nil {
		return fmt.Errorf("fail to call xrp.SubmitTransaction() %w", err)
	}
	if strings.Contains(sentTx.ResultCode, "UNFUNDED_PAYMENT") {
		return fmt.Errorf(
			"fail to call SubmitTransaction. resultCode: %s, resultMessage: %s",
			sentTx.ResultCode, sentTx.ResultMessage)
	}

	// validate transaction
	_, err = xrpAPI.WaitValidation(context.TODO(), sentTx.TxJSON.LastLedgerSequence)
	if err != nil {
		return fmt.Errorf("fail to call xrp.WaitValidation() %w", err)
	}

	// get transaction info
	txInfo, err := xrpAPI.GetTransaction(context.TODO(), txID, earlistLedgerVersion)
	if err != nil {
		return fmt.Errorf("fail to call xrp.GetTransaction() %w", err)
	}
	fmt.Printf("transaction Info: %v\n", txInfo)

	// get receiver info
	accountInfo, err := xrpAPI.GetAccountInfo(context.TODO(), receiverAddr)
	if err != nil {
		errStatus, _ := status.FromError(err)
		return fmt.Errorf(
			"fail to call xrp.GetAccountInfo() code: %d, message: %s",
			errStatus.Code(), errStatus.Message())
	}
	fmt.Printf("receiver account Info: %v\n", accountInfo)

	return nil
}
