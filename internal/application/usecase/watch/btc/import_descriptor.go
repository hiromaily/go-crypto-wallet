package btc

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	portsbtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/retry"
)

type importDescriptorUseCase struct {
	btcClient portsbtc.Bitcoiner
	parser    *btc.DescriptorParser
	chainConf *chaincfg.Params
	addrRepo  watchrepo.AddressRepositorier
	coinType  domainCoin.CoinTypeCode
}

// NewImportDescriptorUseCase creates a descriptor import use case.
func NewImportDescriptorUseCase(
	btcClient portsbtc.Bitcoiner,
	parser *btc.DescriptorParser,
	chainConf *chaincfg.Params,
	addrRepo watchrepo.AddressRepositorier,
	coinType domainCoin.CoinTypeCode,
) watchusecase.ImportDescriptorUseCase {
	return &importDescriptorUseCase{
		btcClient: btcClient,
		parser:    parser,
		chainConf: chainConf,
		addrRepo:  addrRepo,
		coinType:  coinType,
	}
}

func (u *importDescriptorUseCase) Import(
	ctx context.Context,
	input watchusecase.ImportDescriptorInput,
) (watchusecase.ImportDescriptorOutput, error) {
	descriptorStrs, err := u.readDescriptors(input.FilePath)
	if err != nil {
		return watchusecase.ImportDescriptorOutput{}, fmt.Errorf("read descriptors: %w", err)
	}

	var (
		imported  int
		generated int
		errors    []string
	)

	// Parse all descriptors first to validate them
	descriptors := make([]*domainWallet.Descriptor, 0, len(descriptorStrs))
	for _, descStr := range descriptorStrs {
		descriptor, parseErr := u.parser.Parse(descStr)
		if parseErr != nil {
			errors = append(errors, fmt.Sprintf("parse descriptor '%s': %v", descStr, parseErr))
			continue
		}

		if input.ValidateOnly {
			logger.Info("descriptor validated", "descriptor", descStr)
			continue
		}

		descriptors = append(descriptors, descriptor)
	}

	// If validate-only mode, return early
	if input.ValidateOnly {
		return watchusecase.ImportDescriptorOutput{
			DescriptorsImported: len(descriptors),
			AddressesGenerated:  0,
			Errors:              errors,
		}, nil
	}

	// Import descriptors to Bitcoin Core first
	if len(descriptors) > 0 {
		importErrs := u.importToBitcoinCore(descriptorStrs, input.AccountType)
		errors = append(errors, importErrs...)
	}

	// Then derive addresses locally and store them in the database
	for i, descriptor := range descriptors {
		descStr := descriptorStrs[i]

		addrs, deriveErr := u.deriveAddresses(descriptor, input.StartIndex, input.Count)
		if deriveErr != nil {
			errors = append(errors, fmt.Sprintf("derive addresses for descriptor '%s': %v", descStr, deriveErr))
			continue
		}

		if len(addrs) == 0 {
			continue
		}

		if err := u.storeAddresses(ctx, input.AccountType, addrs); err != nil {
			errors = append(errors, fmt.Sprintf("store addresses for descriptor '%s': %v", descStr, err))
			continue
		}

		imported++
		generated += len(addrs)
		logger.Info("descriptor imported", "descriptor", descStr, "addresses_generated", len(addrs))
	}

	return watchusecase.ImportDescriptorOutput{
		DescriptorsImported: imported,
		AddressesGenerated:  generated,
		Errors:              errors,
	}, nil
}

