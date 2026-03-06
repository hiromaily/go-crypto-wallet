package btc

import (
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetDescriptorInfo analyses a descriptor
func (b *Bitcoin) GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error) {
	result, err := b.pkgrpc.GetDescriptorInfo(descriptor)
	if err != nil {
		return nil, fmt.Errorf("fail to call pkgrpc.GetDescriptorInfo(): %w", err)
	}
	return &dtobtc.DescriptorInfo{
		Descriptor:     result.Descriptor,
		Checksum:       result.Checksum,
		IsRange:        result.IsRange,
		IsSolvable:     result.IsSolvable,
		HasPrivateKeys: result.HasPrivateKeys,
	}, nil
}

// ListDescriptors lists descriptors imported into a descriptor-enabled wallet
func (b *Bitcoin) ListDescriptors(privateDescriptors bool) ([]dtobtc.DescriptorInfo, error) {
	result, err := b.pkgrpc.ListDescriptors(privateDescriptors)
	if err != nil {
		return nil, fmt.Errorf("fail to call pkgrpc.ListDescriptors(): %w", err)
	}
	return toDescriptorInfoSlice(result), nil
}

func toDescriptorInfoSlice(r *btcrpc.ListDescriptorsResult) []dtobtc.DescriptorInfo {
	if r == nil {
		return nil
	}
	infos := make([]dtobtc.DescriptorInfo, len(r.Descriptors))
	for i, d := range r.Descriptors {
		infos[i] = dtobtc.DescriptorInfo{
			Descriptor: d.Desc,
			Internal:   d.Internal,
		}
	}
	return infos
}
