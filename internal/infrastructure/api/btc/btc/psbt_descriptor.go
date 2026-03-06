// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// This file contains descriptor derivation functions for PSBT.
// See psbt.go for the main PSBT documentation.
package btc

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// deriveRedeemScriptForAddress derives a redeemScript for a P2SH address by:
// 1. Finding the matching descriptor for the account
// 2. Searching for the address index within the descriptor range
// 3. Deriving the redeemScript at that index
//
// Performance characteristics:
//   - Lists all wallet descriptors (typically < 10)
//   - Iterates through up to 1,000 address indices per descriptor
//   - For each index, derives redeemScript and checks address match
//   - Time complexity: O(n * m) where n = descriptors, m = search range (1000)
//   - Typical time: < 1 second for addresses in first 100 indices
//   - Worst case: Several seconds for high-index addresses
//
// This is used as a fallback when Bitcoin Core's listunspent doesn't return
// redeemScript for descriptor-based addresses. It's called once per PSBT creation
// for P2SH multisig inputs.
func (b *Bitcoin) deriveRedeemScriptForAddress(
	address string, senderAccount domainAccount.AccountType,
) ([]byte, error) {
	// List all descriptors
	descriptorList, err := b.ListDescriptors(false)
	if err != nil {
		return nil, fmt.Errorf("failed to list descriptors: %w", err)
	}

	// Determine expected account index
	var expectedAccountIndex uint32
	switch senderAccount {
	case domainAccount.AccountTypeDeposit:
		expectedAccountIndex = 0
	case domainAccount.AccountTypePayment:
		expectedAccountIndex = 1
	case domainAccount.AccountTypeStored:
		expectedAccountIndex = 2
	case domainAccount.AccountTypeClient,
		domainAccount.AccountTypeAuthorization,
		domainAccount.AccountTypeAuth1, domainAccount.AccountTypeAuth2, domainAccount.AccountTypeAuth3,
		domainAccount.AccountTypeAuth4, domainAccount.AccountTypeAuth5, domainAccount.AccountTypeAuth6,
		domainAccount.AccountTypeAuth7, domainAccount.AccountTypeAuth8, domainAccount.AccountTypeAuth9,
		domainAccount.AccountTypeAuth10, domainAccount.AccountTypeAuth11, domainAccount.AccountTypeAuth12,
		domainAccount.AccountTypeAuth13, domainAccount.AccountTypeAuth14, domainAccount.AccountTypeAuth15,
		domainAccount.AccountTypeAnonymous, domainAccount.AccountTypeTest:
		return nil, fmt.Errorf("unsupported account type for this operation: %s", senderAccount.String())
	default:
		return nil, fmt.Errorf("unknown account type: %s", senderAccount.String())
	}

	logger.Debug("Searching for descriptor matching address",
		"address", address,
		"account", senderAccount.String(),
		"account_index", expectedAccountIndex)

	// Find matching descriptor and derive redeemScript
	for _, desc := range descriptorList {
		// Skip internal (change) descriptors
		if desc.Internal != nil && *desc.Internal {
			continue
		}

		// Check if descriptor matches the account
		if !b.descriptorMatchesAccountIndex(desc.Descriptor, expectedAccountIndex) {
			continue
		}

		logger.Debug("Found matching descriptor, searching for address",
			"descriptor_len", len(desc.Descriptor),
			"account_index", expectedAccountIndex)

		// Try to find the address by deriving up to 1000 addresses (default range)
		for i := range uint32(1000) {
			redeemScript, err := b.DeriveRedeemScriptFromDescriptor(desc.Descriptor, address, i)
			if err != nil {
				// This index doesn't work, try next
				continue
			}

			// Verify the derived redeemScript matches the address
			derivedAddr, err := b.deriveP2SHAddressFromRedeemScript(redeemScript)
			if err != nil {
				continue
			}

			if derivedAddr == address {
				logger.Info("Successfully derived redeemScript for address",
					"address", address,
					"descriptor_index", i,
					"script_len", len(redeemScript))
				return redeemScript, nil
			}
		}
	}

	return nil, fmt.Errorf("no matching descriptor found for address %s (account=%s)", address, senderAccount.String())
}

