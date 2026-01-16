package btc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ListUnspentResult is response type of PRC `listunspent`
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Label         string  `json:"label"`
	RedeemScript  string  `json:"redeemScript"`
	WitnessScript string  `json:"witnessScript"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
	Desc          string  `json:"desc"`
	Safe          bool    `json:"safe"`
}

// listUnspentRaw calls RPC `listunspent` and returns raw ListUnspentResult without conversion to DTO.
// This is used internally to avoid redundant data conversions.
func (b *Bitcoin) listUnspentRaw(confirmationNum uint64) ([]ListUnspentResult, error) {
	input, err := json.Marshal(confirmationNum)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(): %w", err)
	}
	rawResult, err := b.Client.RawRequest("listunspent", []json.RawMessage{input})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(listunspent): %w", err)
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal(rawResult, &listunspentResult)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(): %w", err)
	}

	return listunspentResult, nil
}

// ListUnspent call RPC `listunspent`
func (b *Bitcoin) ListUnspent(confirmationNum uint64) ([]dtobtc.UnspentOutput, error) {
	logger.Debug("call ListUnspent()", "confirmation", b.confirmationBlock)

	listunspentResult, err := b.listUnspentRaw(confirmationNum)
	if err != nil {
		return nil, err
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return ToUnspentOutputList(listunspentResult, b)
}

// ListUnspentByAccount gets listunspent by account.
//
// This method uses two strategies to find unspent outputs for an account:
// 1. Label-based lookup: Fast path using getaddressesbylabel (works for single-sig and labeled multisig)
// 2. Descriptor-based matching: Fallback for multisig descriptors (ranged descriptors without labels)
//
// For multisig wallets using ranged descriptors, Bitcoin Core doesn't support labels on the descriptor itself.
// While individual addresses are labeled after import, there may be cases where labeling fails or addresses
// are derived after import. In such cases, we fall back to matching addresses by their descriptors.
func (b *Bitcoin) ListUnspentByAccount(
	accountType domainAccount.AccountType, confirmationNum uint64,
) ([]dtobtc.UnspentOutput, error) {
	label := accountType.String()
	addrs, err := b.GetAddressesByLabel(label)

	// If we found addresses by label, use them (fast path)
	if err == nil && len(addrs) > 0 {
		logger.Debug("found labeled addresses for account",
			"account_type", accountType,
			"label", label,
			"address_count", len(addrs))

		unspentList, err := b.listUnspentByAccount(addrs, confirmationNum)
		if err != nil {
			return nil, fmt.Errorf("fail to call btc.listUnspentByAccount(): %w", err)
		}

		// If we found UTXOs, return them
		if len(unspentList) > 0 {
			logger.Debug("found unspent outputs via label lookup",
				"account_type", accountType,
				"unspent_count", len(unspentList),
				"confirmation_required", confirmationNum)

			// sort amount by ascending (small to big)
			sort.Slice(unspentList, func(i, j int) bool {
				return unspentList[i].Amount < unspentList[j].Amount
			})

			return ToUnspentOutputList(unspentList, b)
		}

		// Labeled addresses exist but have no UTXOs - fall back to descriptor matching
		logger.Info("labeled addresses found but no UTXOs, falling back to descriptor matching",
			"account_type", accountType,
			"label", label,
			"address_count", len(addrs))
	}

	// No addresses found by label (or error occurred) - fall back to descriptor matching
	if err != nil {
		logger.Info("label lookup failed, falling back to descriptor matching",
			"account_type", accountType,
			"label", label,
			"error", err)
	} else {
		logger.Info("no labeled addresses found, falling back to descriptor matching",
			"account_type", accountType,
			"label", label)
	}

	unspentList, err := b.listUnspentByDescriptorMatching(accountType, confirmationNum)
	if err != nil {
		return nil, fmt.Errorf("fail to call btc.listUnspentByDescriptorMatching(): %w", err)
	}

	if len(unspentList) == 0 {
		logger.Info("no unspent outputs found for account",
			"account_type", accountType,
			"method", "descriptor_matching")
		return nil, nil
	}

	logger.Debug("found unspent outputs via descriptor matching",
		"account_type", accountType,
		"unspent_count", len(unspentList),
		"confirmation_required", confirmationNum)

	// sort amount by ascending (small to big)
	sort.Slice(unspentList, func(i, j int) bool {
		return unspentList[i].Amount < unspentList[j].Amount
	})

	return ToUnspentOutputList(unspentList, b)
}

// GetUnspentListAddrs returns address from unspentList
func (*Bitcoin) GetUnspentListAddrs(
	unspentList []dtobtc.UnspentOutput, accountType domainAccount.AccountType,
) []string {
	addrs := make([]string, 0, len(unspentList))
	for _, unspent := range unspentList {
		if unspent.Label != accountType.String() {
			logger.Warn("listUnspentByAccount() returns address for wrong account",
				"got", unspent.Label,
				"want", accountType.String())
		}
		addrs = append(addrs, unspent.Address)
	}
	return addrs
}

func (b *Bitcoin) listUnspentByAccount(addrs []btcutil.Address, confirmationNum uint64) ([]ListUnspentResult, error) {
	input1, err := json.Marshal(confirmationNum)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(confirmationBlock): %w", err)
	}

	input2, err := json.Marshal(uint64(9999999))
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(9999999): %w", err)
	}

	// address
	strAddrs := make([]string, len(addrs))
	for idx, addr := range addrs {
		strAddrs[idx] = addr.String()
	}

	input3, err := json.Marshal(strAddrs)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Marchal(addresses): %w", err)
	}

	rawResult, err := b.Client.RawRequest("listunspent", []json.RawMessage{input1, input2, input3})
	if err != nil {
		return nil, fmt.Errorf("fail to call json.RawRequest(listunspent): %w", err)
	}

	var listunspentResult []ListUnspentResult
	err = json.Unmarshal(rawResult, &listunspentResult)
	if err != nil {
		return nil, fmt.Errorf("fail to call json.Unmarshal(rawResult): %w", err)
	}

	if len(listunspentResult) == 0 {
		return nil, nil
	}

	return listunspentResult, nil
}

// listUnspentByDescriptorMatching finds unspent outputs for an account by matching descriptors.
//
// This method is used as a fallback when getaddressesbylabel returns no results,
// which happens for multisig accounts using ranged descriptors (descriptors with wildcards).
//
// The matching algorithm:
//  1. Get all unspent outputs from the wallet
//  2. Get all descriptors imported into the wallet
//  3. Filter descriptors by account type (using BIP44 account index)
//  4. For each unspent output, check if it matches the account's descriptors
//  5. Return matching unspent outputs
//
// This approach works because:
//   - Every address derived from a descriptor has that descriptor in its address info
//   - The descriptor includes the BIP44 path with the account index
//   - We can match descriptors by comparing the base descriptor (before the address index)
//
// Performance characteristics:
//   - Fetches ALL unspent outputs from wallet (1 RPC call)
//   - Lists all descriptors (1 RPC call, typically < 10 descriptors)
//   - For each unspent, calls getaddressinfo to check account (N RPC calls)
//   - Time complexity: O(n) where n = total unspent outputs in wallet
//   - Typical time: < 1 second for wallets with < 100 UTXOs
//   - Performance degrades with many UTXOs (1000+ UTXOs may take several seconds)
//
// Optimization: This is only used as a fallback when label-based lookup fails.
// For best performance, ensure addresses are properly labeled during descriptor import.
func (b *Bitcoin) listUnspentByDescriptorMatching(
	accountType domainAccount.AccountType,
	confirmationNum uint64,
) ([]ListUnspentResult, error) {
	// Get all unspent outputs (no address filter) - use raw result to avoid conversion
	allUnspent, err := b.listUnspentRaw(confirmationNum)
	if err != nil {
		return nil, fmt.Errorf("failed to list all unspent outputs: %w", err)
	}

	if len(allUnspent) == 0 {
		return nil, nil
	}

	logger.Debug("checking unspent outputs for account",
		"account_type", accountType,
		"total_unspent", len(allUnspent))

	// Get all descriptors from the wallet
	descriptorList, err := b.ListDescriptors(false) // false = public descriptors (watch-only wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to list descriptors: %w", err)
	}

	if len(descriptorList.Descriptors) == 0 {
		logger.Warn("no descriptors found in wallet")
		return nil, nil
	}

	// Get the expected BIP44 account index for this account type
	expectedAccountIndex := accountType.BIP44AccountIndex()

	// Check if any descriptor matches this account's BIP44 index
	// Descriptor format: "sh(wsh(sortedmulti(2,[fingerprint/path]xpub.../0/*,...))"
	// The path contains the account index (e.g., 44h/1h/0h where 0 is the account index)
	hasMatchingDescriptor := false
	for _, desc := range descriptorList.Descriptors {
		// Skip internal (change) descriptors - we only want external (receiving) descriptors
		if desc.Internal != nil && *desc.Internal {
			continue
		}

		// Check if this descriptor matches the account index
		if b.descriptorMatchesAccountIndex(desc.Desc, expectedAccountIndex) {
			hasMatchingDescriptor = true
			break
		}
	}

	if !hasMatchingDescriptor {
		logger.Warn("no matching descriptors found for account",
			"account_type", accountType,
			"expected_account_index", expectedAccountIndex)
		return nil, nil
	}

	logger.Debug("found matching descriptor for account",
		"account_type", accountType)

	// Filter unspent outputs by checking if they belong to this account
	// We do this by getting address info for each unspent and checking the account
	matchingUnspent := make([]ListUnspentResult, 0)
	for _, unspent := range allUnspent {
		// Check if this address belongs to the account
		// Since we've filtered descriptors by account index, any address in our wallet
		// that matches those descriptors belongs to this account
		belongs, err := b.addressBelongsToAccount(unspent.Address, expectedAccountIndex)
		if err != nil {
			logger.Warn("failed to check address ownership, skipping",
				"address", unspent.Address,
				"error", err)
			continue
		}

		if belongs {
			// Set the label for consistency
			unspent.Label = accountType.String()
			matchingUnspent = append(matchingUnspent, unspent)
		}
	}

	logger.Debug("descriptor matching results",
		"account_type", accountType,
		"matching_unspent", len(matchingUnspent),
		"total_checked", len(allUnspent))

	return matchingUnspent, nil
}

// addressBelongsToAccount checks if an address belongs to an account by verifying its BIP44 path.
//
// This gets the address info from Bitcoin Core and checks if the HDKeyPath or descriptor
// contains the expected account index.
//
// For watch-only multisig descriptors, Bitcoin Core doesn't set HDKeyPath, so we need to
// parse the descriptor to extract the account index.
func (b *Bitcoin) addressBelongsToAccount(address string, expectedAccountIndex uint32) (bool, error) {
	addrInfo, err := b.GetAddressInfo(address)
	if err != nil {
		return false, fmt.Errorf("failed to get address info: %w", err)
	}

	// First, try to use HDKeyPath if available (for non-multisig descriptors)
	if addrInfo.HDKeyPath != "" {
		// Parse the HD key path (e.g., "m/44h/0h/1h/0/5")
		// The account index is the 4th component (index 3 after splitting)
		parts := strings.Split(strings.TrimPrefix(addrInfo.HDKeyPath, "m/"), "/")
		if len(parts) >= 4 {
			// Extract account component (e.g., "1h" or "1'")
			accountComponent := parts[2]
			accountComponent = strings.TrimSuffix(accountComponent, "h")
			accountComponent = strings.TrimSuffix(accountComponent, "H")
			accountComponent = strings.TrimSuffix(accountComponent, "'")

			// Parse the account index
			var accountIndex uint32
			_, err = fmt.Sscanf(accountComponent, "%d", &accountIndex)
			if err == nil {
				return accountIndex == expectedAccountIndex, nil
			}
		}
	}

	// If HDKeyPath is not available, try to parse the descriptor
	// This is needed for watch-only multisig descriptors
	if addrInfo.Desc != "" {
		// Use the existing descriptor matching logic
		return b.descriptorMatchesAccountIndex(addrInfo.Desc, expectedAccountIndex), nil
	}

	// No HD key path or descriptor available
	return false, nil
}

// descriptorMatchesAccountIndex checks if a descriptor contains the expected BIP44 account index.
//
// Descriptor formats we need to handle:
//   - pkh([fingerprint/44h/0h/0h]xpub.../0/*)      -> account index is 0
//   - sh(wpkh([fingerprint/49h/0h/1h]xpub.../0/*)) -> account index is 1
//   - wsh(sortedmulti(2,[fp/48h/0h/2h]xpub.../0/*,[fp/48h/0h/2h]xpub.../0/*)) -> account index is 2
//
// The account index is the 3rd component in the BIP44 path (after coin_type).
func (b *Bitcoin) descriptorMatchesAccountIndex(descriptor string, expectedAccountIndex uint32) bool {
	// Use the descriptor parser for robust parsing
	parser := NewDescriptorParser()
	parsed, err := parser.Parse(descriptor)
	if err != nil {
		// If parsing fails, fall back to simple string matching for compatibility
		logger.Debug("failed to parse descriptor, falling back to string matching",
			"error", err,
			"descriptor_len", len(descriptor))
		return b.descriptorMatchesAccountIndexFallback(descriptor, expectedAccountIndex)
	}

	// Extract account index from the first key's OriginPath
	// All keys in a multisig should have the same account index
	if len(parsed.Keys) == 0 {
		logger.Debug("descriptor has no keys")
		return false
	}

	// Parse the OriginPath to extract account index
	// Format: "/44'/1'/0'" where 0 is the account index (3rd component)
	return b.extractAccountIndexFromOriginPath(parsed.Keys[0].OriginPath, expectedAccountIndex)
}

// extractAccountIndexFromOriginPath extracts the account index from a BIP32 origin path.
// Format: "/44'/1'/0'" -> account index is 0 (3rd component after purpose and coin_type)
//
//nolint:revive // receiver unused but method belongs to Bitcoin type
func (b *Bitcoin) extractAccountIndexFromOriginPath(originPath string, expectedAccountIndex uint32) bool {
	// Split path into components
	parts := strings.Split(strings.TrimPrefix(originPath, "/"), "/")

	// We expect: purpose / coin_type / account
	// Example: 44' / 1' / 0'
	if len(parts) < 3 {
		return false
	}

	// The account index is the 3rd component (index 2)
	accountComponent := parts[2]

	// Remove hardening suffix ('h', 'H', or ')
	accountComponent = strings.TrimSuffix(accountComponent, "h")
	accountComponent = strings.TrimSuffix(accountComponent, "H")
	accountComponent = strings.TrimSuffix(accountComponent, "'")

	// Parse to uint32
	var accountIndex uint32
	_, err := fmt.Sscanf(accountComponent, "%d", &accountIndex)
	if err != nil {
		logger.Debug("failed to parse account index from origin path",
			"origin_path", originPath,
			"account_component", accountComponent,
			"error", err)
		return false
	}

	return accountIndex == expectedAccountIndex
}

// descriptorMatchesAccountIndexFallback is a fallback method using string matching
// when descriptor parsing fails. This handles edge cases and malformed descriptors.
func (*Bitcoin) descriptorMatchesAccountIndexFallback(descriptor string, expectedAccountIndex uint32) bool {
	// Find the first key origin pattern: [fingerprint/path]
	start := strings.Index(descriptor, "[")
	if start == -1 {
		return false
	}

	end := strings.Index(descriptor[start:], "]")
	if end == -1 {
		return false
	}

	// Extract the path (e.g., "fingerprint/44h/0h/1h")
	keyOrigin := descriptor[start+1 : start+end]

	// Split by "/" to get path components
	parts := strings.Split(keyOrigin, "/")

	// We expect at least: fingerprint / purpose / coin_type / account
	// Example: deadbeef / 44h / 0h / 1h
	if len(parts) < 4 {
		return false
	}

	// The account index is the 4th component (index 3)
	accountComponent := parts[3]

	// Parse the account index (remove 'h', 'H', or ' suffix)
	accountComponent = strings.TrimSuffix(accountComponent, "h")
	accountComponent = strings.TrimSuffix(accountComponent, "H")
	accountComponent = strings.TrimSuffix(accountComponent, "'")

	// Convert to uint32
	var accountIndex uint32
	_, err := fmt.Sscanf(accountComponent, "%d", &accountIndex)
	if err != nil {
		return false
	}

	return accountIndex == expectedAccountIndex
}

// LockUnspent lock given txID
// 1st param lock (false)
func (b *Bitcoin) LockUnspent(tx *dtobtc.UnspentOutput) error {
	txIDHash, err := chainhash.NewHashFromStr(tx.TxID)
	if err != nil {
		return fmt.Errorf("fail to call chainhash.NewHashFromStr(%s): %w", tx.TxID, err)
	}
	outpoint := wire.NewOutPoint(txIDHash, tx.Vout)
	err = b.Client.LockUnspent(false, []*wire.OutPoint{outpoint})
	if err != nil {
		return err
	}
	return nil
}

// UnlockUnspent unlock locked unspent tx
// 1st param unlock (true)
func (b *Bitcoin) UnlockUnspent() error {
	list, err := b.Client.ListLockUnspent() // []*wire.OutPoint
	if err != nil {
		return fmt.Errorf("fail to call client.ListLockUnspent(): %w", err)
	}

	if len(list) != 0 {
		err = b.Client.LockUnspent(true, list)
		if err != nil {
			return fmt.Errorf("fail to call client.LockUnspent(): %w", err)
		}
	}

	return nil
}
