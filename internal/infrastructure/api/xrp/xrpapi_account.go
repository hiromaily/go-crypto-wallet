package xrp

import (
	"context"
	"errors"
	"fmt"

	xrpclient "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/client"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAccountInfo calls GetAccountInfo API
func (r *XRP) GetAccountInfo(ctx context.Context, address string) (*xrpclient.AccountInfo, error) {
	// validation
	if address == "" {
		return nil, errors.New("address is empty")
	}

	req := protogen.RequestGetAccountInfo_builder{
		Address: address,
	}.Build()

	res, err := r.API.accountClient.GetAccountInfo(ctx, req)
	if err != nil {
		// errStatus, _ := status.FromError(err)
		// errStatus.Message()
		// errStatus.Code()
		return nil, fmt.Errorf("fail to call accountClient.GetAccountInfo(): %w", err)
	}
	logger.Debug("response",
		"Sequence", res.GetSequence(),
		"XrpBalance", res.GetXrpBalance(),
		"OwnerCount", res.GetOwnerCount(),
		"PreviousAffectingTransactionID", res.GetPreviousAffectingTransactionID(),
		"PreviousAffectingTransactionLedgerVersion", res.GetPreviousAffectingTransactionLedgerVersion(),
	)

	// Convert protogen type to pkg client type
	return fromProtoAccountInfo(res), nil
}