// deriveP2SHAddressFromRedeemScript derives a P2SH address from a redeemScript
func (b *Bitcoin) deriveP2SHAddressFromRedeemScript(redeemScript []byte) (string, error) {
	// Hash the redeemScript
	scriptHash := btcutil.Hash160(redeemScript)

	// Create P2SH address from hash
	// Note: Use NewAddressScriptHashFromHash since we already hashed the script
	address, err := btcutil.NewAddressScriptHashFromHash(scriptHash, b.chainConf)
	if err != nil {
		return "", fmt.Errorf("failed to create P2SH address: %w", err)
	}

	return address.EncodeAddress(), nil
}

// deriveNativeSegWitAddress derives a Bech32 address from a Native SegWit descriptor.
// Supports both P2WPKH (single-sig) and P2WSH (multisig) descriptors.
//
// For P2WPKH:
//   - Derives public key at addressIndex
//   - Hashes public key (HASH160)
//   - Creates Bech32 address from 20-byte witness program
//   - Format: bc1q<20-byte-hash> (42 chars mainnet), bcrt1q<20-byte-hash> (regtest)
//
// For P2WSH:
//   - Derives witness script (multisig script) at addressIndex
//   - Hashes witness script (SHA256)
//   - Creates Bech32 address from 32-byte witness program
//   - Format: bc1q<32-byte-hash> (62 chars mainnet), bcrt1q<32-byte-hash> (regtest)
func (b *Bitcoin) deriveNativeSegWitAddress(parsed *domainWallet.Descriptor, addressIndex uint32) (string, error) {
	switch parsed.Type {
	case domainWallet.DescriptorTypeWPKH:
		// P2WPKH: single-sig Native SegWit
		if len(parsed.Keys) != 1 {
			return "", fmt.Errorf("P2WPKH descriptor must have exactly 1 key, got %d", len(parsed.Keys))
		}

		// Derive the public key at this index
		pubKey, err := b.derivePublicKeyFromDescriptorKey(parsed.Keys[0], addressIndex)
		if err != nil {
			return "", fmt.Errorf("failed to derive public key: %w", err)
		}

		// Hash the public key (HASH160 = RIPEMD160(SHA256(pubkey)))
		pubKeyHash := btcutil.Hash160(pubKey)

		// Create P2WPKH address
		address, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, b.chainConf)
		if err != nil {
			return "", fmt.Errorf("failed to create P2WPKH address: %w", err)
		}

		return address.EncodeAddress(), nil

	case domainWallet.DescriptorTypeWSH:
		// P2WSH: multisig Native SegWit
		// Derive the witness script (multisig script)
		witnessScript, err := b.deriveWitnessScriptFromDescriptor(parsed, addressIndex)
		if err != nil {
			return "", fmt.Errorf("failed to derive witness script: %w", err)
		}

		// Hash the witness script (SHA256)
		witnessScriptHash := sha256.Sum256(witnessScript)

		// Create P2WSH address
		address, err := btcutil.NewAddressWitnessScriptHash(witnessScriptHash[:], b.chainConf)
		if err != nil {
			return "", fmt.Errorf("failed to create P2WSH address: %w", err)
		}

		return address.EncodeAddress(), nil

	default:
		return "", fmt.Errorf("unsupported native SegWit descriptor type: %s", parsed.Type)
	}
}

