package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAccountInfo calls GetAccountInfo API
func (r *XRP) GetAccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseGetAccountInfo, error) {
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

	// Convert infrastructure type to DTO
	return ToDTOResponseGetAccountInfo(res), nil
}
