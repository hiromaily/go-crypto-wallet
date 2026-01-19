package xrp

import (
	"context"
	"errors"
	"fmt"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAccountInfo calls GetAccountInfo API
func (r *Ripple) GetAccountInfo(ctx context.Context, address string) (*dtoRipple.ResponseGetAccountInfo, error) {
	// validation
	if address == "" {
		return nil, errors.New("address is empty")
	}

	req := &protogen.RequestGetAccountInfo{
		Address: address,
	}

	res, err := r.API.accountClient.GetAccountInfo(ctx, req)
	if err != nil {
		// errStatus, _ := status.FromError(err)
		// errStatus.Message()
		// errStatus.Code()
		return nil, fmt.Errorf("fail to call accountClient.GetAccountInfo(): %w", err)
	}
	logger.Debug("response",
		"Sequence", res.Sequence,
		"XrpBalance", res.XrpBalance,
		"OwnerCount", res.OwnerCount,
		"PreviousAffectingTransactionID", res.PreviousAffectingTransactionID,
		"PreviousAffectingTransactionLedgerVersion", res.PreviousAffectingTransactionLedgerVersion,
	)

	// Convert infrastructure type to DTO
	return ToDTOResponseGetAccountInfo(res), nil
}
