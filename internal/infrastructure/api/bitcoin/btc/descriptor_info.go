package btc

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