// importToBitcoinCore imports descriptors to Bitcoin Core using the importdescriptors RPC.
//
// This makes the descriptors solvable in Bitcoin Core, allowing transaction creation.
// Returns a list of error messages for any failed imports.
//
// Note: This method adds checksums to descriptors using Bitcoin Core's getdescriptorinfo RPC
// before importing them, since the domain layer's BIP380 checksum implementation is incorrect.
func (u *importDescriptorUseCase) importToBitcoinCore(
	descriptorStrs []string,
	accountType domainAccount.AccountType,
) []string {
	if len(descriptorStrs) == 0 {
		return nil
	}

	// Build import requests for Bitcoin Core
	requests := make([]dtobtc.ImportDescriptorsRequest, 0, len(descriptorStrs))
	for _, descStr := range descriptorStrs {
		// Add checksum using Bitcoin Core if not present
		// (keygen generates descriptors without checksums)
		if !strings.Contains(descStr, "#") {
			info, err := u.btcClient.GetDescriptorInfo(descStr)
			if err != nil {
				logger.Warn("failed to get descriptor checksum, will try importing without it",
					"descriptor", descStr,
					"error", err.Error())
				// Continue with original descriptor (will likely fail at import)
			} else {
				// Use descriptor with correct checksum from Bitcoin Core
				descStr = info.Descriptor
				logger.Debug("added checksum to descriptor",
					"original", descriptorStrs,
					"with_checksum", descStr)
			}
		}

		req := dtobtc.ImportDescriptorsRequest{
			Descriptor: descStr,
			Timestamp:  "now", // Skip rescanning (fastest)
			Active:     true,  // Track outputs for spending
			Internal:   false, // External chain (receiving addresses)
			Watchonly:  true,  // Watch-only wallet (no private keys)
		}

		// Add range for ranged descriptors (with wildcards)
		// Note: Bitcoin Core doesn't allow labels for ranged descriptors
		if strings.Contains(descStr, "/*") {
			req.Range = 1000 // Derive 0-999 addresses
			// Label is intentionally omitted for ranged descriptors
		} else {
			// Only non-ranged descriptors can have labels
			req.Label = accountType.String()
		}

		requests = append(requests, req)
	}

	// Call Bitcoin Core importdescriptors RPC
	responses, err := u.btcClient.ImportDescriptors(requests)
	if err != nil {
		return []string{fmt.Sprintf("failed to call importdescriptors RPC: %v", err)}
	}

	// Process responses and collect errors
	var importErrors []string
	for i, resp := range responses {
		if !resp.Success {
			errMsg := fmt.Sprintf("failed to import descriptor '%s'", descriptorStrs[i])
			if resp.Error != nil {
				errMsg = fmt.Sprintf("%s: [code %d] %s", errMsg, resp.Error.Code, resp.Error.Message)
			}
			importErrors = append(importErrors, errMsg)
			continue
		}

		// Log warnings if any
		if len(resp.Warnings) > 0 {
			logger.Warn("descriptor import warnings",
				"descriptor", descriptorStrs[i],
				"warnings", strings.Join(resp.Warnings, "; "))
		}

		logger.Info("descriptor imported to Bitcoin Core",
			"descriptor", descriptorStrs[i],
			"label", accountType.String())
	}

	return importErrors
}

func (*importDescriptorUseCase) readDescriptors(filePath string) ([]string, error) {
	var data []byte
	var err error

	if strings.TrimSpace(filePath) == "" {
		data, err = ioReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(filePath) //nolint:gosec
	}
	if err != nil {
		return nil, err
	}

	// Try Bitcoin Core JSON format
	var bcFormat []map[string]any
	if jsonErr := json.Unmarshal(data, &bcFormat); jsonErr == nil && len(bcFormat) > 0 {
		var descriptors []string
		for _, item := range bcFormat {
			if desc, ok := item["desc"].(string); ok && desc != "" {
				descriptors = append(descriptors, desc)
			}
		}
		if len(descriptors) > 0 {
			return descriptors, nil
		}
	}

	// Plain text: one descriptor per line, ignoring comments
	var descriptors []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		descriptors = append(descriptors, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan descriptors: %w", err)
	}

	if len(descriptors) == 0 {
		return nil, errors.New("no descriptors found")
	}

	return descriptors, nil
}

