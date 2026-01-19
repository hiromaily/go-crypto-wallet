package xrp

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/protogen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GenerateAddress calls GenerateAddress API
func (r *Ripple) GenerateAddress(ctx context.Context) (*dtoRipple.ResponseGenerateAddress, error) {
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call addressClient.GenerateAddress(): %w", err)
	}
	logger.Debug("response",
		"XAddress", res.GetXAddress(),
		"ClassicAddress", res.GetClassicAddress(),
		"Address", res.GetAddress(),
		"Secret", res.GetSecret(),
	)

	// Convert infrastructure type to DTO
	return ToDTOResponseGenerateAddress(res), nil
}

// GenerateXAddress calls GenerateXAddress API
func (r *Ripple) GenerateXAddress(ctx context.Context) (*dtoRipple.ResponseGenerateXAddress, error) {
	req := &emptypb.Empty{}

	res, err := r.API.addressClient.GenerateXAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fail to call addressClient.GenerateXAddress(): %w", err)
	}
	logger.Debug("response",
		"XAddress", res.GetXAddress(),
		"Secret", res.GetSecret(),
	)

	// Convert infrastructure type to DTO
	return ToDTOResponseGenerateXAddress(res), nil
}

// IsValidAddress calls IsValidAddress API
func (r *Ripple) IsValidAddress(ctx context.Context, addr string) (bool, error) {
	req := protogen.RequestIsValidAddress_builder{
		Address: addr,
	}.Build()

	res, err := r.API.addressClient.IsValidAddress(ctx, req)
	if err != nil {
		return false, fmt.Errorf("fail to call addressClient.IsValidAddress(): %w", err)
	}
	logger.Debug("response",
		"IsValid", res.GetIsValid(),
	)

	return res.GetIsValid(), nil
}
