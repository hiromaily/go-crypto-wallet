package xrp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// GetAccountInfo calls GetAccountInfo API
func (r *Ripple) GetAccountInfo(address string) (*pb.ResponseGetAccountInfo, error) {
	// validation
	if address == "" {
		return nil, errors.New("address is empty")
	}

	ctx := context.Background()
	req := &pb.RequestGetAccountInfo{
		Address: address,
	}

	res, err := r.API.accountClient.GetAccountInfo(ctx, req)
	if err != nil {
		// errStatus, _ := status.FromError(err)
		// errStatus.Message()
		// errStatus.Code()
		return nil, errors.Wrap(err, "fail to call accountClient.GetAccountInfo()")
	}
	r.logger.Debug("response",
		zap.Uint64("Sequence", res.Sequence),
		zap.String("XrpBalance", res.XrpBalance),
		zap.Uint64("OwnerCount", res.OwnerCount),
		zap.String("PreviousAffectingTransactionID", res.PreviousAffectingTransactionID),
		zap.Uint64("PreviousAffectingTransactionLedgerVersion", res.PreviousAffectingTransactionLedgerVersion),
	)

	return res, nil
}
