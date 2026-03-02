// Package btc provides Bitcoin (BTC) infrastructure implementations.
//
// This file contains BIP32 derivation functions for PSBT.
// See psbt.go for the main PSBT documentation.
package btc

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// addBIP32DerivationForInput adds BIP32 derivation path information to a PSBT input.
// This is required for descriptor-based signing to work correctly.
// It extracts the address from the scriptPubKey, queries Bitcoin Core for derivation info,
// and adds it to the PSBT.
func (b *Bitcoin) addBIP32DerivationForInput(
	updater *psbt.Updater,
	scriptPubKey []byte,
	inputIndex int,
	senderAccount domainAccount.AccountType,
) error {
	// Extract address from scriptPubKey
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptPubKey, b.chainConf)
	if err != nil {
		return fmt.Errorf("failed to extract address from scriptPubKey: %w", err)
	}
	if len(addrs) == 0 {
		return errors.New("no address found in scriptPubKey")
	}

	// Get address info to determine the derivation path
	addressStr := addrs[0].EncodeAddress()
	logger.Info("Adding BIP32 derivation for input",
		"input_index", inputIndex,
		"address", addressStr,
		"sender_account", senderAccount.String())

	addressInfo, err := b.GetAddressInfo(addressStr)
	if err != nil {
		logger.Error("GetAddressInfo failed", "address", addressStr, "error", err)
		return fmt.Errorf("failed to get address info for %s: %w", addressStr, err)
	}

	logger.Info("Got address info",
		"address", addressStr,
		"hd_key_path", addressInfo.HDKeyPath,
		"hd_fingerprint", addressInfo.HDMasterFingerprint,
		"pubkey_len", len(addressInfo.Pubkey),
		"is_watch_only", addressInfo.IsWatchOnly,
		"is_script", addressInfo.IsScript)

	// Check if address has derivation path
	// For watch-only multisig descriptors, HDKeyPath and PubKey may be empty
	// Bitcoin Core doesn't populate HDKeyPath for ranged multisig descriptors
	// In this case, derive BIP32 information from the descriptor
	if addressInfo.HDKeyPath == "" {
		if addressInfo.IsScript {
			// BCH uses legacy (non-descriptor) wallets and doesn't support descriptor-based signing
			// Skip BIP32 derivation from descriptor for BCH
			if b.coinTypeCode == domainCoin.BCH {
				logger.Debug("BCH: Skipping BIP32 derivation from descriptor (legacy wallet, not supported)",
					"input_index", inputIndex,
					"address", addressStr,
					"sender_account", senderAccount.String())
				return nil
			}

			logger.Info("Deriving BIP32 information from descriptor for multisig address",
				"address", addressStr,
				"input_index", inputIndex)

			// Derive BIP32 derivation paths from descriptor
			err := b.addBIP32DerivationFromDescriptor(updater, addressStr, inputIndex, senderAccount)
			if err != nil {
				logger.Error("Failed to derive BIP32 from descriptor",
					"address", addressStr,
					"error", err)
				return fmt.Errorf("failed to derive BIP32 from descriptor: %w", err)
			}
			logger.Info("Successfully added BIP32 derivation from descriptor",
				"address", addressStr,
				"input_index", inputIndex)
			return nil
		}
		return errors.New("address has no HD key path (not from descriptor wallet?)")
	}
	// For P2TR (Taproot), Bitcoin Core doesn't return pubkey in getaddressinfo.
	// Taproot uses x-only pubkeys (32 bytes) which are incompatible with AddInBip32Derivation
	// (which expects 33-byte compressed pubkeys). For P2TR single-sig, the keygen wallet
	// can sign without BIP32 derivation info in PSBT, so we skip it.
	if txscript.IsPayToTaproot(scriptPubKey) {
		logger.Debug("Skipping BIP32 derivation for P2TR input (Taproot uses x-only pubkeys)",
			"input_index", inputIndex,
			"address", addressStr)
		return nil
	}

	// Parse the public key for non-Taproot addresses
	if addressInfo.Pubkey == "" {
		return errors.New("address has no public key")
	}
	pubKeyBytes, err := hex.DecodeString(addressInfo.Pubkey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Get the correct fingerprint from imported descriptors, not from Bitcoin Core's internal wallet
	// Query list of descriptors to find the one that matches this address and account
	fingerprint, _, err := b.getDescriptorInfoForAddress(addressStr, addressInfo.HDKeyPath, senderAccount)
	if err != nil {
		return fmt.Errorf("failed to get descriptor info: %w", err)
	}

	// Parse the full derivation path
	// Use addressInfo.HDKeyPath directly as it already contains the full path from getaddressinfo
	derivationPath, err := parseDerivationPath(addressInfo.HDKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse derivation path: %w", err)
	}

	// Add BIP32 derivation to PSBT
	// updater.AddInBip32Derivation signature: (fingerprint uint32, path []uint32, pubkey []byte, inputIndex int)
	if err := updater.AddInBip32Derivation(fingerprint, derivationPath, pubKeyBytes, inputIndex); err != nil {
		return fmt.Errorf("failed to add BIP32 derivation to PSBT: %w", err)
	}

	logger.Debug("Added BIP32 derivation to PSBT input",
		"input_index", inputIndex,
		"address", addrs[0].EncodeAddress(),
		"path", addressInfo.HDKeyPath,
		"fingerprint", addressInfo.HDMasterFingerprint,
	)

	return nil
}

