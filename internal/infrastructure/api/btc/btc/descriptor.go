package btc

// Descriptor Parser - BTC ONLY
//
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ WARNING: DO NOT USE THIS FILE FOR BCH (Bitcoin Cash) IMPLEMENTATIONS        │
// │                                                                             │
// │ BCH does NOT support output descriptors (BIP380).                           │
// │ BCH uses legacy wallet with address export/import workflow instead.         │
// │                                                                             │
// │ For BCH address management:                                                 │
// │   - Export addresses from keygen wallet                                     │
// │   - Import addresses to watch wallet                                        │
// │   - Do NOT use importdescriptors RPC                                        │
// │                                                                             │
// │ For BCH implementations, see:                                               │
// │   - internal/infrastructure/api/btc/bch/                                    │
// │   - docs/chains/bch/README.md                                               │
// └─────────────────────────────────────────────────────────────────────────────┘

import (
	"bytes"
	"crypto/sha256"
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

	btcdescriptor "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/descriptor"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// NewDescriptorParser creates a new descriptor parser (delegates to pkg).
func NewDescriptorParser() *btcdescriptor.DescriptorParser {
	return btcdescriptor.NewDescriptorParser()
}

// FormatDescriptor formats a Descriptor object into a descriptor string (delegates to pkg).
func FormatDescriptor(desc *btcdescriptor.Descriptor) (string, error) {
	return btcdescriptor.FormatDescriptor(desc)
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
	case btcdescriptor.DescriptorTypeSH:
		// Legacy P2SH multisig
		redeemScript, err = b.deriveMultisigRedeemScript(parsed, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive multisig redeemScript for address %s at index %d: %w",
				address, addressIndex, err)
		}

	case btcdescriptor.DescriptorTypeSHWPKH:
		// P2SH-wrapped P2WPKH (BIP49 Nested SegWit)
		redeemScript, err = b.deriveP2WPKHRedeemScript(parsed, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive P2WPKH redeemScript for address %s at index %d: %w",
				address, addressIndex, err)
		}

	case btcdescriptor.DescriptorTypeSHWSH:
		// P2SH-wrapped P2WSH multisig (BIP49)
		// For sh(wsh(sortedmulti(...))), the redeemScript is OP_0 <witnessScript hash>
		// The witnessScript is the inner multisig script
		redeemScript, err = b.deriveP2WSHRedeemScript(parsed, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive P2WSH redeemScript for address %s at index %d: %w",
				address, addressIndex, err)
		}

	case btcdescriptor.DescriptorTypePKH,
		btcdescriptor.DescriptorTypeWPKH,
		btcdescriptor.DescriptorTypeTR,
		btcdescriptor.DescriptorTypeWSH,
		btcdescriptor.DescriptorTypeUnknown:
		// These descriptor types don't require redeemScript derivation:
		// - PKH: P2PKH addresses don't use redeemScript
		// - WPKH: Native SegWit (bech32) addresses don't use redeemScript
		// - TR: Taproot addresses don't use redeemScript
		// - WSH: Native SegWit multisig, not wrapped in P2SH
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
	parsed *btcdescriptor.Descriptor, addressIndex uint32,
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
	parsed *btcdescriptor.Descriptor, addressIndex uint32,
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

// deriveP2WSHRedeemScript derives the redeemScript for sh(wsh(...)) descriptors.
// For sh(wsh(sortedmulti(M, key1, key2, ...))):
//   - The witnessScript is the multisig script: OP_M <pk1> <pk2> ... OP_N OP_CHECKMULTISIG
//   - The redeemScript is: OP_0 <32-byte-witnessScript-hash>
//
// This is used for P2SH-wrapped P2WSH multisig (BIP49).
func (b *Bitcoin) deriveP2WSHRedeemScript(
	parsed *btcdescriptor.Descriptor, addressIndex uint32,
) ([]byte, error) {
	logger.Debug("Deriving P2WSH redeemScript",
		"descriptor_type", parsed.Type,
		"address_index", addressIndex,
		"num_keys", len(parsed.Keys))

	// Extract M and N from the descriptor script
	requiredSigs, totalSigs, err := b.extractMultisigParams(parsed.Script)
	if err != nil {
		return nil, fmt.Errorf("failed to extract multisig params: %w", err)
	}

	if len(parsed.Keys) != totalSigs {
		return nil, fmt.Errorf("key count mismatch: parsed %d keys but descriptor has %d",
			len(parsed.Keys), totalSigs)
	}

	logger.Info("Building P2WSH redeemScript",
		"required_sigs", requiredSigs,
		"total_sigs", totalSigs,
		"address_index", addressIndex)

	// Derive all public keys at the specified index
	pubKeys := make([][]byte, 0, totalSigs)
	for keyIdx, keyInfo := range parsed.Keys {
		pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive public key %d: %w", keyIdx, err)
		}
		pubKeys = append(pubKeys, pubKey)

		logger.Debug("Derived pubkey for P2WSH multisig",
			"key_index", keyIdx,
			"pubkey", hex.EncodeToString(pubKey))
	}

	// For sortedmulti, sort the public keys lexicographically (BIP67)
	if strings.Contains(parsed.Script, "sortedmulti") {
		logger.Debug("Sorting public keys for sortedmulti")
		sort.Slice(pubKeys, func(i, j int) bool {
			return bytes.Compare(pubKeys[i], pubKeys[j]) < 0
		})
	}

	// Build the witnessScript: OP_M <pk1> <pk2> ... <pkN> OP_N OP_CHECKMULTISIG
	witnessBuilder := txscript.NewScriptBuilder()
	witnessBuilder.AddInt64(int64(requiredSigs))

	for _, pubKey := range pubKeys {
		witnessBuilder.AddData(pubKey)
	}

	witnessBuilder.AddInt64(int64(totalSigs))
	witnessBuilder.AddOp(txscript.OP_CHECKMULTISIG)

	witnessScript, err := witnessBuilder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build witnessScript: %w", err)
	}

	logger.Debug("Built witnessScript",
		"script_len", len(witnessScript),
		"script_hex", hex.EncodeToString(witnessScript))

	// Hash the witnessScript: SHA256(witnessScript)
	witnessScriptHash := sha256.Sum256(witnessScript)

	logger.Debug("Computed witnessScript hash",
		"hash", hex.EncodeToString(witnessScriptHash[:]))

	// Build the redeemScript: OP_0 <32-byte-witnessScript-hash>
	// This is what goes into the P2SH
	redeemBuilder := txscript.NewScriptBuilder()
	redeemBuilder.AddOp(txscript.OP_0)
	redeemBuilder.AddData(witnessScriptHash[:])

	redeemScript, err := redeemBuilder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build P2WSH redeemScript: %w", err)
	}

	logger.Info("Built P2WSH redeemScript",
		"script_len", len(redeemScript),
		"script_hex", hex.EncodeToString(redeemScript),
		"witness_script_len", len(witnessScript))

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
	keyInfo btcdescriptor.DescriptorKey, addressIndex uint32,
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
