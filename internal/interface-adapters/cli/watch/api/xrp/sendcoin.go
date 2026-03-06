package xrp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bookerzzz/grok"
	"google.golang.org/grpc/status"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// sendCoinClient defines the interface for XRP operations needed by runSendCoin.
type sendCoinClient interface {
	apixrp.TransactionPreparer
	apixrp.TransactionSigner
	apixrp.TransactionSubmitter
	apixrp.LedgerPoller
	apixrp.AccountInfoProvider
}

// waitForValidation polls ledger_current until the ledger advances past targetVersion.
func waitForValidation(ctx context.Context, client apixrp.LedgerPoller, targetVersion uint64) error {
	const maxRetries = 30
	for range maxRetries {
		res, err := client.LedgerCurrent(ctx)
		if err != nil {
			return fmt.Errorf("fail to call ledger_current: %w", err)
		}
		if res.Result.LedgerCurrentIndex >= targetVersion {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("timeout waiting for ledger validation")
}

func runSendCoin(xrpAPI sendCoinClient, txData *config.RippleTxData, receiverAddr string, amount float64) error {
	// validator
	if receiverAddr == "" {
		return errors.New("address option [-address] is invalid")
	}

	// send coin
	instructions := &dtoxrp.Instructions{
		MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
	}
	fmt.Printf("sender: %s, receiver: %s, amount: %v\n", txData.Account, receiverAddr, amount)
	txJSON, _, err := xrpAPI.CreateRawTransaction(context.TODO(), txData.Account, receiverAddr, amount, instructions)
	if err != nil {
		return fmt.Errorf("fail to call xrp.CreateRawTransaction() %w", err)
	}
	grok.Value(txJSON)

	// SignTransaction
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

	// Wait for ledger to advance past the transaction's LastLedgerSequence
	if err = waitForValidation(context.TODO(), xrpAPI, sentTx.TxJSON.LastLedgerSequence); err != nil {
		logger.Warn("ledger validation timed out", "error", err)
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
