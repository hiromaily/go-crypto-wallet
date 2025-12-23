package xrp

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GenerateAddress calls GenerateAddress API
func (r *Ripple) GenerateAddress(ctx context.Context) (*ResponseGenerateAddress, error) {
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call addressClient.GenerateAddress(): %w", err)
	}
	logger.Debug("response",
		"XAddress", res.XAddress,
		"ClassicAddress", res.ClassicAddress,
		"Address", res.Address,
		"Secret", res.Secret,
	)

	return res, nil
}

// GenerateXAddress calls GenerateXAddress API
func (r *Ripple) GenerateXAddress(ctx context.Context) (*ResponseGenerateXAddress, error) {
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateXAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call addressClient.GenerateXAddress(): %w", err)
	}
	logger.Debug("response",
		"XAddress", res.XAddress,
		"Secret", res.Secret,
	)

	return res, nil
}

// IsValidAddress calls IsValidAddress API
func (r *Ripple) IsValidAddress(ctx context.Context, addr string) (bool, error) {
	req := &RequestIsValidAddress{
		Address: addr,
	}

	res, err := r.API.addressClient.IsValidAddress(ctx, req)
	if err != nil {
		return false, fmt.Errorf("fail to call addressClient.IsValidAddress(): %w", err)
	}
	logger.Debug("response",
		"IsValid", res.IsValid,
	)

	return res.IsValid, nil
}
