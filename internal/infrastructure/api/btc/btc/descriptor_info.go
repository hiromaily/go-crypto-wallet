package btc

import (
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetDescriptorInfo analyses a descriptor
func (b *Bitcoin) GetDescriptorInfo(descriptor string) (*btcrpc.DescriptorInfo, error) {
	return b.pkgrpc.GetDescriptorInfo(descriptor)
}

// ListDescriptors lists descriptors imported into a descriptor-enabled wallet
func (b *Bitcoin) ListDescriptors(privateDescriptors bool) (*btcrpc.ListDescriptorsResult, error) {
	return b.pkgrpc.ListDescriptors(privateDescriptors)
}