// addBIP32DerivationFromDescriptor adds BIP32 derivation paths to PSBT from descriptor.
// This is required for descriptor-based signing when Bitcoin Core doesn't return HDKeyPath.
func (b *Bitcoin) addBIP32DerivationFromDescriptor(
	updater *psbt.Updater,
	address string,
	inputIndex int,
	senderAccount domainAccount.AccountType,
) error {
	logger.Debug("Deriving BIP32 paths from descriptor",
		"address", address,
		"account", senderAccount.String())

	// List all descriptors to find the original descriptor with xpubs
	descriptorList, err := b.ListDescriptors(false)
	if err != nil {
		return fmt.Errorf("failed to list descriptors: %w", err)
	}

	// Get address info to check which descriptor this address belongs to
	addressInfo, err := b.GetAddressInfo(address)
	if err != nil {
		return fmt.Errorf("failed to get address info: %w", err)
	}

	if addressInfo.Desc == "" {
		return errors.New("address info does not contain descriptor (not a descriptor wallet address?)")
	}

	logger.Debug("Got descriptor from getaddressinfo (with raw keys)",
		"address", address,
		"descriptor_len", len(addressInfo.Desc))

	// Find and process matching wallet descriptor
	for _, desc := range descriptorList.Descriptors {
		foundIndex, err := b.findAddressIndexInDescriptor(desc.Desc, address)
		if err == nil {
			// Found matching descriptor - process it
			return b.processMatchedDescriptor(updater, address, inputIndex, desc.Desc, foundIndex, desc.Internal)
		}
	}

	return fmt.Errorf("no matching wallet descriptor found for address %s", address)
}

// processMatchedDescriptor handles BIP32 derivation for a matched descriptor.
func (b *Bitcoin) processMatchedDescriptor(
	updater *psbt.Updater,
	address string,
	inputIndex int,
	walletDescriptor string,
	addressIndex uint32,
	internal *bool,
) error {
	isInternal := internal != nil && *internal
	logger.Info("Found matching wallet descriptor",
		"address", address,
		"address_index", addressIndex,
		"is_internal", isInternal,
		"descriptor_prefix", walletDescriptor[:80]+"...")

	// Verify with Bitcoin Core
	b.logDescriptorVerification(walletDescriptor, addressIndex, address)

	// Parse and validate descriptor
	parsed, err := b.parseAndValidateDescriptor(walletDescriptor)
	if err != nil {
		return err
	}

	// Verify address matches descriptor
	if err := b.verifyAddressMatchesDescriptor(parsed, walletDescriptor, address, addressIndex); err != nil {
		return err
	}

	logger.Info("Verified address matches descriptor, adding BIP32 derivation for all keys",
		"address", address,
		"descriptor_type", parsed.Type,
		"descriptor_index", addressIndex,
		"num_keys", len(parsed.Keys))

	// Add BIP32 derivation for each key
	if err := b.addBIP32DerivationForKeys(updater, parsed, inputIndex, addressIndex); err != nil {
		return err
	}

	logger.Info("Successfully added BIP32 derivation for all descriptor keys",
		"address", address,
		"descriptor_type", parsed.Type,
		"input", inputIndex,
		"num_keys", len(parsed.Keys))
	return nil
}

