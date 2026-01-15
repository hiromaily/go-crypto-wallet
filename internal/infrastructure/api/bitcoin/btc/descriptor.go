package btc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/txscript"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
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
		originPath := match[2]    // Origin path (e.g., "/44'/1'/1'")
		xpub := match[3]          // Extended public key
		remainingPath := match[4] // Path after xpub (derivation path FROM xpub)

		// Store both origin and remaining paths separately:
		// - OriginPath: Used for PSBT BIP32 derivation metadata (full path from master)
		// - DerivationPath: Used for key derivation from xpub (xpub is already at origin level)
		// Example: [4a91842f/84'/1'/1']tpubDDed.../0/*
		//   OriginPath = "/84'/1'/1'" (where xpub was derived from master)
		//   DerivationPath = "/0/*" (what to derive from xpub)
		//   For PSBT: combine both = "/84'/1'/1'/0/0" (after replacing wildcard)

		key := domainWallet.DescriptorKey{
			Fingerprint:    fingerprint,
			OriginPath:     originPath,
			DerivationPath: remainingPath,
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

// DeriveRedeemScriptFromDescriptor derives a P2SH redeemScript from a descriptor at a specific index.
// This is used when Bitcoin Core's listunspent doesn't return the redeemScript for descriptor-based addresses.
//
// Supported descriptor types:
//
//   - sh(multi(...)): Legacy P2SH multisig
//     Example: sh(multi(2,[fp1/44'/1'/1']tpub.../0/*,[fp2/44'/1'/1']tpub.../0/*))
//     RedeemScript: OP_2 <pk1> <pk2> OP_2 OP_CHECKMULTISIG
//
//   - sh(wpkh(...)): P2SH-wrapped P2WPKH (BIP49 Nested SegWit)
//     Example: sh(wpkh([fp/49'/1'/1']tpub.../0/*))
//     RedeemScript: OP_0 <20-byte-pubkey-hash>
func (b *Bitcoin) DeriveRedeemScriptFromDescriptor(
	descriptor string, address string, addressIndex uint32,
) ([]byte, error) {
	logger.Debug("Deriving redeemScript from descriptor",
		"descriptor_len", len(descriptor),
		"address", address,
		"index", addressIndex)

	// Parse using the existing descriptor parser
	parser := NewDescriptorParser()
	parsed, err := parser.Parse(descriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse descriptor: %w", err)
	}

	logger.Info("Parsed descriptor",
		"type", parsed.Type,
		"total_keys", len(parsed.Keys),
		"script_preview", parsed.Script[:min(100, len(parsed.Script))])

	// Route to appropriate handler based on descriptor type
	var redeemScript []byte
	switch parsed.Type {
	case domainWallet.DescriptorTypeSH:
		// Legacy P2SH multisig
		redeemScript, err = b.deriveMultisigRedeemScript(parsed, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive multisig redeemScript: %w", err)
		}

	case domainWallet.DescriptorTypeSHWPKH:
		// P2SH-wrapped P2WPKH (BIP49 Nested SegWit)
		redeemScript, err = b.deriveP2WPKHRedeemScript(parsed, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive P2WPKH redeemScript: %w", err)
		}

	case domainWallet.DescriptorTypePKH,
		domainWallet.DescriptorTypeWPKH,
		domainWallet.DescriptorTypeTR,
		domainWallet.DescriptorTypeWSH,
		domainWallet.DescriptorTypeSHWSH,
		domainWallet.DescriptorTypeUnknown:
		// These descriptor types don't require redeemScript derivation:
		// - PKH: P2PKH addresses don't use redeemScript
		// - WPKH: Native SegWit (bech32) addresses don't use redeemScript
		// - TR: Taproot addresses don't use redeemScript
		// - WSH/SHWSH: Not yet implemented
		return nil, fmt.Errorf("unsupported descriptor type for redeemScript derivation: %s", parsed.Type)
	}

	// DIAGNOSTIC: Hash the redeemScript to get the P2SH address
	redeemScriptHash := btcutil.Hash160(redeemScript)
	p2shAddr, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, b.chainConf)
	var derivedAddrStr string
	if err == nil {
		derivedAddrStr = p2shAddr.EncodeAddress()
	}

	logger.Info("Derived redeemScript from descriptor",
		"target_address", address,
		"derived_address", derivedAddrStr,
		"match", derivedAddrStr == address,
		"index", addressIndex,
		"descriptor_type", parsed.Type,
		"script_len", len(redeemScript),
		"script_hex", hex.EncodeToString(redeemScript))

	return redeemScript, nil
}

