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

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importDescriptorUseCase struct {
	parser    *btc.DescriptorParser
	chainConf *chaincfg.Params
	addrRepo  watchrepo.AddressRepositorier
	coinType  domainCoin.CoinTypeCode
}

// NewImportDescriptorUseCase creates a descriptor import use case.
func NewImportDescriptorUseCase(
	parser *btc.DescriptorParser,
	chainConf *chaincfg.Params,
	addrRepo watchrepo.AddressRepositorier,
	coinType domainCoin.CoinTypeCode,
) watchusecase.ImportDescriptorUseCase {
	return &importDescriptorUseCase{
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

	// Currently support single-key descriptors; multisig descriptors are not yet supported.
	key := descriptor.Keys[0]
	if len(descriptor.Keys) > 1 && descriptor.Type == domainWallet.DescriptorTypeWSH {
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
	case domainWallet.DescriptorTypeWSH:
		return "", errors.New("WSH descriptor type requires multisig handling, not supported for single key derivation")
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

		addresses = append(addresses, wshAddr.EncodeAddress())
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

	return u.addrRepo.InsertBulk(ctx, items)
}

// ioReadAll is a wrapper for io.ReadAll to allow testing via injection.
var ioReadAll = io.ReadAll