// logDescriptorVerification logs Bitcoin Core deriveaddresses verification.
func (b *Bitcoin) logDescriptorVerification(descriptor string, addressIndex uint32, targetAddress string) {
	addresses, err := b.deriveAddressesFromDescriptor(descriptor, addressIndex, addressIndex)
	if err == nil && len(addresses) > 0 {
		logger.Info("Bitcoin Core deriveaddresses verification",
			"descriptor_index", addressIndex,
			"bitcoin_core_address", addresses[0],
			"target_address", targetAddress,
			"match", addresses[0] == targetAddress)
	} else {
		logger.Warn("Failed to verify with Bitcoin Core deriveaddresses", "error", err)
	}
}

// parseAndValidateDescriptor parses descriptor and validates its type.
func (*Bitcoin) parseAndValidateDescriptor(walletDescriptor string) (*domainWallet.Descriptor, error) {
	parser := NewDescriptorParser()
	parsed, err := parser.Parse(walletDescriptor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wallet descriptor: %w", err)
	}

	// Support all descriptor types: P2SH-wrapped and Native SegWit
	supportedTypes := map[domainWallet.DescriptorType]bool{
		domainWallet.DescriptorTypeSH:     true,
		domainWallet.DescriptorTypeSHWPKH: true,
		domainWallet.DescriptorTypeSHWSH:  true,
		domainWallet.DescriptorTypeWPKH:   true,
		domainWallet.DescriptorTypeWSH:    true,
	}

	if !supportedTypes[parsed.Type] {
		return nil, fmt.Errorf(
			"unsupported descriptor type: %s (only sh/sh(wpkh)/sh(wsh)/wpkh/wsh supported)",
			parsed.Type,
		)
	}

	return parsed, nil
}

// verifyAddressMatchesDescriptor verifies the derived address matches target.
func (b *Bitcoin) verifyAddressMatchesDescriptor(
	parsed *domainWallet.Descriptor,
	walletDescriptor, address string,
	addressIndex uint32,
) error {
	var derivedAddr string
	var err error

	if parsed.Type == domainWallet.DescriptorTypeWSH || parsed.Type == domainWallet.DescriptorTypeWPKH {
		// Native SegWit - derive Bech32 address from witness script/program
		derivedAddr, err = b.deriveNativeSegWitAddress(parsed, addressIndex)
		if err != nil {
			return fmt.Errorf("failed to derive native SegWit address: %w", err)
		}
	} else {
		// P2SH - derive P2SH address from redeem script
		redeemScript, err := b.DeriveRedeemScriptFromDescriptor(walletDescriptor, address, addressIndex)
		if err != nil {
			return fmt.Errorf("failed to derive redeemScript: %w", err)
		}

		derivedAddr, err = b.deriveP2SHAddressFromRedeemScript(redeemScript)
		if err != nil {
			return fmt.Errorf("failed to derive P2SH address from redeemScript: %w", err)
		}
	}

	if derivedAddr != address {
		return fmt.Errorf("derived address %s does not match target %s", derivedAddr, address)
	}

	return nil
}

