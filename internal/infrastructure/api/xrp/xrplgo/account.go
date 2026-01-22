package xrplgo

import (
	"context"
	"fmt"
	"strconv"

	xrpl "github.com/xrpscan/xrpl-go"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
)

// accountInfoResponse represents the response from account_info command.
type accountInfoResponse struct {
	AccountData struct {
		Account           string `json:"Account"`
		Balance           string `json:"Balance"`
		Sequence          uint64 `json:"Sequence"`
		OwnerCount        uint64 `json:"OwnerCount"`
		PreviousTxnID     string `json:"PreviousTxnID"`
		PreviousTxnLgrSeq uint64 `json:"PreviousTxnLgrSeq"`
	} `json:"account_data"`
	LedgerCurrentIndex uint64 `json:"ledger_current_index"`
	Validated          bool   `json:"validated"`
}

// GetAccountInfo retrieves account information from the XRPL.
// Implements AccountInfoProvider interface.
func (c *Client) GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error) {
	req := xrpl.BaseRequest{
		"command":      "account_info",
		"account":      address,
		"ledger_index": "validated",
	}

	resp, err := c.request(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	var result accountInfoResponse
	if err := extractResult(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse account info response: %w", err)
	}

	return &dtoxrp.ResponseGetAccountInfo{
		Sequence:                       result.AccountData.Sequence,
		XrpBalance:                     dropsToXRP(result.AccountData.Balance),
		OwnerCount:                     result.AccountData.OwnerCount,
		PreviousAffectingTransactionID: result.AccountData.PreviousTxnID,
		PreviousAffectingTransactionLedgerVersion: result.AccountData.PreviousTxnLgrSeq,
	}, nil
}

// GetBalance retrieves the XRP balance for an address.
// Implements BalanceChecker interface.
func (c *Client) GetBalance(ctx context.Context, addr string) (float64, error) {
	info, err := c.GetAccountInfo(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	balance, err := strconv.ParseFloat(info.XrpBalance, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse balance: %w", err)
	}

	return balance, nil
}

// GetTotalBalance retrieves the total XRP balance across multiple addresses.
// Implements BalanceChecker interface.
func (c *Client) GetTotalBalance(ctx context.Context, addrs []string) float64 {
	var total float64
	for _, addr := range addrs {
		balance, err := c.GetBalance(ctx, addr)
		if err != nil {
			continue
		}
		total += balance
	}
	return total
}

// dropsToXRP converts drops (1/1,000,000 XRP) to XRP string.
func dropsToXRP(drops string) string {
	dropsInt, err := strconv.ParseInt(drops, 10, 64)
	if err != nil {
		return "0"
	}
	xrp := float64(dropsInt) / 1_000_000
	return strconv.FormatFloat(xrp, 'f', -1, 64)
}