func (u *importDescriptorUseCase) deriveAddresses(
	descriptor *domainWallet.Descriptor,
	start uint32,
	count uint32,
) ([]string, error) {
	if len(descriptor.Keys) == 0 {
		return nil, errors.New("descriptor has no keys")
	}

	// Check if multisig descriptor (WSH or SH-WSH)
	key := descriptor.Keys[0]
	isMultisig := len(descriptor.Keys) > 1 &&
		(descriptor.Type == domainWallet.DescriptorTypeWSH ||
			descriptor.Type == domainWallet.DescriptorTypeSHWSH)
	if isMultisig {
		return deriveMultisigAddresses(descriptor, start, count, u.chainConf)
	}

	return deriveSingleKeyAddresses(key, descriptor.Type, start, count, u.chainConf)
}

func deriveChildKeys(key domainWallet.DescriptorKey, start, count uint32) ([]*hdkeychain.ExtendedKey, error) {
	xpub, err := hdkeychain.NewKeyFromString(key.ExtendedPubKey)
	if err != nil {
		return nil, fmt.Errorf("parse extended key: %w", err)
	}

	pathParts := strings.Split(strings.Trim(key.DerivationPath, "/"), "/")

	base := xpub
	for _, part := range pathParts {
		if part == "" || part == "*" {
			continue
		}
		if strings.HasSuffix(part, "'") || strings.HasSuffix(part, "h") || strings.HasSuffix(part, "H") {
			// Assume hardened parts are already baked into the xpub.
			continue
		}
		idx, convErr := strconv.Atoi(part)
		if convErr != nil {
			return nil, fmt.Errorf("invalid derivation segment %s: %w", part, convErr)
		}
		base, err = base.Derive(uint32(idx))
		if err != nil {
			return nil, fmt.Errorf("derive child %d: %w", idx, err)
		}
	}

	children := make([]*hdkeychain.ExtendedKey, 0, count)
	for i := start; i < start+count; i++ {
		child, derr := base.Derive(i)
		if derr != nil {
			return nil, fmt.Errorf("derive address index %d: %w", i, derr)
		}
		children = append(children, child)
	}

	return children, nil
}

func deriveAddressForType(
	key *hdkeychain.ExtendedKey,
	descType domainWallet.DescriptorType,
	chain *chaincfg.Params,
) (string, error) {
	switch descType {
	case domainWallet.DescriptorTypePKH:
		addr, err := key.Address(chain)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case domainWallet.DescriptorTypeWPKH:
		pubKey, err := key.ECPubKey()
		if err != nil {
			return "", err
		}
		hash := btcutil.Hash160(pubKey.SerializeCompressed())
		addr, err := btcutil.NewAddressWitnessPubKeyHash(hash, chain)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case domainWallet.DescriptorTypeSHWPKH:
		pubKey, err := key.ECPubKey()
		if err != nil {
			return "", err
		}
		witHash := btcutil.Hash160(pubKey.SerializeCompressed())
		witAddr, err := btcutil.NewAddressWitnessPubKeyHash(witHash, chain)
		if err != nil {
			return "", err
		}
		script, err := txscript.PayToAddrScript(witAddr)
		if err != nil {
			return "", err
		}
		shAddr, err := btcutil.NewAddressScriptHash(script, chain)
		if err != nil {
			return "", err
		}
		return shAddr.EncodeAddress(), nil
	case domainWallet.DescriptorTypeTR:
		pubKey, err := key.ECPubKey()
		if err != nil {
			return "", err
		}
		tapKey := txscript.ComputeTaprootKeyNoScript(pubKey)
		addr, err := btcutil.NewAddressTaproot(tapKey.SerializeCompressed(), chain)
		if err != nil {
			return "", err
		}
		return addr.EncodeAddress(), nil
	case domainWallet.DescriptorTypeSH:
		return "", errors.New(
			"SH descriptor type requires multisig handling, not supported for single key derivation")
	case domainWallet.DescriptorTypeWSH:
		return "", errors.New(
			"WSH descriptor type requires multisig handling, not supported for single key derivation")
	case domainWallet.DescriptorTypeSHWSH:
		return "", errors.New(
			"SH-WSH descriptor type requires multisig handling, not supported for single key derivation")
	case domainWallet.DescriptorTypeUnknown:
		return "", errors.New("unknown descriptor type")
	default:
		return "", fmt.Errorf("unsupported descriptor type: %s", descType.String())
	}
}

