package xrp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GenerateAddress calls GenerateAddress API
func (r *Ripple) GenerateAddress() (*ResponseGenerateAddress, error) {
	ctx := context.Background()
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateAddress(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call addressClient.GenerateAddress()")
	}
	r.logger.Debug("response",
		"XAddress", res.XAddress,
		"ClassicAddress", res.ClassicAddress,
		"Address", res.Address,
		"Secret", res.Secret,
	)

	return res, nil
}

// GenerateXAddress calls GenerateXAddress API
func (r *Ripple) GenerateXAddress() (*ResponseGenerateXAddress, error) {
	ctx := context.Background()
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateXAddress(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call addressClient.GenerateXAddress()")
	}
	r.logger.Debug("response",
		"XAddress", res.XAddress,
		"Secret", res.Secret,
	)

	return res, nil
}

// IsValidAddress calls IsValidAddress API
func (r *Ripple) IsValidAddress(addr string) (bool, error) {
	ctx := context.Background()
	req := &RequestIsValidAddress{
		Address: addr,
	}

	res, err := r.API.addressClient.IsValidAddress(ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "fail to call addressClient.IsValidAddress()")
	}
	r.logger.Debug("response",
		"IsValid", res.IsValid,
	)

	return res.IsValid, nil
}
