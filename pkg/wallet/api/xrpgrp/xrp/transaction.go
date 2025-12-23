package xrp

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/status"
)

// CreateRawTransaction creates raw transaction
// - https://xrpl.org/ja/send-xrp.html
func (r *Ripple) CreateRawTransaction(
	ctx context.Context, senderAccount, receiverAccount string, amount float64, instructions *Instructions,
) (*TxInput, string, error) {
	// validation
	if senderAccount == "" {
		return nil, "", errors.New("senderAccount is empty")
	}
	if receiverAccount == "" {
		return nil, "", errors.New("receiverAccount is empty")
	}

	// get balance
	// xrp.MinimumReserve
	accountInfo, err := r.GetAccountInfo(ctx, senderAccount)
	if err != nil {
		errStatus, _ := status.FromError(err)
		return nil, "", fmt.Errorf(
			"fail to call GetAccountInfo() code: %d, message: %s",
			errStatus.Code(), errStatus.Message())
	}
	if amount != 0 && (ToFloat64(accountInfo.XrpBalance)-MinimumReserve) <= amount {
		return nil, "", fmt.Errorf("balance is short to send %s", accountInfo.XrpBalance)
	}

	// get fee
	txJSON, stringJSON, err := r.PrepareTransaction(ctx, senderAccount, receiverAccount, amount, instructions)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call PrepareTransaction(): %w", err)
	}
	calculatedAmount := ToFloat64(accountInfo.XrpBalance) - MinimumReserve - XRPToDrops(ToFloat64(txJSON.Fee))
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
