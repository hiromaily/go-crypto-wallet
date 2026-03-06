package xrppublic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

const dropsPerXRP = 1_000_000

// GetAccountInfo retrieves account information via WebSocket.
func (p *PublicXRP) GetAccountInfo(ctx context.Context, address string) (*dtoxrp.AccountInfo, error) {
	if address == "" {
		return nil, errors.New("address is empty")
	}

	res, err := p.PublicRPC.AccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call AccountInfo(): %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call AccountInfo(): %s", res.Error)
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

// GetBalance returns amount of address.
func (p *PublicXRP) GetBalance(ctx context.Context, addr string) (float64, error) {
	accountInfo, err := p.GetAccountInfo(ctx, addr)
	if err != nil {
		return 0, err
	}
	return xrpkg.ToFloat64(accountInfo.XrpBalance), nil
}

// GetTotalBalance returns total amount in address list.
func (p *PublicXRP) GetTotalBalance(ctx context.Context, addrs []string) float64 {
	var total float64
	for _, addr := range addrs {
		amt, err := p.GetBalance(ctx, addr)
		if err == nil {
			total += amt
		}
	}
	return total
}

// fromWSAccountInfo converts a WebSocket ResponseAccountInfo to the application DTO.
func fromWSAccountInfo(res *xrprpcpublic.ResponseAccountInfo) (*dtoxrp.AccountInfo, error) {
	drops, err := strconv.ParseInt(res.Result.AccountData.Balance, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account balance %q: %w", res.Result.AccountData.Balance, err)
	}
	xrpBalance := strconv.FormatFloat(float64(drops)/dropsPerXRP, 'f', -1, 64)
	return &dtoxrp.AccountInfo{
		Sequence:                       uint64(res.Result.AccountData.Sequence),
		XrpBalance:                     xrpBalance,
		OwnerCount:                     uint64(res.Result.AccountData.OwnerCount),
		PreviousAffectingTransactionID: res.Result.AccountData.PreviousTxnID,
		PreviousAffectingTransactionLedgerVersion: uint64(res.Result.AccountData.PreviousTxnLgrSeq),
	}, nil
}