// addBIP32DerivationForKeys adds BIP32 derivation for each key in the descriptor.
func (b *Bitcoin) addBIP32DerivationForKeys(
	updater *psbt.Updater,
	parsed *domainWallet.Descriptor,
	inputIndex int,
	addressIndex uint32,
) error {
	for keyIdx, keyInfo := range parsed.Keys {
		pubKey, err := b.derivePublicKeyFromDescriptorKey(keyInfo, addressIndex)
		if err != nil {
			return fmt.Errorf("failed to derive public key %d: %w", keyIdx, err)
		}

		fingerprint, err := hex.DecodeString(keyInfo.Fingerprint)
		if err != nil {
			return fmt.Errorf("invalid fingerprint: %w", err)
		}

		// Build the full derivation path
		relativePath := keyInfo.DerivationPath
		if relativePath == "" {
			relativePath = fmt.Sprintf("/%d", addressIndex)
		} else {
			relativePath = strings.ReplaceAll(relativePath, "/*", fmt.Sprintf("/%d", addressIndex))
		}
		fullPath := keyInfo.OriginPath + relativePath

		pathIndices, err := parseDerivationPath(fullPath)
		if err != nil {
			return fmt.Errorf("failed to parse derivation path %s: %w", fullPath, err)
		}

		fingerprintUint32 := binary.LittleEndian.Uint32(fingerprint)
		if err := updater.AddInBip32Derivation(fingerprintUint32, pathIndices, pubKey, inputIndex); err != nil {
			return fmt.Errorf("failed to add BIP32 derivation for key %d: %w", keyIdx, err)
		}

		logger.Debug("Added BIP32 derivation for descriptor key",
			"input", inputIndex,
			"descriptor_type", parsed.Type,
			"key_index", keyIdx,
			"pubkey", hex.EncodeToString(pubKey),
			"fingerprint", keyInfo.Fingerprint,
			"full_path", fullPath,
			"origin", keyInfo.OriginPath,
			"relative", relativePath)
	}

	return nil
}

// parseDerivationPath converts a BIP32 derivation path string to a uint32 array.
// Input format: "m/44h/1h/1h/0/0" or "m/44'/1'/1'/0/0"
// Output format: []uint32 with hardened keys having high bit set (0x80000000)
func parseDerivationPath(path string) ([]uint32, error) {
	// Remove "m/" prefix if present
	if len(path) >= 2 && path[:2] == "m/" {
		path = path[2:]
	}

	// Split by "/"
	parts := make([]string, 0)
	for _, part := range splitPath(path) {
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return nil, errors.New("empty derivation path")
	}

	// Parse each component
	result := make([]uint32, len(parts))
	for i, part := range parts {
		// Check if hardened (ends with 'h' or apostrophe)
		hardened := false
		if len(part) > 0 && (part[len(part)-1] == 'h' || part[len(part)-1] == '\'') {
			hardened = true
			part = part[:len(part)-1]
		}

		// Parse the number
		var num uint32
		_, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			return nil, fmt.Errorf("failed to parse path component %q: %w", part, err)
		}

		// Set hardened bit if needed
		if hardened {
			num |= 0x80000000 // Set high bit for hardened derivation
		}

		result[i] = num
	}

	return result, nil
}