func deriveSingleKeyAddresses(
	key domainWallet.DescriptorKey,
	descType domainWallet.DescriptorType,
	start uint32,
	count uint32,
	chain *chaincfg.Params,
) ([]string, error) {
	paths, err := deriveChildKeys(key, start, count)
	if err != nil {
		return nil, err
	}

	addresses := make([]string, 0, len(paths))
	for _, child := range paths {
		addr, addrErr := deriveAddressForType(child, descType, chain)
		if addrErr != nil {
			return nil, addrErr
		}
		addresses = append(addresses, addr)
	}

	return addresses, nil
}

func deriveMultisigAddresses(
	descriptor *domainWallet.Descriptor,
	start uint32,
	count uint32,
	chain *chaincfg.Params,
) ([]string, error) {
	requiredSigs, err := parseRequiredSignatures(descriptor.Script)
	if err != nil {
		return nil, err
	}

	if requiredSigs <= 0 || requiredSigs > len(descriptor.Keys) {
		return nil, fmt.Errorf("invalid required signatures: %d", requiredSigs)
	}

	addresses := make([]string, 0, count)
	for idx := start; idx < start+count; idx++ {
		pubKeys := make([]*btcec.PublicKey, 0, len(descriptor.Keys))
		for _, key := range descriptor.Keys {
			children, derr := deriveChildKeys(key, idx, 1)
			if derr != nil {
				return nil, derr
			}
			child := children[0]
			pubKey, derr := child.ECPubKey()
			if derr != nil {
				return nil, derr
			}
			pubKeys = append(pubKeys, pubKey)
		}

		// sortedmulti requires lexicographic sorting of pubkeys
		sort.Slice(pubKeys, func(i, j int) bool {
			return bytes.Compare(pubKeys[i].SerializeCompressed(), pubKeys[j].SerializeCompressed()) < 0
		})

		addrPubs := make([]*btcutil.AddressPubKey, 0, len(pubKeys))
		for _, pk := range pubKeys {
			addrPub, addrErr := btcutil.NewAddressPubKey(pk.SerializeCompressed(), chain)
			if addrErr != nil {
				return nil, addrErr
			}
			addrPubs = append(addrPubs, addrPub)
		}

		script, err := txscript.MultiSigScript(addrPubs, requiredSigs)
		if err != nil {
			return nil, fmt.Errorf("failed to create multisig script: %w", err)
		}

		hash := sha256.Sum256(script)
		wshAddr, err := btcutil.NewAddressWitnessScriptHash(hash[:], chain)
		if err != nil {
			return nil, err
		}

		// For P2SH-P2WSH (sh(wsh(...))), wrap the witness script hash in a P2SH
		if descriptor.Type == domainWallet.DescriptorTypeSHWSH {
			// Create P2WSH script (witness program)
			witnessProgram, err := txscript.PayToAddrScript(wshAddr)
			if err != nil {
				return nil, fmt.Errorf("failed to create witness program: %w", err)
			}

			// Wrap in P2SH
			shAddr, err := btcutil.NewAddressScriptHash(witnessProgram, chain)
			if err != nil {
				return nil, fmt.Errorf("failed to create P2SH address: %w", err)
			}

			addresses = append(addresses, shAddr.EncodeAddress())
		} else {
			// Native P2WSH
			addresses = append(addresses, wshAddr.EncodeAddress())
		}
	}

	return addresses, nil
}