// deriveMultisigRedeemScript derives a redeemScript for sh(multi(...)) or sh(sortedmulti(...)) descriptors.
// This handles legacy P2SH multisig addresses.
//
// For a descriptor like:
//
//	sh(sortedmulti(3,[fp1/path]xpub1/0/*,[fp2/path]xpub2/0/*,[fp3/path]xpub3/0/*))
//
// At index 0, it derives the three public keys and builds:
//
//	OP_3 <pk1> <pk2> <pk3> OP_3 OP_CHECKMULTISIG
//
// For sortedmulti, public keys are sorted lexicographically (BIP67).
//

func (b *Bitcoin) deriveMultisigRedeemScript(
	parsed *domainWallet.Descriptor, addressIndex uint32,
) ([]byte, error) {
	// Extract multisig parameters from the descriptor
	requiredSigs, totalSigs, err := b.extractMultisigParams(parsed.Script)
	if err != nil {
		return nil, fmt.Errorf("failed to extract multisig params: %w", err)
	}

	// Check if it's sorted or not
	isSorted := strings.Contains(parsed.Script, "sortedmulti(")

	logger.Info("Parsed multisig descriptor",
		"required_sigs", requiredSigs,
		"total_keys", len(parsed.Keys),
		"is_sorted", isSorted)

	// Derive public keys at the specified index
	pubKeys := make([][]byte, 0, len(parsed.Keys))
	for i, keyInfo := range parsed.Keys {
		pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive public key %d: %w", i, err)
		}
		pubKeys = append(pubKeys, pubKey)
		logger.Debug("Derived public key from descriptor",
			"key_index", i,
			"xpub_first_10", keyInfo.ExtendedPubKey[:10]+"...",
			"derivation_path", keyInfo.DerivationPath,
			"address_index", addressIndex,
			"pubkey", hex.EncodeToString(pubKey))
	}

	// Build multisig redeemScript: OP_M <pk1> <pk2> ... <pkN> OP_N OP_CHECKMULTISIG
	return b.buildMultisigRedeemScript(requiredSigs, totalSigs, pubKeys, isSorted)
}

// deriveP2WPKHRedeemScript derives a redeemScript for sh(wpkh(...)) descriptors.
// This handles P2SH-wrapped P2WPKH addresses (BIP49 Nested SegWit).
//
// For a descriptor like:
//
//	sh(wpkh([fingerprint/49'/1'/1']tpubDDYi.../0/*))
//
// At index 0, it:
// 1. Derives the public key at path /0/0 from the xpub
// 2. Hashes the public key: HASH160(pubkey) = RIPEMD160(SHA256(pubkey))
// 3. Builds P2WPKH script: OP_0 <20-byte-hash>
//
// This P2WPKH script becomes the redeemScript for the P2SH wrapper.
// When spending:
//   - scriptSig: <redeemScript>
//   - witness: [<signature>, <pubkey>]
//

func (b *Bitcoin) deriveP2WPKHRedeemScript(
	parsed *domainWallet.Descriptor, addressIndex uint32,
) ([]byte, error) {
	if len(parsed.Keys) != 1 {
		return nil, fmt.Errorf("sh(wpkh) descriptor must have exactly 1 key, got %d", len(parsed.Keys))
	}

	keyInfo := parsed.Keys[0]

	// Derive the public key at the specified index
	pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %w", err)
	}

	logger.Debug("Derived public key for P2WPKH",
		"xpub_first_10", keyInfo.ExtendedPubKey[:10]+"...",
		"derivation_path", keyInfo.DerivationPath,
		"address_index", addressIndex,
		"pubkey", hex.EncodeToString(pubKey))

	// Hash the public key: RIPEMD160(SHA256(pubkey))
	pubKeyHash := btcutil.Hash160(pubKey)

	logger.Debug("Computed pubkey hash",
		"pubkey_hash", hex.EncodeToString(pubKeyHash))

	// Build P2WPKH script: OP_0 <20-byte-hash>
	// This is the redeemScript for the P2SH wrapper
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_0)
	builder.AddData(pubKeyHash)

	redeemScript, err := builder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build P2WPKH redeemScript: %w", err)
	}

	logger.Info("Built P2WPKH redeemScript",
		"script_len", len(redeemScript),
		"script_hex", hex.EncodeToString(redeemScript))

	return redeemScript, nil
}