// splitPath splits a derivation path by "/" separator
func splitPath(path string) []string {
	parts := make([]string, 0)
	current := ""
	for _, ch := range path {
		if ch == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// getDescriptorInfoForAddress queries loaded descriptors to find the correct fingerprint
// for the given address. This is necessary because getaddressinfo returns Bitcoin Core's
// internal wallet fingerprint, not the imported descriptor's fingerprint.
func (b *Bitcoin) getDescriptorInfoForAddress(
	address, fullPath string, senderAccount domainAccount.AccountType,
) (uint32, string, error) {
	// Call listdescriptors to get public descriptors
	result, err := b.ListDescriptors(false)
	if err != nil {
		return 0, "", fmt.Errorf("failed to call listdescriptors: %w", err)
	}

	// Extract account path from full path
	accountPath, accountPathApostrophe, err := extractAccountPath(fullPath)
	if err != nil {
		return 0, "", err
	}

	expectedAccountIndex := senderAccount.BIP44AccountIndex()
	logger.Debug("Searching for descriptor",
		"account", senderAccount.String(),
		"expected_account_index", expectedAccountIndex,
		"account_path", accountPath,
		"normalized_path", accountPathApostrophe,
	)

	// Find the descriptor that matches this account path
	for _, desc := range result.Descriptors {
		if !descriptorContainsPath(desc.Desc, accountPath, accountPathApostrophe) {
			continue
		}

		// Skip internal (change) descriptors
		if desc.Internal != nil && *desc.Internal {
			logger.Debug("Skipping internal descriptor", "desc", desc.Desc[:50])
			continue
		}

		// Try to extract and match descriptor info
		fingerprint, basePath, ok := parseDescriptorFingerprintAndPath(desc.Desc)
		if !ok {
			continue
		}

		// Verify account index matches
		if matchDescriptorAccountIndex(basePath, expectedAccountIndex, address, senderAccount) {
			return fingerprint, basePath, nil
		}
	}

	return 0, "", fmt.Errorf(
		"no matching descriptor found for address %s with path %s and account %s (expected account index %d)",
		address, fullPath, senderAccount.String(), expectedAccountIndex,
	)
}

// extractAccountPath extracts account-level path from full BIP44 path.
func extractAccountPath(fullPath string) (accountPath, accountPathApostrophe string, err error) {
	pathParts := strings.Split(strings.TrimPrefix(fullPath, "m/"), "/")
	if len(pathParts) < 3 {
		return "", "", fmt.Errorf("invalid derivation path: %s", fullPath)
	}
	// Get the account-level path (e.g., "44h/1h/1h")
	accountPath = strings.Join(pathParts[:len(pathParts)-2], "/")
	// Both 'h' notation and apostrophe notation are valid
	accountPathApostrophe = strings.ReplaceAll(accountPath, "h", "'")
	return accountPath, accountPathApostrophe, nil
}

// descriptorContainsPath checks if descriptor contains the account path.
func descriptorContainsPath(desc, accountPath, accountPathApostrophe string) bool {
	return strings.Contains(desc, accountPath) || strings.Contains(desc, accountPathApostrophe)
}

// parseDescriptorFingerprintAndPath extracts fingerprint and base path from descriptor.
// Returns (fingerprint, basePath, ok).
func parseDescriptorFingerprintAndPath(desc string) (uint32, string, bool) {
	// Extract fingerprint - Format: [fingerprint/44h/1h/1h]
	start := strings.Index(desc, "[")
	end := strings.Index(desc, "/")
	if start == -1 || end == -1 || end <= start {
		return 0, "", false
	}

	fingerprintHex := desc[start+1 : end]
	fingerprintBytes, err := hex.DecodeString(fingerprintHex)
	if err != nil || len(fingerprintBytes) != 4 {
		return 0, "", false
	}

	// Convert to uint32 (big-endian)
	fingerprint := uint32(fingerprintBytes[0])<<24 |
		uint32(fingerprintBytes[1])<<16 |
		uint32(fingerprintBytes[2])<<8 |
		uint32(fingerprintBytes[3])

	// Extract base path
	pathStart := strings.Index(desc, "/")
	pathEnd := strings.Index(desc, "]")
	if pathStart == -1 || pathEnd == -1 || pathEnd <= pathStart {
		return 0, "", false
	}
	basePath := "m" + desc[pathStart:pathEnd]

	return fingerprint, basePath, true
}

// matchDescriptorAccountIndex verifies the account index in basePath matches expected.
func matchDescriptorAccountIndex(
	basePath string,
	expectedAccountIndex uint32,
	address string,
	senderAccount domainAccount.AccountType,
) bool {
	basePathParts := strings.Split(strings.TrimPrefix(basePath, "m/"), "/")
	if len(basePathParts) < 3 {
		return false
	}

	// Parse the account component (e.g., "1h" -> 1)
	accountComponent := basePathParts[2]
	accountComponent = strings.TrimSuffix(accountComponent, "h")
	accountComponent = strings.TrimSuffix(accountComponent, "'")

	var parsedAccountIndex uint32
	_, err := fmt.Sscanf(accountComponent, "%d", &parsedAccountIndex)
	if err != nil || parsedAccountIndex != expectedAccountIndex {
		logger.Debug("Account index mismatch",
			"parsed", parsedAccountIndex,
			"expected", expectedAccountIndex,
			"base_path", basePath,
		)
		return false
	}

	logger.Info("Successfully matched descriptor for sender account",
		"address", address,
		"account", senderAccount.String(),
		"base_path", basePath,
		"account_index", parsedAccountIndex,
	)
	return true
}
