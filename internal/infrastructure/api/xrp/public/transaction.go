package xrp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// CreateRawTransaction creates raw transaction
// - https://xrpl.org/ja/send-xrp.html
func (r *publicRPC) CreateRawTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *dtoxrp.Instructions,
) (*dtoxrp.TxInput, string, error) {
	// validation
	if senderAccount == "" {
		return nil, "", errors.New("senderAccount is empty")
	}
	if receiverAccount == "" {
		return nil, "", errors.New("receiverAccount is empty")
	}

	// get balance
	accountInfo, err := r.GetAccountInfo(ctx, senderAccount)
	if err != nil {
		errStatus, _ := status.FromError(err)
		return nil, "", fmt.Errorf(
			"fail to call GetAccountInfo() code: %d, message: %s",
			errStatus.Code(), errStatus.Message())
	}
	if amount != 0 && (xrpkg.ToFloat64(accountInfo.XrpBalance)-xrpkg.MinimumReserve) <= amount {
		return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
	}

	// get fee
	txJSON, stringJSON, err := r.PrepareTransaction(ctx, senderAccount, receiverAccount, amount, instructions)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call PrepareTransaction(): %w", err)
	}
	feeDrops := xrpkg.XRPToDrops(xrpkg.ToFloat64(txJSON.Fee))
	calculatedAmount := xrpkg.ToFloat64(accountInfo.XrpBalance) - xrpkg.MinimumReserve - feeDrops
	if amount == 0 {
		// send all, but fee should be calculated first
		if calculatedAmount <= 0 {
			return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
		}
		// re-run
		txJSON, stringJSON, err = r.PrepareTransaction(
			ctx, senderAccount, receiverAccount, calculatedAmount, instructions)
		if err != nil {
			return nil, "", fmt.Errorf("fail to call PrepareTransaction(): %w", err)
		}
	} else if calculatedAmount < amount {
		return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
	}

	return txJSON, stringJSON, nil
}

// WaitValidation waits until the ledger advances past targetLedgerVersion.
// In standalone mode it calls ledger_accept via the admin WebSocket to advance the ledger.
// In non-standalone environments the ledger advances naturally and this polls ledger_current.
func (r *publicRPC) WaitValidation(ctx context.Context, targetLedgerVersion uint64) (uint64, error) {
	// Advance the ledger via the admin connection (standalone mode).
	// Silently ignore errors — in production environments this command is unavailable.
	if w.admin != nil {
		if _, err := r.caller.LedgerAccept(ctx, w.admin); err != nil {
			logger.Warn("ledger_accept failed (non-critical; may not be standalone mode)", "error", err)
		}
	}

	// Poll ledger_current until it reaches or exceeds targetLedgerVersion.
	// In standalone mode each iteration calls ledger_accept to advance the ledger.
	const maxRetries = 30
	for range maxRetries {
		res, err := r.caller.LedgerCurrent(ctx, w.public)
		if err != nil {
			return 0, fmt.Errorf("fail to call ledger_current: %w", err)
		}
		currentLedger := res.Result.LedgerCurrentIndex
		logger.Info("WaitValidation polling",
			"currentLedger", currentLedger,
			"targetLedgerVersion", targetLedgerVersion,
		)
		if currentLedger >= targetLedgerVersion {
			return currentLedger, nil
		}
		// Advance the ledger in standalone mode and retry.
		if w.admin != nil {
			if _, err := r.caller.LedgerAccept(ctx, w.admin); err != nil {
				logger.Warn("failed to advance ledger", "error", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return 0, errors.New("timeout waiting for ledger validation")
}

// GetTransaction retrieves a validated transaction by hash via WebSocket.
func (r *publicRPC) GetTransaction(
	ctx context.Context, txID string, targetLedgerVersion uint64,
) (*xrplgo.TxInfo, error) {
	res, err := r.caller.GetTx(ctx, w.public, txID, targetLedgerVersion)
	if err != nil {
		return nil, fmt.Errorf("fail to call client.GetTransaction(): %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to get transaction info by %s: %s", txID, res.Error)
	}

	logger.Debug("response of getTransaction",
		"hash", res.Result.Hash,
		"ledger_index", res.Result.LedgerIndex,
		"validated", res.Result.Validated,
		"TransactionResult", res.Result.Meta.TransactionResult,
	)

	return &xrplgo.TxInfo{
		ID: res.Result.Hash,
		Outcome: xrplgo.TxOutcome{
			Result:        res.Result.Meta.TransactionResult,
			LedgerVersion: int(res.Result.LedgerIndex),
		},
	}, nil
}