// extractMultisigParams extracts M and N from a multisig descriptor
// Format: sh(multi(M,key1,key2,...)) or sh(sortedmulti(M,key1,key2,...))
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) extractMultisigParams(descriptorScript string) (requiredSigs, totalSigs int, err error) {
	// Match multi(M, or sortedmulti(M,
	multiRegex := regexp.MustCompile(`(?:sorted)?multi\((\d+),`)
	matches := multiRegex.FindStringSubmatch(descriptorScript)
	if len(matches) < 2 {
		return 0, 0, errors.New("failed to extract multisig M parameter")
	}

	requiredSigs, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid M parameter: %s", matches[1])
	}

	// Count the number of keys
	// Keys have format: [fingerprint/path]xpub or just xpub
	keyPattern := regexp.MustCompile(`\[[0-9a-fA-F]{8}[^\]]*\][xyztYZ]pub`)
	keyMatches := keyPattern.FindAllString(descriptorScript, -1)
	totalSigs = len(keyMatches)

	if totalSigs == 0 {
		return 0, 0, errors.New("no keys found in multisig descriptor")
	}

	return requiredSigs, totalSigs, nil
}

// derivePublicKeyFromDescriptorKey derives a child public key from a descriptor key
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) derivePublicKeyFromDescriptorKey(
	keyInfo domainWallet.DescriptorKey, addressIndex uint32,
) ([]byte, error) {
	// Parse the derivation path
	// Format: "/0/*" or "/1/*" or "/*"
	// Replace wildcard with actual address index before parsing
	derivationPath := strings.ReplaceAll(keyInfo.DerivationPath, "/*", fmt.Sprintf("/%d", addressIndex))
	derivationPath = strings.ReplaceAll(derivationPath, "*", strconv.FormatUint(uint64(addressIndex), 10))
	pathParts := strings.Split(strings.TrimPrefix(derivationPath, "/"), "/")

	// Parse extended public key
	extKey, err := hdkeychain.NewKeyFromString(keyInfo.ExtendedPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xpub: %w", err)
	}

	// Derive child keys according to the path
	currentKey := extKey
	for _, part := range pathParts {
		if part == "" {
			continue
		}

		// Parse the index
		index, parseErr := strconv.ParseUint(part, 10, 32)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid path component %s: %w", part, parseErr)
		}

		// Derive the key
		currentKey, err = currentKey.Derive(uint32(index))
		if err != nil {
			return nil, fmt.Errorf("failed to derive index %d: %w", index, err)
		}
	}

	// Get the public key
	pubKey, err := currentKey.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	return pubKey.SerializeCompressed(), nil
}

// buildMultisigRedeemScript builds a multisig redeemScript from public keys
// Format: OP_M <pk1> <pk2> ... <pkN> OP_N OP_CHECKMULTISIG
//
// For sortedmulti descriptors (BIP67), public keys are sorted lexicographically
// before being added to the script. This ensures deterministic multisig addresses
// regardless of the order keys are specified in the descriptor.
func (b *Bitcoin) buildMultisigRedeemScript(
	requiredSigs, totalSigs int, pubKeys [][]byte, sorted bool,
) ([]byte, error) {
	if requiredSigs <= 0 || requiredSigs > totalSigs {
		return nil, fmt.Errorf("invalid multisig parameters: %d-of-%d", requiredSigs, totalSigs)
	}

	if len(pubKeys) != totalSigs {
		return nil, fmt.Errorf("pubkey count mismatch: expected %d, got %d", totalSigs, len(pubKeys))
	}

	// Sort public keys lexicographically if this is a sortedmulti descriptor (BIP67)
	if sorted {
		b.sortPublicKeys(pubKeys)
		logger.Debug("Sorted public keys for sortedmulti descriptor",
			"count", len(pubKeys))
	}

	// Build the script
	builder := txscript.NewScriptBuilder()

	// Add OP_M
	builder.AddInt64(int64(requiredSigs))

	// Add public keys
	for _, pubKey := range pubKeys {
		builder.AddData(pubKey)
	}

	// Add OP_N
	builder.AddInt64(int64(totalSigs))

	// Add OP_CHECKMULTISIG
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	return builder.Script()
}

// sortPublicKeys sorts public keys lexicographically (BIP67).
// This is used for sortedmulti descriptors to ensure deterministic ordering.
// The keys are sorted in-place.
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) sortPublicKeys(pubKeys [][]byte) {
	sort.Slice(pubKeys, func(i, j int) bool {
		return bytes.Compare(pubKeys[i], pubKeys[j]) < 0
	})
}