// deriveWitnessScriptFromDescriptor derives a witness script (multisig script) from a P2WSH descriptor.
// The witness script defines the spending conditions (e.g., 2-of-3 multisig).
//
// Format for sortedmulti(M, key1, key2, ...):
//   - M (threshold as OP_N)
//   - For each key: <pubkey>
//   - N (total keys as OP_N)
//   - OP_CHECKMULTISIG
//
// Example for 2-of-3:
//
//	OP_2 <pubkey1> <pubkey2> <pubkey3> OP_3 OP_CHECKMULTISIG
func (b *Bitcoin) deriveWitnessScriptFromDescriptor(
	parsed *domainWallet.Descriptor,
	addressIndex uint32,
) ([]byte, error) {
	if len(parsed.Keys) < 2 {
		return nil, fmt.Errorf("multisig descriptor must have at least 2 keys, got %d", len(parsed.Keys))
	}

	// Derive all public keys at this index
	pubKeys := make([][]byte, len(parsed.Keys))
	for i, keyInfo := range parsed.Keys {
		pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to derive public key %d: %w", i, err)
		}
		pubKeys[i] = pubKey
	}

	// Sort public keys lexicographically (BIP67 sortedmulti)
	// This ensures deterministic key ordering
	sortedPubKeys := make([][]byte, len(pubKeys))
	copy(sortedPubKeys, pubKeys)
	sort.Slice(sortedPubKeys, func(i, j int) bool {
		return bytes.Compare(sortedPubKeys[i], sortedPubKeys[j]) < 0
	})

	// Extract threshold from descriptor script
	// Format: wsh(sortedmulti(M, key1, key2, ...))
	// We need to parse "sortedmulti(M," to extract M
	threshold := 2 // Default for 2-of-3
	re := regexp.MustCompile(`sortedmulti\((\d+),`)
	matches := re.FindStringSubmatch(parsed.Script)
	if len(matches) > 1 {
		parsedThreshold, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse threshold from descriptor: %w", err)
		}
		threshold = parsedThreshold
	}

	// Build witness script using txscript
	builder := txscript.NewScriptBuilder()
	builder.AddInt64(int64(threshold)) // M (threshold)

	for _, pubKey := range sortedPubKeys {
		builder.AddData(pubKey)
	}

	builder.AddInt64(int64(len(sortedPubKeys))) // N (total keys)
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	witnessScript, err := builder.Script()
	if err != nil {
		return nil, fmt.Errorf("failed to build witness script: %w", err)
	}

	return witnessScript, nil
}

// findAddressIndexInDescriptor uses Bitcoin Core's deriveaddresses RPC to find
// the index of an address within a descriptor's range.
//
// Performance characteristics:
//   - Time complexity: O(n) where n is the address index
//   - Searches in chunks of 100 addresses at a time
//   - Maximum search range: 10,000 addresses
//   - Early termination when address is found
//   - Average case: ~50 RPC calls for index 5000 (50 chunks * 100 addresses)
//   - Best case: 1 RPC call (address in first chunk)
//   - Worst case: 100 RPC calls (address at index 9,999)
//
// Optimization opportunities:
//   - Use Bitcoin Core's native address index lookup if available
//   - Cache descriptor derivation results
//   - Binary search if descriptor supports random access
func (b *Bitcoin) findAddressIndexInDescriptor(descriptor string, targetAddress string) (uint32, error) {
	// Search in chunks to avoid deriving too many addresses at once
	const chunkSize = 100
	const maxSearchRange = 10000 // Maximum range to search

	for startIdx := uint32(0); startIdx < maxSearchRange; startIdx += chunkSize {
		endIdx := startIdx + chunkSize - 1

		// Call deriveaddresses with range
		addresses, err := b.deriveAddressesFromDescriptor(descriptor, startIdx, endIdx)
		if err != nil {
			return 0, fmt.Errorf("failed to derive addresses [%d,%d]: %w", startIdx, endIdx, err)
		}

		// Search for target address in the chunk
		for i, addr := range addresses {
			if addr == targetAddress {
				return startIdx + uint32(i), nil
			}
		}
	}

	return 0, fmt.Errorf("address %s not found in descriptor range [0,%d]", targetAddress, maxSearchRange)
}

// deriveAddressesFromDescriptor derives addresses from a descriptor within the specified range.
func (b *Bitcoin) deriveAddressesFromDescriptor(descriptor string, startIdx, endIdx uint32) ([]string, error) {
	addresses, err := b.pkgrpc.DeriveAddresses(descriptor, startIdx, endIdx)
	if err != nil {
		return nil, fmt.Errorf("deriveaddresses RPC failed: %w", err)
	}

	logger.Debug("Derived addresses from descriptor",
		"range", fmt.Sprintf("[%d,%d]", startIdx, endIdx),
		"count", len(addresses))

	return addresses, nil
}
