package btc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// DescriptorParser parses Bitcoin output descriptors according to BIP380.
//
// The parser extracts:
//   - Descriptor type (pkh, sh(wpkh), wpkh, tr, wsh)
//   - Keys with fingerprints and derivation paths
//   - Checksum (if present)
//
// Example descriptors:
//   - P2PKH: pkh([a1b2c3d4/44'/0'/0']xpub.../0/*)
//   - Bech32: wpkh([a1b2c3d4/84'/0'/0']xpub.../0/*)
//   - Taproot: tr([a1b2c3d4/86'/0'/0']xpub.../0/*)
type DescriptorParser struct {
	// keyRegex matches key patterns in descriptors with full metadata
	// Format: [fingerprint/derivation/path]xpub.../path/*
	keyRegex *regexp.Regexp

	// simpleKeyRegex matches keys without fingerprint/path metadata
	// Format: xpub... or xpub.../path
	simpleKeyRegex *regexp.Regexp
}

// NewDescriptorParser creates a new descriptor parser.
func NewDescriptorParser() *DescriptorParser {
	// Regex to match descriptor keys with full metadata
	// Captures:
	// 1. Optional [fingerprint/path]
	// 2. Extended public key (xpub/tpub/etc.)
	// 3. Optional /remaining/path
	//nolint:revive // Regex pattern necessarily long for descriptor parsing
	keyRegex := regexp.MustCompile(`\[([0-9a-fA-F]{8})((?:/\d+['h]?)*)\]([xyztYZ]pub[1-9A-HJ-NP-Za-km-z]+)((?:/\d+['h]?|/\*)*)`)

	// Regex to match simple keys without metadata
	simpleKeyRegex := regexp.MustCompile(`([xyztYZ]pub[1-9A-HJ-NP-Za-km-z]+)`)

	return &DescriptorParser{
		keyRegex:       keyRegex,
		simpleKeyRegex: simpleKeyRegex,
	}
}

// Parse parses a descriptor string into a Descriptor object.
//
// The parser:
//  1. Determines the descriptor type
//  2. Extracts keys with their metadata
//  3. Extracts checksum if present
//  4. Validates the extracted data
//
// Returns an error if the descriptor is malformed or unsupported.
func (p *DescriptorParser) Parse(descriptorStr string) (*domainWallet.Descriptor, error) {
	if strings.TrimSpace(descriptorStr) == "" {
		return nil, errors.New("descriptor string is empty")
	}

	// Strip whitespace
	descriptorStr = strings.TrimSpace(descriptorStr)

	// Extract checksum if present (format: descriptor#checksum)
	var checksum string
	if idx := strings.LastIndex(descriptorStr, "#"); idx != -1 {
		checksum = descriptorStr[idx+1:]
		descriptorStr = descriptorStr[:idx]
	}

	// Determine descriptor type
	descType, err := determineType(descriptorStr)
	if err != nil {
		return nil, fmt.Errorf("failed to determine descriptor type: %w", err)
	}

	// Extract keys
	keys, err := p.extractKeys(descriptorStr)
	if err != nil {
		return nil, fmt.Errorf("failed to extract keys: %w", err)
	}

	if len(keys) == 0 {
		return nil, errors.New("no keys found in descriptor")
	}

	descriptor := &domainWallet.Descriptor{
		Type:     descType,
		Script:   descriptorStr,
		Keys:     keys,
		Checksum: checksum,
	}

	// Validate the parsed descriptor
	if err := domainWallet.ValidateDescriptor(descriptor); err != nil {
		return nil, fmt.Errorf("invalid descriptor: %w", err)
	}

	return descriptor, nil
}

