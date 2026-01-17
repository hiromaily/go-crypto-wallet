package btc

// Descriptor Info - BTC ONLY (BCH does NOT support descriptors)
// See descriptor_service.go for full warning.

import (
	"encoding/json"
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
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
func (b *Bitcoin) GetDescriptorInfo(descriptor string) (*dtobtc.DescriptorInfo, error) {
	params := make([]json.RawMessage, 0, 1)

	// Add descriptor parameter
	bDescriptor, err := json.Marshal(descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal descriptor: %w", err)
	}
	params = append(params, bDescriptor)

	// Call RPC
	rawResult, err := b.Client.RawRequest("getdescriptorinfo", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call getdescriptorinfo RPC: %w", err)
	}

	// Parse result
	var result dtobtc.DescriptorInfo
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal getdescriptorinfo result: %w", err)
	}

	return &result, nil
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
func (b *Bitcoin) ListDescriptors(privateDescriptors bool) (*dtobtc.ListDescriptorsResult, error) {
	params := make([]json.RawMessage, 0, 1)

	// Add private parameter
	bPrivate, err := json.Marshal(privateDescriptors)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private parameter: %w", err)
	}
	params = append(params, bPrivate)

	// Call RPC
	rawResult, err := b.Client.RawRequest("listdescriptors", params)
	if err != nil {
		return nil, fmt.Errorf("failed to call listdescriptors RPC: %w", err)
	}

	// Parse result
	var result dtobtc.ListDescriptorsResult
	if err := json.Unmarshal(rawResult, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listdescriptors result: %w", err)
	}

	return &result, nil
}
