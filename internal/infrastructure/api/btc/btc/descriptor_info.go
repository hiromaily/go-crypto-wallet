package btc

// Descriptor Info - BTC ONLY (BCH does NOT support descriptors)
// See descriptor_service.go for full warning.

import (
	"fmt"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetDescriptorInfo calls Bitcoin Core's getdescriptorinfo RPC to analyze a descriptor
// and calculate its BIP380 checksum.
//
// Bitcoin Core RPC: getdescriptorinfo "descriptor"
//
// This method is used to:
//   - Calculate the correct BIP380 checksum for a descriptor
//   - Validate descriptor syntax
//   - Check if the descriptor is solvable
//
// Returns the descriptor with checksum and analysis results.
func (b *Bitcoin) GetDescriptorInfo(descriptor string) (*btcrpc.DescriptorInfo, error) {
	result, err := btcrpc.GetDescriptorInfo(b.Client, descriptor)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetDescriptorInfo(): %w", err)
	}

	return result, nil
}

// ListDescriptors calls Bitcoin Core's listdescriptors RPC to list all imported descriptors.
//
// Bitcoin Core RPC: listdescriptors ( private )
//
// This method returns all descriptors that have been imported into the wallet,
// including their status (active, internal) and range information for ranged descriptors.
//
// Parameters:
//   - privateDescriptors: if true, return private descriptors (requires wallet with private keys)
//     For watch-only wallets, this must be false.
//
// Returns the list of descriptors with their metadata.
func (b *Bitcoin) ListDescriptors(privateDescriptors bool) (*btcrpc.ListDescriptorsResult, error) {
	result, err := btcrpc.ListDescriptors(b.Client, privateDescriptors)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ListDescriptors(): %w", err)
	}

	return result, nil
}