// determineType determines the descriptor type from the descriptor string.
func determineType(descriptorStr string) (domainWallet.DescriptorType, error) {
	// Check for different descriptor types based on prefix
	// Note: Order matters - check more specific patterns first (sh(wsh) before sh(wpkh))
	switch {
	case strings.HasPrefix(descriptorStr, "pkh("):
		return domainWallet.DescriptorTypePKH, nil
	case strings.HasPrefix(descriptorStr, "sh(wsh("):
		return domainWallet.DescriptorTypeSHWSH, nil
	case strings.HasPrefix(descriptorStr, "sh(wpkh("):
		return domainWallet.DescriptorTypeSHWPKH, nil
	case strings.HasPrefix(descriptorStr, "sh(multi("), strings.HasPrefix(descriptorStr, "sh(sortedmulti("):
		return domainWallet.DescriptorTypeSH, nil
	case strings.HasPrefix(descriptorStr, "wpkh("):
		return domainWallet.DescriptorTypeWPKH, nil
	case strings.HasPrefix(descriptorStr, "tr("):
		return domainWallet.DescriptorTypeTR, nil
	case strings.HasPrefix(descriptorStr, "wsh("):
		return domainWallet.DescriptorTypeWSH, nil
	default:
		return domainWallet.DescriptorTypeUnknown, fmt.Errorf("unknown descriptor type: %s", descriptorStr)
	}
}

// extractKeys extracts keys from the descriptor string.
//
// Keys in descriptors have the format:
//
//	[fingerprint/derivation/path]xpub.../remaining/path
//
// Example:
//
//	[a1b2c3d4/44'/0'/0']xpub6ERApfZw.../0/*
func (p *DescriptorParser) extractKeys(descriptorStr string) ([]domainWallet.DescriptorKey, error) {
	matches := p.keyRegex.FindAllStringSubmatch(descriptorStr, -1)

	if len(matches) == 0 {
		// Try to find keys without fingerprint/path metadata
		// Format: just xpub... or xpub.../path
		simpleMatches := p.simpleKeyRegex.FindAllString(descriptorStr, -1)

		if len(simpleMatches) == 0 {
			return nil, errors.New("no keys found")
		}

		// Create keys without fingerprint/path metadata
		var keys []domainWallet.DescriptorKey
		for _, xpub := range simpleMatches {
			keys = append(keys, domainWallet.DescriptorKey{
				Fingerprint:    "",
				DerivationPath: "",
				ExtendedPubKey: xpub,
			})
		}
		return keys, nil
	}

	keys := make([]domainWallet.DescriptorKey, 0, len(matches))
	for _, match := range matches {
		if len(match) < 5 {
			continue
		}

		fingerprint := match[1]
		// match[2] is accountPath (origin info) - not used, already baked into xpub
		xpub := match[3]          // Extended public key
		remainingPath := match[4] // Path after xpub (derivation path FROM xpub)

		// Use only remainingPath (derivation path FROM the xpub).
		// match[2] is origin metadata (where the xpub was derived from master key),
		// which is already baked into the xpub and should not be derived again.
		// Example: [4a91842f/84'/1'/1]tpubDDed.../0/*
		//   match[2] = "/84'/1'/1" (origin - already in xpub)
		//   remainingPath = "/0/*" (what to derive from xpub)
		// Bitcoin Core's deriveaddresses interprets this the same way.
		fullPath := remainingPath

		key := domainWallet.DescriptorKey{
			Fingerprint:    fingerprint,
			DerivationPath: fullPath,
			ExtendedPubKey: xpub,
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// FormatDescriptor formats a Descriptor object into a descriptor string.
//
// This is the inverse operation of Parse().
func FormatDescriptor(desc *domainWallet.Descriptor) (string, error) {
	if err := domainWallet.ValidateDescriptor(desc); err != nil {
		return "", fmt.Errorf("invalid descriptor: %w", err)
	}

	// For now, return the original script
	// In a full implementation, we would reconstruct the descriptor from components
	result := desc.Script

	if desc.Checksum != "" {
		result += "#" + desc.Checksum
	}

	return result, nil
}
