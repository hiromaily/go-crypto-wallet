package xrp

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	xrpclient "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/client"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

const dropsPerXRP = 1_000_000

// GetAccountInfo retrieves account information via WebSocket.
func (r *XRP) GetAccountInfo(ctx context.Context, address string) (*xrpclient.AccountInfo, error) {
	if address == "" {
		return nil, errors.New("address is empty")
	}

	res, err := xrprpc.AccountInfo(ctx, r.wsPublic, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call accountClient.GetAccountInfo(): %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call accountClient.GetAccountInfo(): %s", res.Error)
	}

	info, err := fromWSAccountInfo(res)
	if err != nil {
		return nil, fmt.Errorf("failed to convert account info: %w", err)
	}
	logger.Debug("response",
		"Sequence", info.Sequence,
		"XrpBalance", info.XrpBalance,
		"OwnerCount", info.OwnerCount,
		"PreviousAffectingTransactionID", info.PreviousAffectingTransactionID,
		"PreviousAffectingTransactionLedgerVersion", info.PreviousAffectingTransactionLedgerVersion,
	)

	return info, nil
}

// fromWSAccountInfo converts a WebSocket ResponseAccountInfo to the pkg client AccountInfo type.
// The Balance field in the WebSocket response is in drops; this converts it to XRP.
func fromWSAccountInfo(res *xrprpc.ResponseAccountInfo) (*xrpclient.AccountInfo, error) {
	drops, err := strconv.ParseInt(res.Result.AccountData.Balance, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account balance %q: %w", res.Result.AccountData.Balance, err)
	}
	xrpBalance := strconv.FormatFloat(float64(drops)/dropsPerXRP, 'f', -1, 64)
	return &xrpclient.AccountInfo{
		Sequence:                       uint64(res.Result.AccountData.Sequence),
		XrpBalance:                     xrpBalance,
		OwnerCount:                     uint64(res.Result.AccountData.OwnerCount),
		PreviousAffectingTransactionID: res.Result.AccountData.PreviousTxnID,
		PreviousAffectingTransactionLedgerVersion: uint64(res.Result.AccountData.PreviousTxnLgrSeq),
	}, nil
}