func parseRequiredSignatures(script string) (int, error) {
	start := strings.Index(script, "sortedmulti(")
	if start == -1 {
		return 0, errors.New("sortedmulti expression not found for multisig descriptor")
	}
	after := script[start+len("sortedmulti("):]
	parts := strings.SplitN(after, ",", 2)
	if len(parts) < 2 {
		return 0, errors.New("invalid sortedmulti expression")
	}
	required, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("parse required signatures: %w", err)
	}
	return required, nil
}

func (u *importDescriptorUseCase) storeAddresses(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addrs []string,
) error {
	if len(addrs) == 0 {
		return nil
	}

	items := make([]*domainAddress.Address, 0, len(addrs))
	for _, addr := range addrs {
		entity, err := domainAddress.NewAddress(u.coinType, accountType, addr, false)
		if err != nil {
			return fmt.Errorf("build address entity: %w", err)
		}
		items = append(items, entity)
	}

	// Store addresses in database
	if err := u.addrRepo.InsertBulk(ctx, items); err != nil {
		return err
	}

	// Label addresses in Bitcoin Core for getaddressesbylabel lookup
	// Note: This is necessary because ranged descriptor import doesn't support labels
	// (Bitcoin Core limitation), so we must label each derived address individually.
	label := accountType.String()

	var labelErrors []string
	for _, addr := range addrs {
		if err := u.setLabelWithRetry(addr, label); err != nil {
			labelErrors = append(labelErrors, fmt.Sprintf("%s: %v", addr, err))
			logger.Error("failed to label address after retries",
				"address", addr,
				"label", label,
				"error", err)
		}
	}

	// If all addresses failed to be labeled, return error
	if len(labelErrors) == len(addrs) {
		return fmt.Errorf("failed to label all addresses: %s", strings.Join(labelErrors, "; "))
	}

	// If some addresses failed, log warning but continue
	if len(labelErrors) > 0 {
		logger.Warn("some addresses failed to be labeled",
			"failed_count", len(labelErrors),
			"total_count", len(addrs))
	}

	return nil
}

// setLabelWithRetry attempts to set a label with retry and verification.
// It uses the pkg/retry package with exponential backoff.
func (u *importDescriptorUseCase) setLabelWithRetry(addr, label string) error {
	cfg := retry.Config{
		MaxRetries:     2, // 2 retries = 3 total attempts (initial + 2 retries)
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		Multiplier:     2.0,
		Jitter:         false, // Keep false to match original behavior
	}

	operation := func(attempt uint) error {
		// Try to set the label
		if err := u.btcClient.SetLabel(addr, label); err != nil {
			logger.Debug("setLabel failed, will retry",
				"address", addr,
				"label", label,
				"attempt", attempt,
				"error", err)
			return fmt.Errorf("setlabel call failed: %w", err)
		}

		// Verify the label was actually set by checking getaddressinfo
		info, err := u.btcClient.GetAddressInfo(addr)
		if err != nil {
			logger.Debug("label verification failed, will retry",
				"address", addr,
				"label", label,
				"attempt", attempt,
				"error", err)
			return fmt.Errorf("failed to verify label (getaddressinfo failed): %w", err)
		}

		// Check if the label is present in the address labels
		labelFound := false
		for _, addrLabel := range info.Labels {
			if addrLabel == label {
				labelFound = true
				break
			}
		}

		if !labelFound {
			logger.Debug("label not found in address info, will retry",
				"address", addr,
				"expected_label", label,
				"actual_labels", info.Labels,
				"attempt", attempt)
			return fmt.Errorf("label not found after setlabel (expected %s, got %v)", label, info.Labels)
		}

		// Success - label was set and verified
		logger.Debug("successfully labeled address",
			"address", addr,
			"label", label,
			"attempt", attempt)
		return nil
	}

	// Use context with reasonable timeout for label operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return retry.Retry(ctx, cfg, operation)
}

// ioReadAll is a wrapper for io.ReadAll to allow testing via injection.
var ioReadAll = io.ReadAll
